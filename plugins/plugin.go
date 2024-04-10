package plugins

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"math/rand"
	"time"

	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/prices"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	colorReset string = "\033[0m"
	colorRed   string = "\033[31m"
	colorGreen string = "\033[32m"
	colorBlue  string = "\033[34m"
	colorBlack string = "\033[1;30m"
	colorBold  string = ""
)

func getRandomPluginId(base string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("base-%d", r.Intn(10000))
}

type PluginMap map[string]Plugin

func (pm *PluginMap) Select(confs []config.PluginConf) *PluginSelection {
	ps := PluginSelection{}
	for _, conf := range confs {
		p := getPlugin(conf.Name)
		p.SetArgs(conf.Args)
		if conf.Desc != "" {
			p.Desc = conf.Desc
		}
		ps = append(ps, NamedPlugin{getRandomPluginId(conf.Name), conf.Name, *p})
	}
	return &ps
}

func (pm *PluginMap) SelectAll() *PluginSelection {
	ps := PluginSelection{}
	sortedKeys := []string{}
	for k := range *pm {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		ps = append(ps, NamedPlugin{getRandomPluginId(k), k, (*pm)[k]})
	}
	return &ps
}

func (ps *PluginSelection) DocList(doEval bool) string {
	s := ""
	for _, pair := range *ps {
		p := pair.Plugin
		s += fmt.Sprintf("%s%s%-30s%s%s%-20s%s\n", colorGreen, colorBold, pair.Name, colorReset, colorBlue, p.Desc, colorReset)
		if p.ArgDescs != nil {
			s += "  args:\n"
			for _, a := range p.ArgDescs {
				s += fmt.Sprintf("   - %s (%T), default %v\n", a.Desc, a.Default, a.Default)
			}
		}
		if doEval {
			p.SetArgs(p.ArgDescs.ExamplesIf())
			p.Eval()
			if p.Error() != "" {
				s += fmt.Sprintf("  error: %s%s%s\n", colorRed, p.Error(), colorReset)
			} else {
				s += fmt.Sprintf("  output: %s%s%s\n", colorGreen, p.Output(), colorReset)
			}
		}

	}
	return s
}

func (ps *PluginSelection) DocConfig() string {
	s := "plugins:\n"
	for _, p := range *ps {
		s += fmt.Sprintf("  - name: %s\n", p.Name)
		if p.Plugin.ArgDescs != nil {
			s += "    args:\n"
			for _, a := range p.Plugin.ArgDescs {
				s += fmt.Sprintf("      - %s: %v\n", a.Desc, a.Default)
			}
		}

	}
	return s
}

type NamedPlugin struct {
	Id     string
	Name   string
	Plugin Plugin
}

type PluginSelection []NamedPlugin

func (ps *PluginSelection) GetPluginById(id string) *Plugin {
	for _, p := range *ps {
		if p.Id == id {
			return &p.Plugin
		}
	}
	return nil
}

type RefreshFunc func(...interface{}) (interface{}, error)

type ArgDesc struct {
	Desc    string
	Default interface{}
}

type ArgDescs []ArgDesc

func (a ArgDescs) ExamplesIf() []interface{} {
	examples := []interface{}{}
	for _, arg := range a {
		examples = append(examples, arg.Default)
	}
	return examples
}

func (a ArgDescs) HelpString() string {
	s := []string{}
	for _, arg := range a {
		s = append(s, fmt.Sprintf("%s (%T)", arg.Desc, arg.Default))
	}
	return "[" + strings.Join(s, ", ") + "]"
}

type PluginDesc map[string]struct {
	Name string
	Desc string
	Args []string
}

func ValidateAndExpandArgs(args []interface{}, argDescs ArgDescs) ([]interface{}, error) {
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
							return nil, fmt.Errorf("arg #%d (%s): expected float64, got %T %s", i, argDescs[i].Desc, arg, arg)
						}
						args[i] = f
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

type Plugin struct {
	Desc      string
	args      []interface{}
	Help      string
	Formatter func(interface{}) string
	Opts      *PluginOpts
	ArgDescs  ArgDescs
	Refresh   RefreshFunc
	// will be set by Eval()
	err       string
	rawOutput interface{}
	output    string
}

func (p *Plugin) Output() string {
	return p.output
}

func (p *Plugin) Error() string {
	return p.err
}

func (p *Plugin) RawOutput() interface{} {
	return p.rawOutput
}

type PluginOpts struct {
	MarkOutputGreen bool
	MarkNegativeRed bool
}

func (p *Plugin) SetArgs(args []interface{}) {
	p.args = args
}

func (p *Plugin) GetRaw() (interface{}, error) {
	/*
		item := cache.Cache.Get(p.Name)
		if (item != nil) && (!item.IsExpired()) {
			return item.Value(), nil
		}
	*/
	val, err := p.Refresh(p.args...)
	if err != nil {
		return nil, err
	}
	//cache.Cache.Set(p.Name, val, ttlcache.DefaultTTL)
	return val, nil
}

func (p *Plugin) Eval() {
	raw, err := p.GetRaw()
	if err != nil {
		p.err = err.Error()
	}
	if raw != nil {
		p.rawOutput = raw
		p.output = p.Formatter(raw)
	}
}

func ToPlaintext(pl []Plugin) string {
	s := ""
	for _, p := range pl {
		p.Eval()
		l := fmt.Sprintf("%-25s", p.Desc)
		if p.err != "" {
			l += fmt.Sprintf("Error: %s", p.err)
		}
		if p.Output() != "" {
			l += p.Output()
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
	pr := message.NewPrinter(language.English)
	absVal := math.Abs(f)
	if absVal < 1 {
		return pr.Sprintf("%.6f", f)
	}
	if absVal < 2.5 {
		return pr.Sprintf("%.4f", f)
	}
	if absVal < 100 {
		return pr.Sprintf("%.2f", f)
	}
	return pr.Sprintf("%.0f", f)
}

func UintFormatter(i interface{}) string {
	return fmt.Sprintf("%d", i.(uint64))
}

func getPlugin(key string) *Plugin {
	if p, ok := All[key]; ok {
		return &p
	}
	panic(fmt.Sprintf("Plugin not found: %s", key))
}
