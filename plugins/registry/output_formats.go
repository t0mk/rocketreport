package registry

import (
	"fmt"
	"strings"

	"slices"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/types"
)

const Void = "void"

func buttonize(s1, s2 string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(s1, s2)
}

func getTelegramMessageSubject(labelValue map[string]string) (string, error) {
	subjectFields := config.TelegramHeaderTemplate()
	headerFields := []string{}
	for _, f := range subjectFields {
		if strings.HasPrefix(f, "%") {
			label := strings.TrimPrefix(f, "%")
			if v, ok := labelValue[label]; ok {
				headerFields = append(headerFields, v)
				continue
			}
			if p, ok := AvailablePlugins[label]; ok {
				if len(p.Args) > 0 {
					panic(fmt.Errorf("TELEGRAM_HEADER_TEMPLATE: you can't use plugin %s in header template, it requires arguments", label))
				}
				p.Eval()
				if p.Error() != "" {
					return "", fmt.Errorf("TELEGRAM_HEADER_TEMPLATE: %s", p.Error())
				}
				headerFields = append(headerFields, p.Output())
				continue
			}
		
		} else {
			headerFields = append(headerFields, f)
		}
	}
	return strings.Join(headerFields, " "), nil
}

func (ps *PluginSelection) TelegramFormat(chatId int64) *tgbotapi.MessageConfig {
	rows := [][]tgbotapi.InlineKeyboardButton{}
	labelValue := map[string]string{}
	for _, p := range *ps {
		p.Plugin.Eval()
		if p.Hide {
			continue
		}
		labelValue[p.Label] = p.Plugin.Output()
		row := []tgbotapi.InlineKeyboardButton{}
		row = append(row, buttonize(p.Plugin.Desc, Void))
		row = append(row, buttonize(p.Plugin.Output(), p.Label))
		rows = append(rows, row)
	}

	kb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	subj, err := getTelegramMessageSubject(labelValue)
	if err != nil {
		panic(err)
	}
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
