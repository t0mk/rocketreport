package plugins

import (
	"github.com/t0mk/rocketreport/exchanges"
)

func RegisterExchangeTickerPlugins() {
	AllPlugins = append(AllPlugins, []Plugin{
		{
			Key:       "kraken",
			Desc:      "Kraken ticker",
			Help:      "Get the latest ticker price from Kraken",
			Formatter: SmartFloatFormatter,
			ArgDescs: []ArgDesc{
				{"ticker", "ticker name", "XETHZUSD"},
			},
			Refresh: func(args ...interface{}) (interface{}, error) {
				ticker := args[0].(string)
				ab, err := exchanges.Kraken(ticker)
				if err != nil {
					return nil, err
				}
				return ab.Ask, nil
			},
		},
	}...)

}
