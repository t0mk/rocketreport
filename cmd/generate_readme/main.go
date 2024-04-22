package main

import (
	"os"
	"strings"

	"github.com/t0mk/rocketreport/plugins/registry"
)

const (
	readmeTemplatePath = "templates/README.md.template"
	pluginTableMarker  = "__PLUGIN_TABLE__"
)

func main() {
	// templates/README.md.template

	readmeTemplateBytes, err := os.ReadFile("templates/README.md.template")

	if err != nil {
		panic(err)
	}

	readmeTemplate := string(readmeTemplateBytes)

	registry.RegisterAll()
	registry.All.SelectAll()

	tableChapter := ""

	for _, c := range registry.Categories {
		tableChapter += "### " + string(c) + " Plugins\n\n"
		tableChapter += registry.Selected.Cat(c).MarkdownTable() + "\n\n"
	}

	readme := strings.Replace(readmeTemplate, pluginTableMarker, tableChapter, 1)

	err = os.WriteFile("README.md", []byte(readme), 0644)
	if err != nil {
		panic(err)
	}

}
