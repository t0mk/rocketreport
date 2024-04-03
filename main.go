package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/plugins"
)

type Context struct {
	Debug bool
}

type ServePluginCmd struct {
}

func (s *ServePluginCmd) Run(ctx *Context) error {
	RunBot(&plugins.AllPlugins)
	return nil
}

type ListPluginsCmd struct {
	Eval     bool `short:"e" help:"Evaluate all plugins" xor:"template,eval"`
	Template bool `short:"t" help:"Output in template format" xor:"template,eval"`
}

func (l *ListPluginsCmd) Run(ctx *Context) error {
	for _, p := range plugins.AllPlugins {
		line := fmt.Sprintf("%-20s %-20s", p.Key, p.Desc)
		if l.Eval && p.ArgDescs != nil {
			line += fmt.Sprintf("(%s)", p.ArgDescs.ExamplesString())
		}
		if p.ArgDescs != nil {
			line += "\n  args:\n"
			for _, a := range p.ArgDescs {
				line += fmt.Sprintf("   - %s: %s, for example \"%s\"\n", a.Name, a.Desc, a.Example)
			}
		}
		if l.Template {
			line = p.Key
			if p.ArgDescs != nil {
				line = fmt.Sprintf("%s: %s", p.Key, p.ArgDescs.ExamplesString())
			}
		}
		if l.Eval {
			p.Eval(p.ArgDescs.ExamplesIf()...)
			if p.Err != "" {
				line += fmt.Sprintf(" (%s)", p.Err)
			} else {
				line += fmt.Sprintf(" (%s)", p.Output)
			}
		}
		fmt.Println(line)
	}
	return nil
}

type SendCmd struct {
	DoSend bool `short:"s" help:"Send to Telegram"`
}

func (s *SendCmd) Run(ctx *Context) error {
	fmt.Println("sending")
	if s.DoSend {
		nm := plugins.AllPlugins.TelegramFormat(config.TelegramChatID(), tgMsgSubject())
		_, err := config.TelegramBot().Send(*nm)
		return err
	} else {
		fmt.Println("Not sending to Telegram, use -s to send.")
		txt := plugins.ToPlaintext(plugins.AllPlugins)
		fmt.Println(txt)
	}
	return nil
}

type PrintCmd struct{}

var cli struct {
	Send  SendCmd        `cmd:"" help:"Send to configured telegram chat"`
	Print PrintCmd       `cmd:"" help:"Print to stdout"`
	List  ListPluginsCmd `cmd:"" help:"List all plugins"`
	Serve ServePluginCmd `cmd:"" help:"Serve bot"`
}

func (p *PrintCmd) Run(ctx *Context) error {
	fmt.Println(plugins.AllPlugins.TextFormat())
	return nil
}

func main() {
	config.Setup()
	plugins.RegisterAll()
	ctx := kong.Parse(&cli)
	err := ctx.Run(&Context{Debug: true})
	ctx.FatalIfErrorf(err)
	kong.Parse(&cli)
}
