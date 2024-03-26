package plugins

import (
	"fmt"

	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/prices"
	"github.com/t0mk/rocketreport/utils"
)

const (
	minipoolDetails = "minipoolDetails"
)

func RegisterAll() {
	RegisterBasicPlugins()
	RegisterValidatorPlugins()
}

func RegisterBasicPlugins() {
	Plugins = append(Plugins, []Plugin{
		{
			Key:       "eth1sync",
			Desc:      "Eth1 client",
			Help:      "Check the sync status of the Eth1 client",
			Formatter: StrFormatter,
			Opts:      &PluginOpts{MarkOutputGreen: true},
			Refresh: func() (interface{}, error) {
				ecs := config.EC.CheckStatus(config.RpConfig)
				return utils.EthClientStatusString(ecs), nil
			},
		},
		{
			Key:       "eth2sync",
			Desc:      "Eth2 client",
			Help:      "Check the sync status of the Eth2 client",
			Formatter: StrFormatter,
			Opts:      &PluginOpts{MarkOutputGreen: true},
			Refresh: func() (interface{}, error) {
				bcs := config.BC.CheckStatus()
				return utils.EthClientStatusString(bcs), nil
			},
		},
		{
			Key:       "actualStake",
			Desc:      "Actual stake",
			Help:      "Check the actual RPL stake",
			Formatter: FloatSuffixFormatter(1, "RPL"),
			Refresh:   GetActualStake,
		},
		{
			Key:       "minStake",
			Desc:      "Minimum stake",
			Help:      "Check the minimum RPL stake",
			Formatter: FloatSuffixFormatter(1, "RPL"),
			Refresh:   GetMinStake,
		},
		{
			Key:       "stakeReserve",
			Desc:      "Stake reserve",
			Help:      "Check the reserve of RPL stake",
			Formatter: FloatSuffixFormatter(1, "RPL"),
			Opts:      &PluginOpts{MarkNegativeRed: true},
			Refresh: func() (interface{}, error) {
				actualStakeRaw, err := getPlugin("actualStake").GetRaw()
				if err != nil {
					return nil, err
				}
				actualStake := actualStakeRaw.(float64)

				minStakeRaw, err := getPlugin("minStake").GetRaw()
				if err != nil {
					return nil, err
				}
				minStake := minStakeRaw.(float64)
				return actualStake - minStake, nil
			},
		},
		{
			Key:       "oracleRplPrice",
			Desc:      "Oracle RPL-ETH",
			Help:      "Check the RPL price from the oracle",
			Formatter: FloatSuffixFormatter(6, "ETH"),
			Refresh:   func() (interface{}, error) { return prices.PriRplEthOracle() },
		},
		{
			Key:       "ethPrice",
			Desc:      fmt.Sprintf("ETH-%s", config.ChosenFiat),
			Help:      fmt.Sprintf("Check ETH/%s price", config.ChosenFiat),
			Formatter: FloatSuffixFormatter(0, config.ChosenFiat.String()),
			Refresh:   func() (interface{}, error) { return prices.PriEth(config.ChosenFiat) },
		},
		{
			Key:       "rplPrice",
			Desc:      fmt.Sprintf("RPL-%s", config.ChosenFiat),
			Help:      fmt.Sprintf("Check RPL/%s price", config.ChosenFiat),
			Formatter: FloatSuffixFormatter(2, config.ChosenFiat.String()),
			Refresh:   func() (interface{}, error) { return prices.PriRpl(config.ChosenFiat) },
		},
		{
			Key:       "ownEthDeposit",
			Desc:      "Own ETH deposit",
			Help:      "Check the amount of ETH deposited",
			Formatter: FloatSuffixFormatter(0, "ETH"),
			Refresh: func() (interface{}, error) {
				mpd, err := CachedGetMinipoolDetails(minipoolDetails)
				if err != nil {
					return nil, err
				}
				return mpd.TotalDeposit, nil
			},
		},
		{
			Key:       "rplFiat",
			Desc:      "RPL funds",
			Help:      fmt.Sprintf("Check the amount of RPL in %s", config.ChosenFiat),
			Formatter: FloatSuffixFormatter(0, config.ChosenFiat.String()),
			Refresh: func() (interface{}, error) {
				rplPriceRaw, err := getPlugin("rplPrice").GetRaw()
				if err != nil {
					return nil, err
				}
				rplPrice := rplPriceRaw.(float64)

				actualStakeRaw, err := getPlugin("actualStake").GetRaw()
				if err != nil {
					return nil, err
				}
				actualStake := actualStakeRaw.(float64)

				return rplPrice * actualStake, nil
			},
		},
	}...)
}
