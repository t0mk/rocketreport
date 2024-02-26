package main

import (
	"github.com/t0mk/rocketreport/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Table represents a table to send
type Table struct {
	Rows [][]string
}

func (t *Table) Format(subj string) tgbotapi.MessageConfig {

	rowFunc := tgbotapi.NewInlineKeyboardRow
	b := func(s string) tgbotapi.InlineKeyboardButton {
		return tgbotapi.NewInlineKeyboardButtonData(s, s)
	}
	rows := [][]tgbotapi.InlineKeyboardButton{}
	for _, r := range t.Rows {
		row := []tgbotapi.InlineKeyboardButton{}
		for _, c := range r {
			row = append(row, b(c))
		}
		rows = append(rows, rowFunc(row...))
	}
	kb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	nm := tgbotapi.NewMessage(config.TelegramChatId, subj)
	nm.DisableWebPagePreview = true
	nm.ParseMode = "Markdown"
	nm.ReplyMarkup = kb
	return nm
}
