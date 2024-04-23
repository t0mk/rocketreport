package registry

import (
	"fmt"

	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
)

func (ps *PluginSelection) DocList(doEval bool) string {
	s := ""
	for _, namedPlugin := range *ps {
		p := namedPlugin.Plugin
		s += fmt.Sprintf("%s%s%-30s%s%s%-20s%s\n", formatting.ColorGreen, formatting.ColorBold, namedPlugin.Name, formatting.ColorReset, formatting.ColorBlue, p.Desc, formatting.ColorReset)
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

func (ps *PluginSelection) Cat(cat types.PluginCat) *PluginSelection {
	ret := PluginSelection{}
	for _, p := range *ps {
		if p.Plugin.Cat == cat {
			ret = append(ret, p)
		}
	}
	return &ret
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

	selectionHasArgs := false
	for _, p := range *ps {
		if p.Plugin.ArgDescs != nil {
			selectionHasArgs = true
			break
		}
	}

	s := "| Name | Description | Args | Defaults |\n"
	s += "|------|-------------|------|--------------|\n"
	if !selectionHasArgs {
		s = "| Name | Description |\n"
		s += "|------|-------------|\n"
	}
	for _, p := range *ps {
		if selectionHasArgs {
			s += fmt.Sprintf("| %s | %s | %s | %s |\n", p.Name, p.Plugin.Help, p.Plugin.ArgDescs.HelpStringDoc(), p.Plugin.ArgDescs.ExamplesString())
		} else {
			s += fmt.Sprintf("| %s | %s |\n", p.Name, p.Plugin.Help)
		}
	}

	return s
}

func (ps *PluginSelection) FindByIdOrName(idOrName string) *types.RRPlugin {
	for _, p := range *ps {
		if (p.Id == idOrName) || (p.Name == idOrName) {
			return &p.Plugin
		}
	}
	return nil
}

func GetPluginByIdOrName(idOrName string) *types.RRPlugin {
	return Selected.FindByIdOrName(idOrName)
}
