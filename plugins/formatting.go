package plugins

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func buttonize(s1, s2 string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(s1, s2)
}

func (ps *PluginSet) TelegramFormat(chatId int64, subj string) *tgbotapi.MessageConfig {
	rows := [][]tgbotapi.InlineKeyboardButton{}
	for _, p := range *ps {
		row := []tgbotapi.InlineKeyboardButton{}
		p.Eval()
		row = append(row, buttonize(p.Desc, ""))
		row = append(row, buttonize(p.Output, ""))
		rows = append(rows, row)
	}

	kb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	nm := tgbotapi.NewMessage(chatId, subj)
	nm.DisableWebPagePreview = true
	nm.ParseMode = "Markdown"
	nm.ReplyMarkup = kb
	return &nm
}

func (ps *PluginSet) TextFormat() string {
	s := ""
	for _, p := range *ps {
		p.Eval()
		l := fmt.Sprintf("%-25s", p.Desc)
		if p.Err != "" {
			l += fmt.Sprintf("%sError: %s%s", colorRed, p.Err, colorReset)
		}
		if p.Output != "" {
			if p.Opts != nil && p.Opts.MarkOutputGreen {
				l += fmt.Sprintf("%s%s%s", colorGreen, p.Output, colorReset)
			} else if p.Opts != nil && p.Opts.MarkNegativeRed && p.RawOutput.(float64) < 0 {
				l += fmt.Sprintf("%s%s%s", colorRed, p.Output, colorReset)
			} else {
				l += p.Output
			}
		}
		s += l + "\n"
	}
	return s
}
