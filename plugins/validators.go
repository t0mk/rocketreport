package plugins

import (
	"fmt"

	"github.com/t0mk/rocketreport/config"
)

func RegisterValidatorPlugins() {
	Plugins = append(Plugins, []Plugin{
		{
			Key:       "depositedEthFiat",
			Desc:      "Deposited Funds",
			Help:      fmt.Sprintf("Check the amount of dpeosited ETH in %s", config.ChosenFiat),
			Formatter: FloatSuffixFormatter(0, config.ChosenFiat.String()),
			Refresh: func() (interface{}, error) {
				ethPriceRaw, err := getPlugin("ethPrice").GetRaw()
				if err != nil {
					return nil, err
				}
				ethPrice := ethPriceRaw.(float64)

				ownEthDepositRaw, err := getPlugin("ownEthDeposit").GetRaw()
				if err != nil {
					return nil, err
				}
				ownEthDeposit := ownEthDepositRaw.(float64)

				return ethPrice * ownEthDeposit, nil
			},
		},
		{
			Key:       "earnedConsesusEth",
			Desc:      "Earned consensus ETH",
			Help:      fmt.Sprintf("Check the amount of consensus ETH in %s", config.ChosenFiat),
			Formatter: FloatSuffixFormatter(3, "ETH"),
			Refresh: func() (interface{}, error) {
				details, err := CachedGetMinipoolDetails(minipoolDetails)
				if err != nil {
					return nil, err
				}
				return details.Earned, nil
			},
		},
		{
			Key:       "earnedConsensusFunds",
			Desc:      "Earned consensus funds",
			Help:      fmt.Sprintf("Check the amount of consensus funds in %s", config.ChosenFiat),
			Formatter: FloatSuffixFormatter(0, config.ChosenFiat.String()),
			Refresh: func() (interface{}, error) {
				earnedConsesusEthRaw, err := getPlugin("earnedConsesusEth").GetRaw()
				if err != nil {
					return nil, err
				}
				earnedConsesusEth := earnedConsesusEthRaw.(float64)

				ethPriceRaw, err := getPlugin("ethPrice").GetRaw()
				if err != nil {
					return nil, err
				}
				ethPrice := ethPriceRaw.(float64)

				return earnedConsesusEth * ethPrice, nil
			},
		},
		{
			Key:       "totalFunds",
			Desc:      "Total funds",
			Help:      fmt.Sprintf("Check the total amount of funds in %s", config.ChosenFiat),
			Formatter: FloatSuffixFormatter(0, config.ChosenFiat.String()),
			Refresh: func() (interface{}, error) {
				earnedConsensusFundsRaw, err := getPlugin("earnedConsensusFunds").GetRaw()
				if err != nil {
					return nil, err
				}
				earnedConsensusFunds := earnedConsensusFundsRaw.(float64)

				rplFiatRaw, err := getPlugin("rplFiat").GetRaw()
				if err != nil {
					return nil, err
				}
				rplFiat := rplFiatRaw.(float64)

				depositedEthFiatRaw, err := getPlugin("depositedEthFiat").GetRaw()
				if err != nil {
					return nil, err
				}
				depositedEthFiat := depositedEthFiatRaw.(float64)

				return earnedConsensusFunds + rplFiat + depositedEthFiat, nil
			},
		},
	}...)

}
