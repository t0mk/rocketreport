package plugins

import (
	"encoding/json"

	"math/big"

	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/t0mk/rocketreport/utils"
)

const bcGasUrl = "https://beaconcha.in/api/v1/execution/gasnow"

func GasPlugins() map[string]Plugin {
	return map[string]Plugin{
		"gasPrice": {
			Desc:      "Gas Price",
			Help:      "Get the latest gas price",
			Formatter: SmartFloatFormatter,
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
