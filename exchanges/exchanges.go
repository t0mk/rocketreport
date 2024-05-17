package exchanges

import (
	"encoding/json"
	"strconv"

	"github.com/t0mk/rocketreport/utils"
)

type BitfinexTicker []float64

type AskBid struct {
	Ask float64
	Bid float64
}

type ExchangeGetter func(string) (*AskBid, error)

type BinanceTicker struct {
	Ask string `json:"askPrice"`
	Bid string `json:"bidPrice"`
}

func Binance(ticker string) (*AskBid, error) {
	url := "https://api.binance.com/api/v3/ticker/bookTicker?symbol=" + ticker
	body, err := utils.GetHTTPResponseBodyFromUrl(url)
	if err != nil {
		return nil, err
	}
	var tickerData BinanceTicker
	err = json.Unmarshal(body, &tickerData)
	if err != nil {
		return nil, err
	}
	a := tickerData.Ask
	b := tickerData.Bid
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

func Bitfinex(ticker string) (*AskBid, error) {
	url := "https://api-pub.bitfinex.com/v2/ticker/t" + ticker
	body, err := utils.GetHTTPResponseBodyFromUrl(url)
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

func Kraken(ticker string) (*AskBid, error) {
	url := "https://api.kraken.com/0/public/Ticker?pair=" + ticker
	body, err := utils.GetHTTPResponseBodyFromUrl(url)
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

func Coinmate(ticker string) (*AskBid, error) {
	url := "https://coinmate.io/api/ticker?currencyPair=" + ticker
	body, err := utils.GetHTTPResponseBodyFromUrl(url)
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
