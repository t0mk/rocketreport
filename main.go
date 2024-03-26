package main

import (
	"fmt"
	"time"

	"github.com/alecthomas/kong"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/plugins"
	"github.com/t0mk/rocketreport/prices"
	"github.com/t0mk/rocketreport/zaplog"
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
	log := zaplog.New()
	ethFiat, err := prices.PriEth(config.ChosenFiat())
	if err != nil {
		log.Error("Error getting eth price", err)
	}
	suff := fmt.Sprintf("%s/Ξ", config.ChosenFiat().String())
	ethFiatStr := plugins.FloatSuffixFormatter(0, suff)(ethFiat)

	fmt.Println("sending")
	if s.DoSend {
		ts := time.Now().Format("Mon 02-Jan 15:04")
		subj := fmt.Sprintf("%s - %s", ts, ethFiatStr)
		nm := plugins.Plugins.TelegramFormat(subj)
		_, err := config.TelegramBot().Send(nm)
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
	fmt.Println(plugins.Plugins.TermText())
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
