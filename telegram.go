package main

import (
	"fmt"
	"log"
	"time"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func botInit() error {
	var err error
	bot, err = tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		return err
	}
	return nil
}

func tgString(msg string) error {
	if !sendTg {
		debug("NOT sending tgString():", msg)
		return nil
	}
	log.Println("Sending tgString():", msg)
	if bot != nil {
		ts := time.Now().Format("Mon 02-Jan 15:04")
		hea := fmt.Sprintf("%s, %s", ts, msg)
		nm := tgbotapi.NewMessage(tgChatID, hea)
		_, err := Bot.Send(nm)
		return err
	}
	return nil
}