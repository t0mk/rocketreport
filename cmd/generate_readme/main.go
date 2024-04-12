package main

import (
	"os"
	"strings"

	"github.com/t0mk/rocketreport/plugins"
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

	plugins.RegisterAll()
	pluginTable := plugins.All.SelectAll().MarkdownTable()

	readme := strings.Replace(readmeTemplate, pluginTableMarker, pluginTable, 1)

	err = os.WriteFile("README.md", []byte(readme), 0644)
	if err != nil {
		panic(err)
	}

}
