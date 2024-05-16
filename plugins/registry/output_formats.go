package registry

import (
	"fmt"

	"slices"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
)

const Void = "void"

func buttonize(s1, s2 string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(s1, s2)
}

func (ps *PluginSelection) TelegramFormat(chatId int64, subj string) *tgbotapi.MessageConfig {
	rows := [][]tgbotapi.InlineKeyboardButton{}
	for _, p := range *ps {
		p.Plugin.Eval()
		if p.Hide {
			continue
		}
		row := []tgbotapi.InlineKeyboardButton{}
		row = append(row, buttonize(p.Plugin.Desc, Void))
		row = append(row, buttonize(p.Plugin.Output(), p.Label))
		rows = append(rows, row)
	}

	kb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	nm := tgbotapi.NewMessage(chatId, subj)
	nm.DisableWebPagePreview = true
	nm.ParseMode = "Markdown"
	nm.ReplyMarkup = kb

	return &nm
}

func (ps *PluginSelection) TextFormat() string {
	s := ""
	for _, p := range *ps {
		p.Plugin.Eval()
		if p.Hide {
			continue
		}
		descLen := len(p.Plugin.Desc)
		indent := max(25, descLen)
		fString := fmt.Sprintf("%%-%ds ", indent)
		desc := fmt.Sprintf(fString, p.Plugin.Desc)
		var value string
		if p.Plugin.Output() != "" {
			out := p.Plugin.Output()
			value = fmt.Sprintf("%s%s%s", formatting.ColorBlue, out, formatting.ColorReset)
			if p.Plugin.Opts != nil {
				if slices.Contains(p.Plugin.Opts, types.OptOkGreen) {
					value = fmt.Sprintf("%s%s%s", formatting.ColorGreen, out, formatting.ColorReset)
				} else if slices.Contains(p.Plugin.Opts, types.OptNegativeRed) {
					if p.Plugin.RawOutput().(float64) < 0 {
						value = fmt.Sprintf("%s%s%s", formatting.ColorRed, out, formatting.ColorReset)
					}
				} else if slices.Contains(p.Plugin.Opts, types.OptRedIfLessThan10) {
					if p.Plugin.RawOutput().(float64) < 10 {
						value = fmt.Sprintf("%s%s%s", formatting.ColorRed, out, formatting.ColorReset)
					}
				}
			}
		}
		if p.Plugin.Error() != "" {
			value = fmt.Sprintf("%sError: %s%s", formatting.ColorRed, p.Plugin.Error(), formatting.ColorReset)
		}
		s += fmt.Sprintf("%s\t%s\n", desc, value)
	}
	return s
}
