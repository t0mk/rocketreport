package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/plugins"
	"github.com/t0mk/rocketreport/utils"
)

type Context struct {
	Debug bool
}

type ServePluginCmd struct {
}

func (s *ServePluginCmd) Run(ctx *Context) error {
	RunBot(plugins.All.Select(config.Plugins))
	return nil
}

type ListPluginsCmd struct {
	Eval           bool `short:"e" help:"Evaluate all plugins" xor:"config,eval"`
	ConfigTemplate bool `short:"c" help:"Output list of all plugins in yaml for plugins.yml" xor:"config,eval"`
}

func (l *ListPluginsCmd) Run(ctx *Context) error {
	allPlugins := plugins.All.SelectAll()
	if l.ConfigTemplate {
		fmt.Println(allPlugins.DocConfig())
		return nil
	}
	fmt.Println(allPlugins.DocList(l.Eval))
	return nil
}

type SendCmd struct {
	DoSend bool `short:"s" help:"Send to Telegram"`
}

func (s *SendCmd) Run(ctx *Context) error {
	fmt.Println("sending")
	if s.DoSend {
		nm := plugins.All.Select(config.Plugins).TelegramFormat(config.TelegramChatID(), tgMsgSubject())
		_, err := config.TelegramBot().Send(*nm)
		return err
	} else {
		fmt.Println("Not sending to Telegram, use -s to send.")
		txt := plugins.All.Select(config.Plugins).TextFormat()
		fmt.Println(txt)
	}
	return nil
}

type PrintCmd struct{}

type RunPluginCmd struct {
	PluginCommandLine []string `arg:"" type:"string" help:"Plugin name and arguments"`
}

func (r *RunPluginCmd) Run(ctx *Context) error {
	pluginName := r.PluginCommandLine[0]
	pluginArgs := r.PluginCommandLine[1:]
	p, ok := plugins.All[pluginName]
	if !ok {
		return fmt.Errorf("plugin %s not found", pluginName)
	}
	p.SetArgs(utils.ToIfSlice(pluginArgs))
	p.Eval()
	if p.Error() != "" {
		return fmt.Errorf("error: %s", p.Error())
	}
	fmt.Println(p.Output())
	return nil
}

var cli struct {
	ConfigFile       string         `type:"path" name:"config" help:"Config file"`
	PluginConfigFile string         `type:"path" name:"plugins" help:"Plugins config" default:"plugins.yaml"`
	Send             SendCmd        `cmd:"" help:"Send to configured telegram chat"`
	Print            PrintCmd       `cmd:"" help:"Print to stdout"`
	Plugins          ListPluginsCmd `cmd:"" help:"List all plugins"`
	Plugin           RunPluginCmd   `cmd:"" help:"List all plugins"`
	Serve            ServePluginCmd `cmd:"" help:"Serve bot"`
}

func (p *PrintCmd) Run(ctx *Context) error {
	fmt.Println(plugins.All.Select(config.Plugins).TextFormat())
	return nil
}

func main() {
	plugins.RegisterAll()
	ctx := kong.Parse(&cli)
	config.Setup(cli.ConfigFile, cli.PluginConfigFile)
	err := ctx.Run(&Context{Debug: true})
	ctx.FatalIfErrorf(err)
	kong.Parse(&cli)
}
