package prices

import (
	"fmt"
	"strings"

	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/t0mk/rocketreport/cache"
	"github.com/t0mk/rocketreport/config"
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
	if f, ok := config.XchMap[denom]; ok {
		pri, err := f(denom)
		if err != nil {
			return 0, fmt.Errorf("error getting price: %v", err)
		}
		cache.Cache.Set("price"+denom, pri, 60*60)
		return pri, nil
	}
	return 0, fmt.Errorf("unsupported denominating currency: %s", denom)
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
