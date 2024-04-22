package rocket

import (
	"fmt"

	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
)

func ValidatorPlugins() map[string]types.RRPlugin {
	return map[string]types.RRPlugin{
		"rpEarnedConsesusEth": {
			Cat:       types.PluginCatRocket,
			Desc:      "Earned consensus ETH",
			Help:      fmt.Sprintf("Check the amount of consensus ETH in %s*", config.ChosenFiat()),
			Formatter: formatting.FloatSuffix(5, "ETH"),
			Refresh: func(...interface{}) (interface{}, error) {
				details, err := CachedGetMinipoolDetails(minipoolDetails)
				if err != nil {
					return nil, err
				}
				return details.Earned, nil
			},
		},
	}
}
