package plugins

import (
	"maps"
)

var All PluginMap

func RegisterAll() {
	All = PluginMap{}
	maps.Copy(All, BasicPlugins())
	maps.Copy(All, ExchangeTickerPlugins())
	maps.Copy(All, ValidatorPlugins())
	maps.Copy(All, GasPlugins())
}