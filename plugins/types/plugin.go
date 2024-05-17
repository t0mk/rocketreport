package types

import (
	"fmt"
	"strings"

	"github.com/t0mk/rocketreport/zaplog"
	"go.uber.org/zap"
)

type PluginCat string

const (
	PluginCatRocket    PluginCat = "Rocketpool"
	PluginCatMeta      PluginCat = "Meta"
	PluginCatCommon    PluginCat = "Common"
	PluginCatExchange  PluginCat = "Exchange"
	OptOkGreen                   = "okgreen"
	OptNegativeRed               = "negativered"
	OptRedIfLessThan10           = "rediflessthan10"
)

type RRPlugin struct {
	Cat       PluginCat
	Desc      string
	Args      []interface{}
	Help      string
	Formatter func(interface{}) string
	Opts      []string
	ArgDescs  ArgDescs
	Refresh   RefreshFunc
	// will be set by Eval()
	err       string
	rawOutput interface{}
	output    string
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

func (a ArgDescs) ExamplesString() string {
	if len(a) == 0 {
		return ""
	}
	s := []string{}
	for _, arg := range a {
		s = append(s, fmt.Sprintf("%v", arg.Default))
	}
	return strings.Join(s, ", ")
}

func (a ArgDescs) HelpString() string {
	s := []string{}
	for _, arg := range a {
		s = append(s, fmt.Sprintf("%s (%T)", arg.Desc, arg.Default))
	}
	return "[" + strings.Join(s, ", ") + "]"
}

func (a ArgDescs) HelpStringDoc() string {
	s := []string{}
	if len(a) == 0 {
		return ""
	}
	for _, arg := range a {
		s = append(s, fmt.Sprintf("%s (%T)", arg.Desc, arg.Default))
	}
	return strings.Join(s, ", ")
}

type PluginDesc map[string]struct {
	Name string
	Desc string
	Args []string
}

func (p *RRPlugin) Output() string {
	return p.output
}

func (p *RRPlugin) Error() string {
	return p.err
}

func (p *RRPlugin) RawOutput() interface{} {
	return p.rawOutput
}

func (p *RRPlugin) SetArgs(args []interface{}) {
	p.Args = args
}

func (p *RRPlugin) GetRaw() (interface{}, error) {
	log := zaplog.New()
	log.Debug("Refreshing plugin", zap.String("plugin", p.Desc))

	val, err := p.Refresh(p.Args...)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (p *RRPlugin) Eval() {
	raw, err := p.GetRaw()
	if err != nil {
		p.err = err.Error()
	}
	if raw != nil {
		p.rawOutput = raw
		p.output = p.Formatter(raw)
	}
}

func ToPlaintext(pl []RRPlugin) string {
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
