package rocket

import (
	"time"

	"github.com/rocket-pool/rocketpool-go/rewards"
	"github.com/t0mk/rocketreport/cache"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
	"github.com/t0mk/rocketreport/utils"
)

const (
	minipoolDetails = "minipoolDetails"
)


func RewardPlugins() map[string]types.RRPlugin {
	floatRewardPluginNameDescHelpUnits := [][]string{
		{"rpEffectiveRplStake", "Effective RPL stake", "Check the effective RPL stake of Rocketpool node", "RPL"},
		{"rpTotalRplStake", "Total RPL stake", "Check the total RPL stake of Rocketpool node", "RPL"},
		{"rpEstimatedRewards", "Estimated rewards", "Check the estimated rewards of Rocketpool node", "RPL"},
		{"rpCumulativeRplRewards", "Cumulative RPL rewards", "Check the cumulative RPL rewards of Rocketpool node", "RPL"},
		{"rpCumulativeEthRewards", "Cumulative ETH rewards", "Check the cumulative ETH rewards of Rocketpool node", "ETH"},
		{"rpUnclaimedRplRewards", "Unclaimed RPL rewards", "Check the unclaimed RPL rewards of Rocketpool node", "RPL"},
		{"rpUnclaimedEthRewards", "Unclaimed ETH rewards", "Check the unclaimed ETH rewards of Rocketpool node", "ETH"},
		{"rpBeaconRewards", "Beacon rewards", "Check the beacon rewards of Rocketpool node", "ETH"},
	}
	rewardPlugins := map[string]types.RRPlugin{}
	for _, frp := range floatRewardPluginNameDescHelpUnits {
		rewardPlugins[frp[0]] = GetFloatRewardPlugin(frp[0], frp[1], frp[2], frp[3])
	}
	return rewardPlugins
}

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
			Refresh:   cache.FloatWrap("rplEthOraclePrice", RplEthOraclePrice),
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
		"rpIntervalEnd": types.RRPlugin{
			Cat:       types.PluginCatRocket,
			Desc:      "End of current RP interval",
			Help:      "Check the end of the current Rocketpool interval",
			Formatter: formatting.Time,
			Refresh: func(args ...interface{}) (interface{}, error) {
				start, err := rewards.GetClaimIntervalTimeStart(config.RP(), nil)
				if err != nil {
					return nil, err
				}
				duration, err := rewards.GetClaimIntervalTime(config.RP(), nil)
				if err != nil {
					return nil, err
				}
				return start.Add(duration).UTC(), nil
			},
		},
		"rpUntilIntervalEnd": types.RRPlugin{
			Cat:       types.PluginCatRocket,
			Desc:      "Until end of RP interval",
			Help:      "Check the time until the end of the current Rocketpool interval",
			Formatter: formatting.Duration,
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

				return timeUntilEnd, nil
			},
		},
		"rpOracleRplPriceUpdate": types.RRPlugin{
			Cat:       types.PluginCatRocket,
			Desc:      "Oracle RPL price update",
			Help:      "Time of next RPL price update in Rocketpool oracle",
			Formatter: formatting.Time,
			Refresh:   cache.TimeWrap("rplEthOraclePriceUpdate", NextRplPriceUpdate),
		},
		"rpUntilOracleRplPriceUpdate": types.RRPlugin{
			Cat:       types.PluginCatRocket,
			Desc:      "Until RPL price update",
			Help:      "Time until next RPL price update in Rocketpool oracle",
			Formatter: formatting.Duration,
			Refresh: func(args ...interface{}) (interface{}, error) {
				updateTimeRaw, err := cache.TimeWrap("rplEthOraclePriceUpdate", NextRplPriceUpdate)()
				if err != nil {
					return nil, err
				}
				updateTime := updateTimeRaw.(time.Time)
				timeUntilUpdate := time.Until(updateTime.UTC())
				return timeUntilUpdate, nil
			},
		},
		"rpFeeDistributorBalance": {
			Cat:       types.PluginCatRocket,
			Desc:      "Fee distributor balance",
			Help:      "Check the balance of the Rocketpool fee distributor",
			Formatter: formatting.SmartFloatSuffix("ETH"),
			Refresh:   cache.FloatWrap("feeDistributorBalance", GetFeeDistributorBalance),
		},
	}
}
