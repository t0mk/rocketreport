package prices

import (
	"fmt"
	"strings"

	"github.com/t0mk/rocketreport/cache"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/exchanges"
	"github.com/t0mk/rocketreport/zaplog"
)

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

func PriEth() (float64, error) {
	denom := config.ChosenFiat()
	log := zaplog.New()
	log.Debug("priEthtypes.denom", denom)
	item := cache.Cache.Get("price" + denom)
	if (item != nil) && (!item.IsExpired()) {
		return item.Value().(float64), nil
	}
	var f func(string) (*exchanges.AskBid, error)
	ticker := "ETH" + denom
	switch denom {
	case config.USD:
		f = exchanges.Bitfinex
	case config.EUR:
		f = exchanges.Bitfinex
	case config.GBP:
		f = exchanges.Bitfinex
	case config.JPY:
		f = exchanges.Bitfinex
	case config.AUD:
		f = exchanges.Kraken
	case config.CHF:
		f = exchanges.Kraken
	case config.CZK:
		f = exchanges.Coinmate
		ticker = "ETH_" + denom
	case config.USDT:
		f = exchanges.Binance

	default:
		return 0, fmt.Errorf("unsupported denominating currency: %s", denom)
	}
	ab, err := f(ticker)
	if err != nil {
		return 0, err
	}
	cache.Cache.Set(ticker, ab.Ask, 60)
	return ab.Ask, nil
}
