package registry

import (
	"fmt"
	"maps"
	"math/rand"
	"sort"
	"time"

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

var PluginConfigurations types.PluginConfs
var AvailablePlugins PluginMap

var Selected *PluginSelection
var Categories = []types.PluginCat{
	types.PluginCatRocket,
	types.PluginCatExchange,
	types.PluginCatMeta,
	types.PluginCatCommon,
}

func RegisterAll() {
	AvailablePlugins = PluginMap{}
	maps.Copy(AvailablePlugins, common.ExchangeTickerPlugins())
	maps.Copy(AvailablePlugins, common.GasPlugins())
	maps.Copy(AvailablePlugins, common.PricePlugins())
	maps.Copy(AvailablePlugins, common.BalancePlugins())
	maps.Copy(AvailablePlugins, MetaPlugins())
	maps.Copy(AvailablePlugins, rocket.BasicPlugins())
	maps.Copy(AvailablePlugins, rocket.ValidatorPlugins())
	maps.Copy(AvailablePlugins, rocket.RewardPlugins())
}

func getRandomPluginId(base string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%s-%d", base, r.Intn(10000))
}

func getPlugin(key string) (*types.RRPlugin, error) {
	if p, ok := AvailablePlugins[key]; ok {
		return &p, nil
	}
	return nil, fmt.Errorf("Plugin not found: %s", key)
}

func (pm *PluginMap) PanickingSelect(confs types.PluginConfs) {
	err := pm.Select(confs)
	if err != nil {
		panic(err)
	}
}

func (pm *PluginMap) Select(confs types.PluginConfs) error {
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
			pluginLabel = conf.Hash()
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

/*
func (pm *PluginMap) Select(confs config.PluginConfs) error {
	ps := PluginSelection{}
	labels := map[string]bool{}
	names := map[string]bool{}
	for _, conf := range confs {
		p, err := getPlugin(conf.Name)
		if err != nil {
			return err
		}
		fmt.Println(p)
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
*/

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
