package prices

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/t0mk/rocketreport/cache"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/types"
	"github.com/t0mk/rocketreport/utils"
	"github.com/t0mk/rocketreport/zaplog"
)

var xchMap = map[types.Denom]func(types.Denom) (float64, error){
	types.USD: BitfinexPri,
	types.EUR: BitfinexPri,
	types.GBP: BitfinexPri,
	types.JPY: BitfinexPri,
	types.AUD: KrakenPri,
	types.CHF: KrakenPri,
	types.CZK: CoinmatePri,
}

var currencySymbols = map[string]rune{
	"USD": '$',
	"EUR": '€',
	"GBP": '£',
	"JPY": '¥',
	"ETH": 'Ξ',
}

func FindAndReplaceAllCurrencyOccurencesBySign(s string) string {
	ret := s
	for k, v := range currencySymbols {
		ret = strings.ReplaceAll(ret, k, string(v))
	}
	return ret
}

func PriRplEthOracle() (float64, error) {
	if config.CachedRplPrice != nil {
		return *config.CachedRplPrice, nil
	}
	rplPrice, err := network.GetRPLPrice(config.RP, nil)
	if err != nil {
		return 0, err
	}

	floatRplPrice, _ := utils.WeiToEther(rplPrice).Float64()

	// Return the price
	config.CachedRplPrice = &floatRplPrice
	return floatRplPrice, nil
}

func PriEth(denom types.Denom) (float64, error) {
	log := zaplog.New()
	log.Debug("priEthtypes.denom", denom)
	item := cache.Cache.Get("price" + denom.String())
	if (item != nil) && (!item.IsExpired()) {
		return item.Value().(float64), nil
	}
	if f, ok := xchMap[denom]; ok {
		pri, err := f(denom)
		if err != nil {
			return 0, fmt.Errorf("error getting price: %v", err)
		}
		cache.Cache.Set("price"+denom.String(), pri, 60*60)
		return pri, nil
	}
	return 0, fmt.Errorf("unsupported denominating currency: %s", denom)
}

func PriRpl(denom types.Denom) (float64, error) {
	//rplPrice, err := PriRplEthReal()
	rplPrice, err := PriRplEthOracle()
	if err != nil {
		return 0, err
	}
	ethPrice, err := PriEth(denom)
	if err != nil {
		return 0, err
	}
	return rplPrice * ethPrice, nil
}

func KrakenPri(denom types.Denom) (float64, error) {
	log := zaplog.New()
	log.Debug("bitfinexPri", denom)
	ticker := "ETH" + string(denom)
	ab, err := KrakenGetter(ticker)
	if err != nil {
		return 0, err
	}
	return ab.Ask, nil
}

func BitfinexPri(denom types.Denom) (float64, error) {
	log := zaplog.New()
	log.Debug("bitfinexPri", denom)
	ticker := "ETH" + string(denom)
	ab, err := BitfinexGetter(ticker)
	if err != nil {
		return 0, err
	}
	return ab.Ask, nil
}

func CoinmatePri(denom types.Denom) (float64, error) {
	log := zaplog.New()
	log.Debug("coinmatePri", denom)
	ticker := "ETH_" + string(denom)
	ab, err := CoinmateGetter(ticker)
	if err != nil {
		return 0, err
	}
	return ab.Ask, nil
}

func getHTTPResponseBodyFromUrl(url string) ([]byte, error) {
	log := zaplog.New()
	log.Debug("getHTTPResponseBodyFromUrl", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Get: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	log.Debug("BODY:", string(body))
	return body, nil
}

type BitfinexTicker []float64

type AskBid struct {
	Ask float64
	Bid float64
}

func BitfinexGetter(ticker string) (*AskBid, error) {
	url := "https://api-pub.bitfinex.com/v2/ticker/t" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	var tickerData BitfinexTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return nil, err
	}
	a := tickerData[2]
	b := tickerData[0]
	return &AskBid{a, b}, nil
}

type KrakenTicker struct {
	Result map[string]struct {
		A []string `json:"a"`
		B []string `json:"b"`
	} `json:"result"`
}

func KrakenGetter(ticker string) (*AskBid, error) {
	url := "https://api.kraken.com/0/public/Ticker?pair=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	var tickerData KrakenTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return nil, err
	}
	a := tickerData.Result[ticker].A[0]
	b := tickerData.Result[ticker].B[0]
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return nil, err
	}
	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return nil, err
	}
	return &AskBid{af, bf}, nil
}

type CoinmateTicker struct {
	Data struct {
		Ask float64 `json:"ask"`
		Bid float64 `json:"bid"`
	} `json:"data"`
}

func CoinmateGetter(ticker string) (*AskBid, error) {
	url := "https://coinmate.io/api/ticker?currencyPair=" + ticker
	body, err := getHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	var tickerData CoinmateTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return nil, err
	}
	return &AskBid{tickerData.Data.Ask, tickerData.Data.Bid}, nil
}
