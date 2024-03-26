package plugins

import (
	"fmt"

	"github.com/jellydator/ttlcache/v3"
	"github.com/t0mk/rocketreport/cache"
	"github.com/t0mk/rocketreport/prices"
	"github.com/t0mk/rocketreport/zaplog"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	colorReset string = "\033[0m"
	colorRed   string = "\033[31m"
	colorGreen string = "\033[32m"
)

var Plugins PluginSet

type PluginSet []Plugin

type RefreshFunc func() (interface{}, error)

type Plugin struct {
	Key       string
	Desc      string
	Help      string
	Formatter func(interface{}) string
	Opts      *PluginOpts
	Refresh   RefreshFunc
	// will be set by Eval()
	Err       string
	RawOutput interface{}
	Output    string
}

type PluginOpts struct {
	MarkOutputGreen bool
	MarkNegativeRed bool
}

func (p *Plugin) GetRaw() (interface{}, error) {
	item := cache.Cache.Get(p.Key)
	if (item != nil) && (!item.IsExpired()) {
		return item.Value(), nil
	}
	val, err := p.Refresh()
	if err != nil {
		return nil, err
	}
	cache.Cache.Set(p.Key, val, ttlcache.DefaultTTL)
	return val, nil
}

func (p *Plugin) Eval() {
	log := zaplog.New()
	log.Debug("Evaluating plugin ", p.Key)
	raw, err := p.GetRaw()
	if err != nil {
		p.Err = err.Error()
	}
	if raw != nil {
		p.RawOutput = raw
		p.Output = p.Formatter(raw)
	}
	log.Debug("Evaluating plugin ", p.Key, " done")
}

func ToStringMatrix(pl []Plugin) [][]string {
	s := make([][]string, len(pl))
	for i, p := range pl {
		p.Eval()
		s[i] = []string{p.Desc, p.Output}
	}
	return s
}

func ToPlaintext(pl []Plugin) string {
	s := ""
	for _, p := range pl {
		p.Eval()
		l := fmt.Sprintf("%-25s", p.Desc)
		if p.Err != "" {
			l += fmt.Sprintf("Error: %s", p.Err)
		}
		if p.Output != "" {
			l += p.Output
		}
		s += l + "\n"
	}
	return s
}


func StrFormatter(i interface{}) string {
	return i.(string)
}

func FloatSuffixFormatter(ndecs int, suffix string) func(interface{}) string {
	return func(i interface{}) string {
		f := message.NewPrinter(language.English)
		replacedSuffix := prices.FindAndReplaceAllCurrencyOccurencesBySign(suffix)
		return f.Sprintf("%.*f %s", ndecs, i.(float64), replacedSuffix)
	}
}

func getPlugin(key string) *Plugin {
	for _, p := range Plugins {
		if p.Key == key {
			return &p
		}
	}
	panic(fmt.Sprintf("Plugin not found: %s", key))
}
