package plugins

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const Void = "void"

func buttonize(s1, s2 string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(s1, s2)
}

func (ps *PluginSelection) TelegramFormat(chatId int64, subj string) *tgbotapi.MessageConfig {
	rows := [][]tgbotapi.InlineKeyboardButton{}
	for _, p := range *ps {
		row := []tgbotapi.InlineKeyboardButton{}
		p.Plugin.Eval()
		row = append(row, buttonize(p.Plugin.Desc, Void))
		row = append(row, buttonize(p.Plugin.Output(), p.Id))
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
		descLen := len(p.Plugin.Desc)
		indent := max(25, descLen)
		fString := fmt.Sprintf("%%-%ds ", indent)
		l := fmt.Sprintf(fString, p.Plugin.Desc)
		if p.Plugin.Error() != "" {
			l += fmt.Sprintf("%sError: %s%s", colorRed, p.Plugin.Error(), colorReset)
		}
		if p.Plugin.Output() != "" {
			if p.Plugin.Opts != nil && p.Plugin.Opts.MarkOutputGreen {
				l += fmt.Sprintf("%s%s%s", colorGreen, p.Plugin.Output(), colorReset)
			} else if p.Plugin.Opts != nil && p.Plugin.Opts.MarkNegativeRed && p.Plugin.RawOutput().(float64) < 0 {
				l += fmt.Sprintf("%s%s%s", colorRed, p.Plugin.Output(), colorReset)
			} else {
				l += p.Plugin.Output()
			}
		}
		s += l + "\n"
	}
	return s
}
