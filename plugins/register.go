package plugins

import (
	"maps"
	"sort"

	"github.com/t0mk/rocketreport/config"
)

type NamedPlugin struct {
	Id     string
	Name   string
	Plugin Plugin
}

type PluginMap map[string]Plugin
type PluginSelection []NamedPlugin

var All PluginMap
var Selected *PluginSelection

func RegisterAll() {
	All = PluginMap{}
	maps.Copy(All, BasicPlugins())
	maps.Copy(All, ExchangeTickerPlugins())
	maps.Copy(All, ValidatorPlugins())
	maps.Copy(All, GasPlugins())
	maps.Copy(All, MetaPlugins())
}

func (pm *PluginMap) Select(confs []config.PluginConf) {
	ps := PluginSelection{}
	for _, conf := range confs {
		p := getPlugin(conf.Name)
		p.SetArgs(conf.Args)
		if conf.Desc != "" {
			p.Desc = conf.Desc
		}
		pluginId := conf.Id
		if ps.FindById(pluginId) != nil {
			panic("Duplicate plugin id: " + pluginId)
		}
		if pluginId == "" {
			pluginId = getRandomPluginId(conf.Name)
		}
		ps = append(ps, NamedPlugin{pluginId, conf.Name, *p})
	}
	Selected = &ps
}

func (pm *PluginMap) SelectAll() {
	ps := PluginSelection{}
	sortedKeys := []string{}
	for k := range *pm {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		ps = append(ps, NamedPlugin{getRandomPluginId(k), k, (*pm)[k]})
	}
	Selected = &ps
}
