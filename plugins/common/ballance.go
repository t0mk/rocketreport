package common

import (
	"fmt"

	"github.com/t0mk/rocketreport/cache"
	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
	"github.com/t0mk/rocketreport/utils"
)

func BallancePlugins() map[string]types.RRPlugin {
	argDescs := types.ArgDescs{
		{Desc: "address", Default: ""},
	}
	return map[string]types.RRPlugin{
		"addressBallance": {
			Cat:       types.PluginCatCommon,
			Desc:      "Address ballance",
			Help:      "Check the ballance of an address",
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
					func() (float64, error) { return utils.AddressBallance(*addr) },
				)
			},
		},
	}
}
