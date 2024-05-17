package rocket

import (
	"time"

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
		{"rpEffectiveRplStake", "Effective RPL stake", "effective RPL stake of Rocketpool node", "RPL"},
		{"rpTotalRplStake", "Total RPL stake", "total RPL stake of Rocketpool node", "RPL"},
		{"rpEstimatedRewards", "Estimated rewards", "estimated rewards of Rocketpool node", "RPL"},
		{"rpCumulativeRplRewards", "Cumulative RPL rewards", "cumulative RPL rewards of Rocketpool node", "RPL"},
		{"rpCumulativeEthRewards", "Cumulative ETH rewards", "cumulative ETH rewards of Rocketpool node", "ETH"},
		{"rpUnclaimedRplRewards", "Unclaimed RPL rewards", "unclaimed RPL rewards of Rocketpool node", "RPL"},
		{"rpUnclaimedEthRewards", "Unclaimed ETH rewards", "unclaimed ETH rewards of Rocketpool node", "ETH"},
		{"rpBeaconRewards", "Beacon rewards", "beacon rewards of Rocketpool node", "ETH"},
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
			Help:      "sync status of Eth1 client (with Rocketpool Golang library)",
			Formatter: formatting.Str,
			Opts:      []string{types.OptOkGreen},
			Refresh: func(...interface{}) (interface{}, error) {
				ecs := config.EC().CheckStatus(config.RpConfig())
				return utils.EthClientStatusString(ecs), nil
			},
		},
		"rpEth2sync": {
			Cat:       types.PluginCatRocket,
			Desc:      "Eth2 client",
			Help:      "sync status of Eth2 client (with Rocketpool Golang library)",
			Formatter: formatting.Str,
			Opts:      []string{types.OptOkGreen},
			Refresh: func(...interface{}) (interface{}, error) {
				bcs := config.BC().CheckStatus()
				return utils.EthClientStatusString(bcs), nil
			},
		},
		"rpEthMatched": {
			Cat:       types.PluginCatRocket,
			Desc:      "Matched ETH",
			Help:      "matched ETH of Rocketpool node",
			Formatter: formatting.FloatSuffix(0, "ETH"),
			Refresh:   cache.FloatWrap("rpEthMatched", GetEthMatched),
		},
		"rpMinStake": {
			Cat:       types.PluginCatRocket,
			Desc:      "Minimum stake",
			Help:      "minimum RPL stake for Rocketpool node",
			Formatter: formatting.SmartFloatSuffix("RPL"),
			Refresh:   cache.FloatWrap("rpMinStake", GetMinStake),
		},
		"rpNodeStake": {
			Cat:       types.PluginCatRocket,
			Desc:      "Node stake",
			Help:      "RPL stake of Rocketpool node",
			Formatter: formatting.SmartFloatSuffix("RPL"),
			Refresh:   cache.FloatWrap("rpNodeStake", GetNodeStake),
		},
		"rpStakeRatio": {
			Cat:       types.PluginCatRocket,
			Desc:      "Stake ratio",
			Help:      "How much % of the borrowed Eth value is staked",
			Opts:      []string{types.OptRedIfLessThan10},
			Formatter: formatting.FloatSuffix(2, "%"),
			Refresh:   cache.FloatWrap("rpStakeRatio", GetStakeRatio),
		},
		"rpOracleRplPrice": {
			Cat:       types.PluginCatRocket,
			Desc:      "Oracle RPL-ETH",
			Help:      "RPL price from Rocketpool oracle",
			Formatter: formatting.FloatSuffix(6, "ETH"),
			Refresh:   cache.FloatWrap("rpEthOraclePrice", GetRplEthOraclePrice),
		},
		"rpOwnEthDeposit": {
			Cat:       types.PluginCatRocket,
			Desc:      "Own ETH deposit",
			Help:      "amount of ETH deposited in Rocketpool node",
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
			Help:      "end of the current Rocketpool interval",
			Formatter: formatting.Time("2006-01-02 15:04:05"),
			Refresh:   cache.TimeWrap("rpIntervalEnd", GetIntervalEnd),
		},
		"rpUntilIntervalEnd": types.RRPlugin{
			Cat:       types.PluginCatRocket,
			Desc:      "Until end of RP interval",
			Help:      "time until the end of the current Rocketpool interval",
			Formatter: formatting.Duration,
			Refresh: func(args ...interface{}) (interface{}, error) {
				intervalEnd, err := cache.Time("rpIntervalEnd", GetIntervalEnd)
				if err != nil {
					return nil, err
				}
				timeUntilEnd := time.Until(intervalEnd.UTC())

				return timeUntilEnd, nil
			},
		},
		"rpOracleRplPriceUpdate": types.RRPlugin{
			Cat:       types.PluginCatRocket,
			Desc:      "Oracle RPL price update",
			Help:      "Time of next RPL price update in Rocketpool oracle",
			Formatter: formatting.Time("2006-01-02 15:04:05"),
			Refresh:   cache.TimeWrap("rpOracleRplPriceUpdate", GetNextRplPriceUpdate),
		},
		"rpUntilOracleRplPriceUpdate": types.RRPlugin{
			Cat:       types.PluginCatRocket,
			Desc:      "Until RPL price update",
			Help:      "Time until next RPL price update in Rocketpool oracle",
			Formatter: formatting.Duration,
			Refresh: func(args ...interface{}) (interface{}, error) {
				updateTime, err := cache.Time("rpOracleRplPriceUpdate", GetNextRplPriceUpdate)
				if err != nil {
					return nil, err
				}
				timeUntilUpdate := time.Until(updateTime.UTC())
				return timeUntilUpdate, nil
			},
		},
		"rpFeeDistributorBalance": {
			Cat:       types.PluginCatRocket,
			Desc:      "Fee distributor balance",
			Help:      "balance of the Rocketpool fee distributor",
			Formatter: formatting.SmartFloatSuffix("ETH"),
			Refresh:   cache.FloatWrap("rpFeeDistributorBalance", GetFeeDistributorBalance),
		},
		"rpNodeBalance": {
			Cat:       types.PluginCatRocket,
			Desc:      "Node balance",
			Help:      "balance of the Rocketpool node",
			Formatter: formatting.SmartFloatSuffix("ETH"),
			Refresh:   cache.FloatWrap("rpNodeBalance", GetNodeBalance),
		},
		"rpWithdrawalAddressBalance": {
			Cat:       types.PluginCatRocket,
			Desc:      "Withdrawal address balance",
			Help:      "balance of the Rocketpool withdrawal address",
			Formatter: formatting.SmartFloatSuffix("ETH"),
			Refresh:   cache.FloatWrap("rpWithdrawalAddressBalance", GetWithdrawalAddressBalance),
		},
	}
}
