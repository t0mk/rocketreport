package registry

import (
	"fmt"
	"strconv"

	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
)

type Reducer func(float64, float64) float64

func CreateMetaPlugin(desc, help string, reducer Reducer) types.RRPlugin {
	return types.RRPlugin{
		Cat:       types.PluginCatMeta,
		Desc:      desc,
		Help:      help,
		Formatter: formatting.SmartFloat,
		ArgDescs: types.ArgDescs{
			{
				Desc:    "list of values - numbers or plugin outputs",
				Default: []interface{}{},
			},
		},
		Refresh: func(args ...interface{}) (interface{}, error) {
			ret, err := GetArgValue(args[0])
			if err != nil {
				return nil, err
			}
			vals, err := ValidateAndExpandMetaArgs(args)
			if err != nil {
				return nil, err
			}
			for _, v := range vals[1:] {
				ret = reducer(ret, v)
			}
			return ret, nil
		},
	}
}

func GetArgValue(arg interface{}) (float64, error) {
	argString, ok := arg.(string)
	if !ok {
		floatVal, ok := arg.(float64)
		if ok {
			return floatVal, nil
		}
		intVal, ok := arg.(int)
		if ok {
			return float64(intVal), nil
		}
	}
	floatVal, err := strconv.ParseFloat(argString, 64)
	if err == nil {
		return floatVal, nil
	}
	intVal, err := strconv.Atoi(argString)
	if err == nil {
		return float64(intVal), nil
	}
	pl, err := GetPluginByLabelOrName(argString)
	if err != nil {
		return 0., fmt.Errorf("error finding plugin by reference \"%s\": %s", arg, err)
	}
	if pl == nil {
		return 0., fmt.Errorf("plugin with Id \"%s\" not found", arg)
	}
	pl.Eval()
	if pl.Error() != "" {
		return 0., fmt.Errorf("error evaluating plugin \"%s\": %s", arg, pl.Error())
	}
	floatResult, ok := pl.RawOutput().(float64)
	if !ok {
		return 0., fmt.Errorf("plugin %s did not return a number", arg)
	}
	return floatResult, nil
}

type MetaVals []float64

func ValidateAndExpandMetaArgs(args []interface{}) (MetaVals, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("expected 2 args, got %d", len(args))
	}
	ret := MetaVals{}
	for i, arg := range args {
		v, err := GetArgValue(arg)
		if err != nil {
			return nil, fmt.Errorf("arg %d (%s): %s", i, arg, err)
		}
		ret = append(ret, v)
	}
	return ret, nil
}

func MetaPlugins() map[string]types.RRPlugin {
	return map[string]types.RRPlugin{
		"add": CreateMetaPlugin(
			"Sum",
			"Sum of given args, either numbers or plugin outputs, adds args and outputs a float",
			func(a, b float64) float64 {
				return a + b
			},
		),
		"mul": CreateMetaPlugin(
			"Multiply",
			"Product of given args, either numbers or plugin outputs, multiplies args and outputs a float",
			func(a, b float64) float64 {
				return a * b
			},
		),
		"sub": CreateMetaPlugin(
			"Subtract",
			"Subtract second arg from first, either numbers or plugin outputs, subtracts args and outputs a float",
			func(a, b float64) float64 {
				return a - b
			},
		),
		"div": CreateMetaPlugin(
			"Divide",
			"Divide first arg by second, either numbers or plugin outputs, divides args and outputs a float",
			func(a, b float64) float64 {
				return a / b
			},
		),
	}
}
