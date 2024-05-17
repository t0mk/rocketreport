package common

import (
	"fmt"

	"github.com/t0mk/rocketreport/cache"
	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
	"github.com/t0mk/rocketreport/utils"
)

func BalancePlugins() map[string]types.RRPlugin {
	argDescs := types.ArgDescs{
		{Desc: "address", Default: ""},
	}
	return map[string]types.RRPlugin{
		"addressBalanceEtherscan": {
			Cat:       types.PluginCatCommon,
			Desc:      "Address balance from Etherscan",
			Help:      "balance of an address using Etherscan",
			ArgDescs:  argDescs,
			Formatter: formatting.SmartFloatSuffix("ETH"),
			Refresh: func(args ...interface{}) (interface{}, error) {
				if len(args) != 1 {
					return "", fmt.Errorf("expected 1 argument, got %d", len(args))
				}
				s := args[0].(string)
				addr, ok := utils.ValidateAndParseAddress(s)
				if !ok {
					return "", fmt.Errorf("invalid address: %s", s)
				}
				return cache.Float(
					"addressBalanceEtherscan"+addr.String(),
					func() (float64, error) { return utils.AddressBalanceEtherscan(*addr) },
				)
			},
		},
		"addressBalance": {
			Cat:       types.PluginCatCommon,
			Desc:      "Address balance",
			Help:      "balance of an address via Execution client",
			ArgDescs:  argDescs,
			Formatter: formatting.SmartFloatSuffix("ETH"),
			Refresh: func(args ...interface{}) (interface{}, error) {
				if len(args) != 1 {
					return "", fmt.Errorf("expected 1 argument, got %d", len(args))
				}
				s := args[0].(string)
				addr, ok := utils.ValidateAndParseAddress(s)
				if !ok {
					return "", fmt.Errorf("invalid address: %s", s)
				}
				return cache.Float(
					"addressBalance"+addr.String(),
					func() (float64, error) { return utils.AddressBalance(*addr) },
				)
			},
		},
	}
}
