package rocket

import (
	"fmt"

	"github.com/rocket-pool/rocketpool-go/rewards"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
	"github.com/t0mk/rocketreport/prices"
	"github.com/t0mk/rocketreport/utils"
)

const (
	minipoolDetails = "minipoolDetails"
)

func BasicPlugins() map[string]types.RRPlugin {
	return map[string]types.RRPlugin{
		"rpEth1sync": {
			Desc:      "[Rocketpool] Eth1 client",
			Help:      "[Rocketpool] Check the sync status of Eth1 client (with Rocketpool Golang library)",
			Formatter: formatting.Str,
			Opts:      &types.Opts{MarkOutputGreen: true},
			Refresh: func(...interface{}) (interface{}, error) {
				ecs := config.EC().CheckStatus(config.RpConfig())
				return utils.EthClientStatusString(ecs), nil
			},
		},
		"rpEth2sync": {
			Desc:      "[Rocketpool] Eth2 client",
			Help:      "[Rocketpool] Check the sync status of Eth2 client (with Rocketpool Golang library)",
			Formatter: formatting.Str,
			Opts:      &types.Opts{MarkOutputGreen: true},
			Refresh: func(...interface{}) (interface{}, error) {
				bcs := config.BC().CheckStatus()
				return utils.EthClientStatusString(bcs), nil
			},
		},
		"rpActualStake": {
			Desc:      "[Rocketpool] Actual stake",
			Help:      "[Rocketpool] Check actual RPL stake of Rocketpool node",
			Formatter: formatting.FloatSuffix(1, "RPL"),
			Refresh:   GetActualStake,
		},
		"rpMinStake": {
			Desc:      "[Rocketpool] Minimum stake",
			Help:      "[Rocketpool] Check the minimum RPL stake for Rocketpool node",
			Formatter: formatting.FloatSuffix(1, "RPL"),
			Refresh:   GetMinStake,
		},
		"rpOracleRplPrice": {
			Desc:      "[Rocketpool] Oracle RPL-ETH",
			Help:      "[Rocketpool] Check the RPL price from Rocketpool oracle",
			Formatter: formatting.FloatSuffix(6, "ETH"),
			Refresh:   func(...interface{}) (interface{}, error) { return prices.PriRplEthOracle() },
		},
		"ethPrice": {
			Desc:      fmt.Sprintf("[Rocketpool] ETH-%s", config.ChosenFiat()),
			Help:      fmt.Sprintf("[Rocketpool] Check ETH/%s* price", config.ChosenFiat()),
			Formatter: formatting.FloatSuffix(0, config.ChosenFiat()),
			Refresh:   func(...interface{}) (interface{}, error) { return prices.PriEth(config.ChosenFiat()) },
		},
		"rplPrice": {
			Desc:      fmt.Sprintf("[Rocketpool] RPL-%s", config.ChosenFiat()),
			Help:      fmt.Sprintf("[Rocketpool] Check RPL/%s* price (RPL/ETH based on Rocketpool Oracle)", config.ChosenFiat()),
			Formatter: formatting.FloatSuffix(2, config.ChosenFiat()),
			Refresh:   func(...interface{}) (interface{}, error) { return prices.PriRpl(config.ChosenFiat()) },
		},
		"ownEthDeposit": {
			Desc:      "[Rocketpool] Own ETH deposit",
			Help:      "[Rocketpool] Check the amount of ETH deposited in Rocketpool node",
			Formatter: formatting.FloatSuffix(0, "ETH"),
			Refresh: func(...interface{}) (interface{}, error) {
				mpd, err := CachedGetMinipoolDetails(minipoolDetails)
				if err != nil {
					return nil, err
				}
				return mpd.TotalDeposit, nil
			},
		},
		"rpIntervalEnds": types.RRPlugin{
			Desc:      "[Rocketpool] End of current interval",
			Help:      "[Rocketpool] Check the end of the current interval",
			Formatter: formatting.Str,
			Refresh: func(args ...interface{}) (interface{}, error) {
				start, err := rewards.GetClaimIntervalTimeStart(config.RP(), nil)
				if err != nil {
					return nil, err
				}
				duration, err := rewards.GetClaimIntervalTime(config.RP(), nil)
				if err != nil {
					return nil, err
				}
				endTimeinString := start.Add(duration).UTC().Format("2006-01-02 15:04:05")
				return endTimeinString, nil
			},
		},
	}
}
