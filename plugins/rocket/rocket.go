package rocket

import (
	"fmt"
	"time"

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
			Cat:       types.PluginCatRocket,
			Desc:      "Eth1 client",
			Help:      "Check the sync status of Eth1 client (with Rocketpool Golang library)",
			Formatter: formatting.Str,
			Refresh: func(...interface{}) (interface{}, error) {
				ecs := config.EC().CheckStatus(config.RpConfig())
				return utils.EthClientStatusString(ecs), nil
			},
		},
		"rpEth2sync": {
			Cat:       types.PluginCatRocket,
			Desc:      "Eth2 client",
			Help:      "Check the sync status of Eth2 client (with Rocketpool Golang library)",
			Formatter: formatting.Str,
			Refresh: func(...interface{}) (interface{}, error) {
				bcs := config.BC().CheckStatus()
				return utils.EthClientStatusString(bcs), nil
			},
		},
		"rpActualStake": {
			Cat:       types.PluginCatRocket,
			Desc:      "Actual stake",
			Help:      "Check actual RPL stake of Rocketpool node",
			Formatter: formatting.FloatSuffix(1, "RPL"),
			Refresh:   GetActualStake,
		},
		"rpMinStake": {
			Cat:       types.PluginCatRocket,
			Desc:      "Minimum stake",
			Help:      "Check the minimum RPL stake for Rocketpool node",
			Formatter: formatting.FloatSuffix(1, "RPL"),
			Refresh:   GetMinStake,
		},
		"rpOracleRplPrice": {
			Cat:       types.PluginCatRocket,
			Desc:      "Oracle RPL-ETH",
			Help:      "Check the RPL price from Rocketpool oracle",
			Formatter: formatting.FloatSuffix(6, "ETH"),
			Refresh:   func(...interface{}) (interface{}, error) { return prices.PriRplEthOracle() },
		},
		"ethPrice": {
			Cat:       types.PluginCatRocket,
			Desc:      fmt.Sprintf("ETH-%s", config.ChosenFiat()),
			Help:      fmt.Sprintf("Check ETH/%s* price", config.ChosenFiat()),
			Formatter: formatting.FloatSuffix(0, config.ChosenFiat()),
			Refresh:   func(...interface{}) (interface{}, error) { return prices.PriEth(config.ChosenFiat()) },
		},
		"rplPrice": {
			Cat:       types.PluginCatRocket,
			Desc:      fmt.Sprintf("RPL-%s", config.ChosenFiat()),
			Help:      fmt.Sprintf("Check RPL/%s* price (RPL/ETH based on Rocketpool Oracle)", config.ChosenFiat()),
			Formatter: formatting.FloatSuffix(2, config.ChosenFiat()),
			Refresh:   func(...interface{}) (interface{}, error) { return prices.PriRpl(config.ChosenFiat()) },
		},
		"rpOwnEthDeposit": {
			Cat:       types.PluginCatRocket,
			Desc:      "Own ETH deposit",
			Help:      "Check the amount of ETH deposited in Rocketpool node",
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
			Cat:       types.PluginCatRocket,
			Desc:      "End of current interval",
			Help:      "Check the end of the current interval",
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
		"rpUntilEndOfInterval": types.RRPlugin{
			Cat:       types.PluginCatRocket,
			Desc:      "Until end of interval",
			Help:      "Check the time until the end of the current interval",
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
				timeUntilEnd := time.Until(start.Add(duration).UTC())

				return formatDuration(timeUntilEnd), nil
			},
		},
		"rpEstimatedRewards": {
			Cat:       types.PluginCatRocket,
			Desc:      "Estimated RPL rewards",
			Help:      "Check the estimated RPL rewards for the current interval",
			Formatter: formatting.FloatSuffix(2, "RPL"),
			Refresh: func(args ...interface{}) (interface{}, error) {
				rewards, err := GetRewards()
				if err != nil {
					return nil, err
				}
				return rewards.EstimatedRewards, nil
			},
		},
	}
}

func formatDuration(duration time.Duration) string {
	days := duration / (time.Hour * 24)
	hours := (duration % (time.Hour * 24)) / time.Hour
	minutes := (duration % time.Hour) / time.Minute

	formatted := ""
	if days > 0 {
		formatted += fmt.Sprintf("%dd ", days)
	}
	if hours > 0 {
		formatted += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 {
		formatted += fmt.Sprintf("%dmin", minutes)
	}

	return formatted
}
