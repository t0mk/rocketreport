package common

import (
	"fmt"
	"strconv"

	"github.com/t0mk/rocketreport/exchanges"
	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
)

func ValidateAndExpandArgs(args []interface{}, argDescs types.ArgDescs) ([]interface{}, error) {
	if len(args) > len(argDescs) {
		return nil, fmt.Errorf("too many arguments, expected %d, got %d", len(argDescs), len(args))
	}
	// fill in defaults
	for i := len(args); i < len(argDescs); i++ {
		args = append(args, argDescs[i].Default)
	}
	for i, arg := range args {
		// check types
		if argDescs[i].Default != nil {
			if _, ok := argDescs[i].Default.(float64); ok {
				if _, ok := arg.(float64); !ok {
					// maybe string that needs to be converted to float
					if s, ok := arg.(string); ok {
						f, err := strconv.ParseFloat(s, 64)
						if err != nil {
							return nil, fmt.Errorf("arg #%d (%s): expected float64, but got %T %s", i, argDescs[i].Desc, arg, arg)
						}
						args[i] = f
					} else if _, ok := arg.(int); ok {
						args[i] = float64(arg.(int))
					} else {
						return nil, fmt.Errorf("arg #%d (%s): expected float64, got %T %s", i, argDescs[i].Desc, arg, arg)
					}
				}
			}
			if _, ok := argDescs[i].Default.(string); ok {
				if _, ok := arg.(string); !ok {
					return nil, fmt.Errorf("arg #%d (%s): expected string, got %T %s", i, argDescs[i].Desc, arg, arg)
				}
			}
		}
	}
	return args, nil
}

func CreateExchangeTickerPlugin(name, exampleTicker string, getter exchanges.ExchangeGetter) types.RRPlugin {
	argDescs := types.ArgDescs{
		{Desc: "ticker", Default: exampleTicker},
		{Desc: "amount", Default: 1.},
	}
	return types.RRPlugin{
		Desc:      fmt.Sprintf("%s ticker", name),
		Help:      fmt.Sprintf("Get the latest ticker price from %s", name),
		Formatter: formatting.SmartFloat,
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

func ExchangeTickerPlugins() map[string]types.RRPlugin {
	return map[string]types.RRPlugin{
		"kraken":   CreateExchangeTickerPlugin("Kraken", "XETHZEUR", exchanges.Kraken),
		"bitfinex": CreateExchangeTickerPlugin("Bitfinex", "ETHEUR", exchanges.Bitfinex),
		"coinmate": CreateExchangeTickerPlugin("Coinmate", "ETH_EUR", exchanges.Coinmate),
	}
}
