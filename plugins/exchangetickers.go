package plugins

import (
	"fmt"

	"github.com/t0mk/rocketreport/exchanges"
)

func CreateExchangeTickerPlugin(name, exampleTicker string, getter exchanges.ExchangeGetter) Plugin {
	argDescs := ArgDescs{
		{"ticker", exampleTicker},
		{"amount", 1.},
	}
	return Plugin{
		Desc:      fmt.Sprintf("%s ticker", name),
		Help:      fmt.Sprintf("Get the latest ticker price from %s", name),
		Formatter: SmartFloatFormatter,
		ArgDescs:  argDescs,
		Refresh: func(args ...interface{}) (interface{}, error) {
			expandedArgs, err := ValidateAndExpandArgs(args, argDescs)
			if err != nil {
				return nil, err
			}
			ticker := expandedArgs[0].(string)
			amount := expandedArgs[1].(float64)
			ab, err := getter(ticker)
			if err != nil {
				return nil, err
			}
			return ab.Ask * amount, nil
		},
	}

}

func ExchangeTickerPlugins() map[string]Plugin {
	return map[string]Plugin{
		"kraken":   CreateExchangeTickerPlugin("Kraken", "XETHZEUR", exchanges.Kraken),
		"bitfinex": CreateExchangeTickerPlugin("Bitfinex", "ETHEUR", exchanges.Bitfinex),
		"coinmate": CreateExchangeTickerPlugin("Coinmate", "ETH_EUR", exchanges.Coinmate),
	}
}
