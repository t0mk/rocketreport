package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/params"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Fiat string

const (
	USD Fiat = "USD"
	EUR Fiat = "EUR"
	GBP Fiat = "GBP"
	JPY Fiat = "JPY"
	AUD Fiat = "AUD"
	CHF Fiat = "CHF"
	CZK Fiat = "CZK"
)

func (f Fiat) String() string {
	return string(f)
}

var fiatCache = make(map[Fiat]float64)

var fiatXchMap = map[Fiat]func(Fiat) (float64, error){
	USD: bitfinexPri,
	EUR: bitfinexPri,
	GBP: bitfinexPri,
	JPY: bitfinexPri,
	AUD: krakenPri,
	CHF: krakenPri,
	CZK: coinmatePri,
}

func weiToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.Ether))
}

func fmtEth(p float64) string {
	return fmt.Sprintf("%.6f", p)
}

func fmtRplEur(p float64) string {
	return fmt.Sprintf("%.2f", p)
}

func fmtRpl(p float64) string {
	return fmt.Sprintf("%.1f", p)
}

func priRplEth() (float64, error) {
	if cachedRplPrice != nil {
		return *cachedRplPrice, nil
	}
	/*
		rplPrice, err := network.GetRPLPrice(rp, nil)
		if err != nil {
			return 0, err
		}
	*/
	rplPrice := big.NewInt(0).SetInt64(0x281ee2f7086259)

	floatRplPrice, _ := weiToEther(rplPrice).Float64()

	// Return the price
	cachedRplPrice = &floatRplPrice
	return floatRplPrice, nil
}

func fmtFiat(p float64) string {
	f := message.NewPrinter(language.English)
	i := int(p)
	return f.Sprintf("%d", i)
}

func priEthFiat(fiat Fiat) (float64, error) {
	if debug {
		log.Println("priEthFiat", fiat)
	}
	if f, ok := fiatCache[fiat]; ok {
		return f, nil
	}
	if f, ok := fiatXchMap[fiat]; ok {
		return f(fiat)
	}
	return 0, fmt.Errorf("unsupported fiat: %s", fiat)
}

func priRplFiat(fiat Fiat) (float64, error) {
	rplPrice, err := priRplEth()
	if err != nil {
		return 0, err
	}
	ethPrice, err := priEthFiat(fiat)
	if err != nil {
		return 0, err
	}
	return rplPrice * ethPrice, nil
}

func krakenPri(fiat Fiat) (float64, error) {
	if debug {
		log.Println("bitfinexPri", fiat)
	}
	ticker := "ETH" + string(fiat)
	ab, err := KrakenGetter(ticker)
	if err != nil {
		return 0, err
	}
	return ab.Ask, nil
}

func bitfinexPri(fiat Fiat) (float64, error) {
	if debug {
		log.Println("bitfinexPri", fiat)
	}
	ticker := "ETH" + string(fiat)
	ab, err := BitfinexGetter(ticker)
	if err != nil {
		return 0, err
	}
	return ab.Ask, nil
}

func coinmatePri(fiat Fiat) (float64, error) {
	if debug {
		log.Println("coinmatePri", fiat)
	}
	ticker := "ETH_" + string(fiat)
	ab, err := CoinmateGetter(ticker)
	if err != nil {
		return 0, err
	}
	return ab.Ask, nil
}

func getHTTPResponseBodyFromUrl(url string) ([]byte, error) {
	if debug {
		log.Println("getHTTPResponseBodyFromUrl", url)
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Get: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	if debug {
		log.Println("BODY:", string(body))
	}
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
