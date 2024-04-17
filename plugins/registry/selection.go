package registry

import (
	"fmt"

	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
)

func (ps *PluginSelection) DocList(doEval bool) string {
	s := ""
	for _, pair := range *ps {
		p := pair.Plugin
		s += fmt.Sprintf("%s%s%-30s%s%s%-20s%s\n", formatting.ColorGreen, formatting.ColorBold, pair.Name, formatting.ColorReset, formatting.ColorBlue, p.Desc, formatting.ColorReset)
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
				s += fmt.Sprintf("  error: %s%s%s\n", formatting.ColorRed, p.Error(), formatting.ColorReset)
			} else {
				s += fmt.Sprintf("  output: %s%s%s\n", formatting.ColorGreen, p.Output(), formatting.ColorReset)
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

func (ps *PluginSelection) MarkdownTable() string {
	s := "| Name | Description | Args | Defaults |\n"
	s += "|------|-------------|------|--------------|\n"
	for _, p := range *ps {
		s += fmt.Sprintf("| %s | %s | %s | %s |\n", p.Name, p.Plugin.Help, p.Plugin.ArgDescs.HelpStringDoc(), p.Plugin.ArgDescs.ExamplesString())
	}
	return s
}

func (ps *PluginSelection) FindById(id string) *types.RRPlugin {
	for _, p := range *ps {
		if p.Id == id {
			return &p.Plugin
		}
	}
	return nil
}

func GetPluginById(id string) *types.RRPlugin {
	return Selected.FindById(id)
}
