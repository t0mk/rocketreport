package plugins

import (
	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
	"github.com/t0mk/rocketreport/utils"
)

func ExtraPlugins() map[string]types.RRPlugin {
	return map[string]types.RRPlugin{
		"smoothingPoolBalance": {
			Desc:      "Smoothing Pool Balance",
			Help:      "ETH in the smoothing pool",
			Formatter: formatting.FloatSuffix(2, "ETH"),
			Refresh: func(...interface{}) (interface{}, error) {
				b, err := utils.SmoothingPoolBalance()
				if err != nil {
					return nil, err
				}
				return *b, nil
			},
		},
	}
}
