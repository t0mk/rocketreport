package registry

import (
	"fmt"
	"maps"
	"math/rand"
	"sort"
	"time"

	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/plugins/common"
	"github.com/t0mk/rocketreport/plugins/rocket"
	"github.com/t0mk/rocketreport/plugins/types"
	"github.com/t0mk/rocketreport/utils"
)

type NamedPlugin struct {
	Label  string
	Name   string
	Hide   bool
	Plugin types.RRPlugin
}

type PluginMap map[string]types.RRPlugin
type PluginSelection []NamedPlugin

var All PluginMap
var Selected *PluginSelection
var Categories = []types.PluginCat{
	types.PluginCatRocket,
	types.PluginCatExchange,
	types.PluginCatMeta,
	types.PluginCatCommon,
}

func RegisterAll() {
	All = PluginMap{}
	maps.Copy(All, common.ExchangeTickerPlugins())
	maps.Copy(All, common.GasPlugins())
	maps.Copy(All, common.PricePlugins())
	maps.Copy(All, common.BalancePlugins())
	maps.Copy(All, MetaPlugins())
	maps.Copy(All, rocket.BasicPlugins())
	maps.Copy(All, rocket.ValidatorPlugins())
	maps.Copy(All, rocket.RewardPlugins())
}

func getRandomPluginId(base string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%s-%d", base, r.Intn(10000))
}

func getPlugin(key string) (*types.RRPlugin, error) {
	if p, ok := All[key]; ok {
		return &p, nil
	}
	return nil, fmt.Errorf("Plugin not found: %s", key)
}

func (pm *PluginMap) Select(confs []config.PluginConf) error {
	ps := PluginSelection{}
	labels := map[string]bool{}
	names := map[string]bool{}
	for _, conf := range confs {
		p, err := getPlugin(conf.Name)
		if err != nil {
			return err
		}
		names[conf.Name] = true
		p.SetArgs(conf.Args)
		//p.Opts = conf.Opts
		if conf.Desc != "" {
			p.Desc = conf.Desc
		} else {
			if conf.Args != nil {
				p.Desc += " " + utils.IfSliceToString(conf.Args)
			}
		}
		pluginLabel := conf.Labl
		if pluginLabel == "" {
			pluginLabel = getRandomPluginId(conf.Name)
		}
		if _, ok := labels[pluginLabel]; ok {
			return fmt.Errorf("Label %s is already used", pluginLabel)
		}
		labels[pluginLabel] = true
		ps = append(ps, NamedPlugin{pluginLabel, conf.Name, conf.Hide, *p})
	}
	for l := range labels {
		if _, ok := names[l]; ok {
			return fmt.Errorf("Label %s clashes with a plugin name", l)
		}
	}
	Selected = &ps
	return nil
}

func (pm *PluginMap) SelectAll() {
	ps := PluginSelection{}
	sortedKeys := []string{}
	for k := range *pm {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		ps = append(ps, NamedPlugin{getRandomPluginId(k), k, false, (*pm)[k]})
	}
	Selected = &ps
}
