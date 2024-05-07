package registry

import (
	"fmt"
	"strings"
	"testing"

	"os"

	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/utils"
)

func contentToTempFile(content string) (string, func()) {
	contentReplacedTabs := strings.Replace(content, "\t", "    ", -1)
	f, err := os.CreateTemp("", "*.yml")
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(contentReplacedTabs)
	if err != nil {
		panic(err)
	}
	return f.Name(), func() {
		f.Close()
		os.Remove(f.Name())
	}
}

type Checker func(string) error

type PluginTestCase struct {
	Content string
	Check   Checker
}

var outputChecks = []PluginTestCase{
	{
		`plugins:
  - name: binance
    labl: ethusdt
    args: [ETHUSDT]
	hide: true
  - name: binance
    labl: rplusdt
    args: [RPLUSDT]
  - name: div
    args: [rplusdt, ethusdt]
    desc: RPLETH
`, utils.ShouldntContain("ETHUSDT"),
	},
	{
		`plugins:
  - name: binance
    labl: ethusdt
    args: [ETHUSDT]
  - name: binance
    labl: rplusdt
    args: [RPLUSDT]
  - name: div
    args: [rplusdt, binance]
    desc: RPLETH
`, utils.ShouldContain("Multiple plugins found for reference"),
	},
	{
		`plugins:
  - name: binance
    labl: ethusdt
    args: [ETHUSDT]
  - name: binance
    labl: rplusdt
    args: [RPLUSDT]
  - name: div
    args: [rplusdt, wronglabel]
    desc: RPLETH
`, utils.ShouldContain("No plugin found"),
	},
}

var selectChecksShouldFail = []string{
	`plugins:
  - name: binance
    labl: ethusdt
    args: [ETHUSDT]
  - name : binance
    labl: rplusdt
    args: [RPLUSDT]
  - name: div
    labl: binance
    args: [rplusdt, ethusdt]
`,
	`plugins:
  - name: binance
    labl: div
    args: [ETHUSDT]
  - name: binance
    labl: rplusdt
    args: [RPLUSDT]
  - name: div
    args: [rplusdt, ethusdt]
`,
}

func PrintlnErr(t string) {
	fmt.Fprintln(os.Stderr, t)
}

func TestSelect(t *testing.T) {
	/*
		testConfFile := os.Getenv("TEST_CONF")
		if testConfFile == "" {
			t.Skip("no TEST_CONF env var")
		}
	*/
	RegisterAll()
	for _, c := range selectChecksShouldFail {
		err := SingleTestSelect(c)
		if err == nil {
			t.Errorf("Expected error, got nil:\n%s", c)
		}
		fmt.Println(err)
	}
}

func SingleTestSelect(pluginConfYaml string) error {
	fname, del := contentToTempFile(pluginConfYaml)
	defer del()
	pluginConf := config.FileToPlugins(fname)
	return All.Select(pluginConf)
}

func TestOutputs(t *testing.T) {
	RegisterAll()
	for _, c := range outputChecks {
		fname, del := contentToTempFile(c.Content)
		defer del()
		pluginConf := config.FileToPlugins(fname)
		err := All.Select(pluginConf)
		if err != nil {
			t.Errorf("got error during Select: %s", err)
		}
		txt := Selected.TextFormat()
		err = c.Check(txt)
		if err != nil {
			PrintlnErr("Plugin config:")
			PrintlnErr(c.Content)
			PrintlnErr("Result:")
			PrintlnErr(txt)
			t.Error(err)
		}
	}
}
