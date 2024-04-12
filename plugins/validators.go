package plugins

import (
	"fmt"

	"github.com/t0mk/rocketreport/config"
)

func ValidatorPlugins() map[string]Plugin {
	return map[string]Plugin{
		"depositedEthFiat": {
			Desc:      "Deposited Funds",
			Help:      fmt.Sprintf("Check the amount of deposited ETH in %s*", config.ChosenFiat()),
			Formatter: FloatSuffixFormatter(0, config.ChosenFiat()),
			Refresh: func(...interface{}) (interface{}, error) {
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
		"earnedConsesusEth": {
			Desc:      "Earned consensus ETH",
			Help:      fmt.Sprintf("Check the amount of consensus ETH in %s*", config.ChosenFiat()),
			Formatter: FloatSuffixFormatter(5, "ETH"),
			Refresh: func(...interface{}) (interface{}, error) {
				details, err := CachedGetMinipoolDetails(minipoolDetails)
				if err != nil {
					return nil, err
				}
				return details.Earned, nil
			},
		},
		"earnedConsensusFunds": {
			Desc:      "Earned consensus funds",
			Help:      fmt.Sprintf("Check the amount of consensus funds in %s*", config.ChosenFiat()),
			Formatter: FloatSuffixFormatter(0, config.ChosenFiat()),
			Refresh: func(...interface{}) (interface{}, error) {
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
		"totalFunds": {
			Desc:      "Total funds",
			Help:      fmt.Sprintf("Check the total amount of funds in %s*", config.ChosenFiat()),
			Formatter: FloatSuffixFormatter(0, config.ChosenFiat()),
			Refresh: func(...interface{}) (interface{}, error) {
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
	}
}
