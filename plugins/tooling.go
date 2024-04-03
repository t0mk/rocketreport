package plugins

import (
	"fmt"
	"math"
	"strings"

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

var AllPlugins PluginSet

type PluginSet []Plugin

type RefreshFunc func(...interface{}) (interface{}, error)

type ArgDesc struct {
	Name    string
	Desc    string
	Example string
}

type ArgDescs []ArgDesc

func (a ArgDescs) ExamplesIf () []interface{} {
	examples := []interface{}{}
	for _, arg := range a {
		examples = append(examples, arg.Example)
	}
	return examples
}

func (a ArgDescs) ExamplesString() string {
	examples := []string{}
	for _, arg := range a {
		examples = append(examples, arg.Example)
	}
	return strings.Join(examples, ", ")
}



type Plugin struct {
	Key       string
	Desc      string
	Help      string
	Formatter func(interface{}) string
	Opts      *PluginOpts
	ArgDescs  ArgDescs
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

func (p *Plugin) GetRaw(args ...interface{}) (interface{}, error) {
	item := cache.Cache.Get(p.Key)
	if (item != nil) && (!item.IsExpired()) {
		return item.Value(), nil
	}
	val, err := p.Refresh(args...)
	if err != nil {
		return nil, err
	}
	cache.Cache.Set(p.Key, val, ttlcache.DefaultTTL)
	return val, nil
}

func (p *Plugin) Eval(args ...interface{}) {
	log := zaplog.New()
	log.Debug("Evaluating plugin ", p.Key)
	raw, err := p.GetRaw(args...)
	if err != nil {
		p.Err = err.Error()
	}
	if raw != nil {
		p.RawOutput = raw
		p.Output = p.Formatter(raw)
	}
	log.Debug("Evaluating plugin ", p.Key, " done")
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

func SmartFloatFormatter(i interface{}) string {
	f := i.(float64)
	absVal := math.Abs(f)
	if absVal < 1 {
		return fmt.Sprintf("%.6f", f)
	}
	if absVal < 2.5 {
		return fmt.Sprintf("%.4f", f)
	}
	if absVal < 100 {
		return fmt.Sprintf("%.2f", f)
	}
	return fmt.Sprintf("%.0f", f)
}

func getPlugin(key string) *Plugin {
	for _, p := range AllPlugins {
		if p.Key == key {
			return &p
		}
	}
	panic(fmt.Sprintf("Plugin not found: %s", key))
}
