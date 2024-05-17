package common

import (
	"context"
	"encoding/json"

	"math/big"

	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
	"github.com/t0mk/rocketreport/utils"
)

const bcGasUrl = "https://beaconcha.in/api/v1/execution/gasnow"

func GasPlugins() map[string]types.RRPlugin {
	return map[string]types.RRPlugin{
		"gasPriceBeaconcha.in": {
			Cat:       types.PluginCatCommon,
			Desc:      "Gas Price",
			Help:      "latest gas price from beaconcha.in",
			Formatter: formatting.SmartFloat,
			Refresh: func(...interface{}) (interface{}, error) {
				body, err := utils.GetHTTPResponseBodyFromUrl(bcGasUrl)
				if err != nil {
					return nil, err
				}
				var gasNow GasNow
				err = json.Unmarshal(body, &gasNow)
				if err != nil {
					return nil, err
				}

				gp := eth.WeiToGwei(big.NewInt(0).SetUint64(gasNow.Data.Rapid))
				return gp, nil
			},
		},
		"gasPriceExecutionClient": {
			Cat:       types.PluginCatCommon,
			Desc:      "Gas Price",
			Help:      "latest gas price from the execution client",
			Formatter: formatting.SmartFloat,
			Refresh: func(...interface{}) (interface{}, error) {
				c := context.Background()
				biGasPrice, err := config.EC().SuggestGasPrice(c)
				if err != nil {
					return nil, err
				}
				gp := eth.WeiToGwei(biGasPrice)
				return gp, nil
			},
		},
	}
}

type GasNow struct {
	Data struct {
		Rapid uint64 `json:"rapid"`
	} `json:"data"`
}
