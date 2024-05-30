package common

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/t0mk/rocketreport/cache"
	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
	"github.com/t0mk/rocketreport/utils"
)

func BalancePlugins() map[string]types.RRPlugin {
	argDescs := types.ArgDescs{
		{Desc: "addresses", Default: []string{}},
	}
	return map[string]types.RRPlugin{
		"addressBalancesEtherscan": {
			Cat:       types.PluginCatCommon,
			Desc:      "Address balance from Etherscan",
			Help:      "Balance of one or more addresses using Etherscan. Use only once, Etherescan has rate limits.",
			ArgDescs:  argDescs,
			Formatter: formatting.SmartFloatSuffix("ETH"),
			Refresh: func(args ...interface{}) (interface{}, error) {
				if len(args) == 0 {
					return "", fmt.Errorf("expected at least 1 argument, got 0")
				}
				addrStrings := []string{}
				for _, arg := range args {
					s, ok := arg.(string)
					if !ok {
						return "", fmt.Errorf("expected string argument, got %T", arg)
					}
					addrStrings = append(addrStrings, s)
				}

				return cache.Float(
					"addressBalanceEtherscan"+strings.Join(addrStrings, ","), func() (float64, error) {
						return utils.AddressBalanceEtherscan(addrStrings)
					},
				)
			},
		},
		"addressBalances": {
			Cat:       types.PluginCatCommon,
			Desc:      "Address balances",
			Help:      "Sum of balances of a list of addresses via Execution client",
			ArgDescs:  argDescs,
			Formatter: formatting.SmartFloatSuffix("ETH"),
			Refresh: func(args ...interface{}) (interface{}, error) {
				addresses := []*common.Address{}
				if len(args) == 0 {
					return "", fmt.Errorf("expected at least 1 argument, got 0")
				}
				for _, arg := range args {
					s, ok := arg.(string)
					if !ok {
						return "", fmt.Errorf("expected string argument, got %T", arg)
					}
					addr, ok := utils.ValidateAndParseAddress(s)
					if !ok {
						return "", fmt.Errorf("invalid address: %s", s)
					}
					addresses = append(addresses, addr)
				}
				total := 0.0
				for _, addr := range addresses {
					balance, err := cache.Float(
						"addressBalance"+addr.String(),
						func() (float64, error) { return utils.AddressBalance(*addr) },
					)
					if err != nil {
						return "", err
					}
					total += balance
				}
				return total, nil
			},
		},
	}
}
