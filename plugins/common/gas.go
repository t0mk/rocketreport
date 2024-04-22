package common

import (
	"encoding/json"

	"math/big"

	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
	"github.com/t0mk/rocketreport/utils"
)

const bcGasUrl = "https://beaconcha.in/api/v1/execution/gasnow"

func GasPlugins() map[string]types.RRPlugin {
	return map[string]types.RRPlugin{
		"gasPrice": {
			Cat:       types.PluginCatCommon,
			Desc:      "Gas Price",
			Help:      "Get the latest gas price",
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
	}
}

type GasNow struct {
	Data struct {
		Rapid uint64 `json:"rapid"`
	} `json:"data"`
}
