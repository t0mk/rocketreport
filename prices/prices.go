package prices

import (
	"fmt"
	"strings"

	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/t0mk/rocketreport/cache"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/exchanges"
	"github.com/t0mk/rocketreport/utils"
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

func PriRplEthOracle() (float64, error) {
	if config.CachedRplPrice != nil {
		return *config.CachedRplPrice, nil
	}
	rplPrice, err := network.GetRPLPrice(config.RP(), nil)
	if err != nil {
		return 0, err
	}

	floatRplPrice, _ := utils.WeiToEther(rplPrice).Float64()

	// Return the price
	config.CachedRplPrice = &floatRplPrice
	return floatRplPrice, nil
}

func PriEth(denom string) (float64, error) {
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

func PriRpl(denom string) (float64, error) {
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
