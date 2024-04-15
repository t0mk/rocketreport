package plugins

import (
	"fmt"
	"strconv"
)

type Reducer func(float64, float64) float64

func CreateMetaPlugin(desc, help string, reducer Reducer, seed float64) Plugin {
	return Plugin{
		Desc:      desc,
		Help:      help,
		Formatter: FloatSuffixFormatter(0, ""),
		ArgDescs: ArgDescs{
			{"list of values - numbers or plugin outputs", []interface{}{}},
		},
		Refresh: func(args ...interface{}) (interface{}, error) {
			ret := seed
			vals, err := ValidateAndExpandMetaArgs(args)
			if err != nil {
				return nil, err
			}
			for _, v := range vals {
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
	pl := GetPluginById(argString)
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

func simpleSum(a, b float64) float64 {
	return a + b
}

func simpleProd(a, b float64) float64 {
	return a * b
}

func MetaPlugins() map[string]Plugin {
	return map[string]Plugin{
		"sum": CreateMetaPlugin(
			"Sum",
			"Sum of given args, either numbers or plugin outputs, adds args and outputs a float",
			simpleSum,
			0,
		),
		"prod": CreateMetaPlugin(
			"Product",
			"Product of given args, either numbers or plugin outputs, multiplies args and outputs a float",
			simpleProd,
			1,
		),
	}
}
