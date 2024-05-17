package types

import (
	"fmt"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/t0mk/rocketreport/utils"
)

type PluginConf struct {
	Name    string        `yaml:"name" json:"name"`
	Desc    string        `yaml:"desc" json:"desc"`
	Labl    string        `yaml:"labl" json:"labl"`
	Args    []interface{} `yaml:"args" json:"args"`
	Hide    bool          `yaml:"hide" json:"hide"`
}

func (pc PluginConf) Hash() string {
	argsString := utils.IfSliceToString(pc.Args)
	return fmt.Sprintf("%s:%s", pc.Name, argsString)
}

type PluginConfs []PluginConf

func (pcs PluginConfs) String() string {
	s := ""
	for i, p := range pcs {
		if p.Name != "" {
			s += fmt.Sprintf("name: %s\n", p.Name)
		}
		if p.Desc != "" {
			s += fmt.Sprintf("desc: %s\n", p.Desc)
		}
		if p.Labl != "" {
			s += fmt.Sprintf("labl: %s\n", p.Labl)
		}
		if p.Args != nil {
			s += "args:\n"
			for _, a := range p.Args {
				s += fmt.Sprintf("  - %v\n", a)
			}
		}
		if p.Hide {
			s += "hide: true\n"
		}
		if i < len(pcs)-1 {
			s += "----\n"
		}
	}
	return s
}

func PluginsString(pcs []PluginConf) string {
	s := "plugins:\n"
	for _, p := range pcs {
		s += fmt.Sprintf("  - name: %s\n", p.Name)
		if p.Args != nil {
			s += "    args:\n"
			for _, a := range p.Args {
				s += fmt.Sprintf("      - %v\n", a)
			}
		}
	}
	return s
}

func FileToPlugins(file string) []PluginConf {
	pluginsWrap := struct {
		Plugins []PluginConf `yaml:"plugins"`
	}{}
	loader := aconfig.LoaderFor(&pluginsWrap, aconfig.Config{
		SkipFlags: true,
		SkipEnv:   true,
		Files:     []string{file},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
			".yml":  aconfigyaml.New(),
		},
	})
	err := loader.Load()
	if err != nil {
		panic(err)
	}
	return pluginsWrap.Plugins

}
