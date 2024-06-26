package common

import (
	"fmt"

	"github.com/t0mk/rocketreport/cache"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/exchanges"
	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
	"github.com/t0mk/rocketreport/prices"
)

func PriRplReal() (float64, error) {
	abRplUsdt, err := exchanges.Binance("RPLUSDT")
	if err != nil {
		return 0, err
	}
	abEthUsdt, err := exchanges.Binance("ETHUSDT")
	if err != nil {
		return 0, err
	}
	return abRplUsdt.Ask / abEthUsdt.Ask, nil
}

func PricePlugins() map[string]types.RRPlugin {
	return map[string]types.RRPlugin{
		"ethPrice": {
			Cat:       types.PluginCatCommon,
			Desc:      fmt.Sprintf("ETH-%s", config.ChosenFiat()),
			Help:      fmt.Sprintf("ETH/%s* price", config.ChosenFiat()),
			//Formatter: formatting.FloatSuffix(0, config.ChosenFiat()),
			Formatter: formatting.SmartFloat,
			//Refresh:   func(...interface{}) (interface{}, error) { return prices.PriEth(config.ChosenFiat()) },
			Refresh:   cache.FloatWrap("ethPrice", prices.PriEth),
		},
		"rplPriceRealtime": {
			Cat:       types.PluginCatCommon,
			Desc:      "RPL-ETH from Exchange",
			Help:      "Realtime RPL-ETH (based on RPL-USDT and ETH-USDT from Binance)",
			Formatter: formatting.SmartFloatSuffix("ETH"),
			Refresh:   cache.FloatWrap("rplPriceRealtime", PriRplReal),
		},
	}
}
