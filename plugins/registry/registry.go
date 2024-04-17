package registry

import (
	"fmt"
	"maps"
	"math/rand"
	"sort"
	"time"

	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/plugins/rocket"
	"github.com/t0mk/rocketreport/plugins/types"
	"github.com/t0mk/rocketreport/plugins/common"
)

type NamedPlugin struct {
	Id     string
	Name   string
	Plugin types.RRPlugin
}

type PluginMap map[string]types.RRPlugin
type PluginSelection []NamedPlugin

var All PluginMap
var Selected *PluginSelection

func RegisterAll() {
	All = PluginMap{}
	maps.Copy(All, common.ExchangeTickerPlugins())
	maps.Copy(All, common.GasPlugins())
	maps.Copy(All, MetaPlugins())
	maps.Copy(All, rocket.BasicPlugins())
	maps.Copy(All, rocket.ValidatorPlugins())
}

func getRandomPluginId(base string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%s-%d", base, r.Intn(10000))
}

func getPlugin(key string) *types.RRPlugin {
	if p, ok := All[key]; ok {
		return &p
	}
	panic(fmt.Sprintf("Plugin not found: %s", key))
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
