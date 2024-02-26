package main

import (
	"fmt"
	"time"

	"github.com/alecthomas/kong"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/plugins"
)

type Context struct {
	Debug bool
}

type ListPluginsCmd struct {
	Eval bool `short:"e" help:"Evaluate all plugins"`
}

func (l *ListPluginsCmd) Run(ctx *Context) error {
	for _, p := range plugins.Plugins {
		line := fmt.Sprintf("%-20s %-20s", p.Key, p.Desc)
		if l.Eval {
			p.Eval()
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
	if s.DoSend && config.Bot != nil {
		sm := plugins.ToStringMatrix(plugins.Plugins)
		tab := Table{Rows: sm}
		ts := time.Now().Format("Mon 02-Jan 15:04")
		nm := tab.Format(ts)
		_, err := config.Bot.Send(nm)
		return err
	} else {
		fmt.Println("Not sending to Telegram, use -s to send.")
		txt := plugins.ToPlaintext(plugins.Plugins)
		fmt.Println(txt)
	}
	return nil
}

type PrintCmd struct{}

var cli struct {
	Send  SendCmd        `cmd:"" help:"Send to configured telegram chat"`
	Print PrintCmd       `cmd:"" help:"Print to stdout"`
	List  ListPluginsCmd `cmd:"" help:"List all plugins"`
}

func (p *PrintCmd) Run(ctx *Context) error {
	fmt.Println(plugins.ToTermText(plugins.Plugins))
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
