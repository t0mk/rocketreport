package exchanges

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/t0mk/rocketreport/zaplog"
)

func KrakenPri(denom string) (float64, error) {
	log := zaplog.New()
	log.Debug("bitfinexPri", denom)
	ticker := "ETH" + string(denom)
	ab, err := KrakenGetter(ticker)
	if err != nil {
		return 0, err
	}
	return ab.Ask, nil
}

func BitfinexPri(denom string) (float64, error) {
	log := zaplog.New()
	log.Debug("bitfinexPri", denom)
	ticker := "ETH" + string(denom)
	ab, err := BitfinexGetter(ticker)
	if err != nil {
		return 0, err
	}
	return ab.Ask, nil
}

func CoinmatePri(denom string) (float64, error) {
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
