package main

import (
	"os"
	"strings"

	"github.com/t0mk/rocketreport/plugins/registry"
)

func getSectionHeader(c string) string {
	return c + " Plugins\n"
}

func getSectionLink(c string) string {
	withDashes := strings.ReplaceAll(strings.ToLower(c), " ", "-")
	return "- [" + getSectionHeader(c) + "](#" + withDashes + "-plugins)\n"
}

func main() {
	registry.RegisterAll()
	registry.All.SelectAll()

	plugins_md := `
# Rocketreport Plugins

`
	for _, c := range registry.Categories {
		plugins_md += getSectionLink(string(c))
	}

	for _, c := range registry.Categories {
		plugins_md += "## " + getSectionHeader(string(c))
		plugins_md += registry.Selected.Cat(c).MarkdownTable() + "\n\n"
	}

	plugins_md += "\n\n&ast; you can use different fiat as quote currency in these plugins if you set \"fiat\" option in config.yml"

	err := os.WriteFile("PLUGINS.md", []byte(plugins_md), 0644)
	if err != nil {
		panic(err)
	}

}
