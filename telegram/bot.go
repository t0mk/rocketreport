package telegram

import (
	"fmt"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/plugins/formatting"
	"github.com/t0mk/rocketreport/plugins/registry"
	"github.com/t0mk/rocketreport/zaplog"
)

func newMsg(chatId int64, text string) *tgbotapi.MessageConfig {
	ret := tgbotapi.NewMessage(chatId, text)
	return &ret
}

func ReportChatID() {
	bot := config.TelegramBot()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}
	botMe, err := bot.GetMe()
	if err != nil {
		panic(err)
	}
	fmt.Println("To find out Chat ID, send a message to your bot (@" + botMe.UserName + ")")
	fmt.Println("https://t.me/" + botMe.UserName)

	update := <-updates
	if update.Message != nil {
		chatId := update.Message.Chat.ID
		txt := fmt.Sprintf("Your Chat ID is:\n%d", chatId)
		fmt.Println(txt)
		msg := newMsg(chatId, txt)
		_, err := bot.Send(*msg)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println()
	fmt.Printf("You can create another Chat at:\nhttps://t.me/%s\n", botMe.UserName)

}

func sendReport(bot *tgbotapi.BotAPI) error {
	ps := registry.Selected

	msg := ps.TelegramFormat(config.TelegramChatID())
	_, err := bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func sendReportAndSchedule(bot *tgbotapi.BotAPI) error {
	err := sendReport(bot)
	if err != nil {
		return err
	}
	if config.TelegramMessageSchedule() != nil {
		// shcedule goroutine for next report
		go func() {
			time.AfterFunc(time.Until(config.TelegramMessageSchedule().Next(time.Now())), func() {
				err := sendReportAndSchedule(bot)
				if err != nil {
					fmt.Println("Error while sending and scheduling: ", err)
				}
			})
		}()
	}
	return nil
}

func RunBot() {
	// header is space-separated keys of either labels or plugins-without args
	log := zaplog.New()
	bot := config.TelegramBot()

	log.Info("Authorized on account: ", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	msg := newMsg(config.TelegramChatID(), "Starting Rocketreport bot")
	kbMarkup := &tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{{
			tgbotapi.KeyboardButton{Text: "Check Now"},
			tgbotapi.KeyboardButton{Text: "Kill Bot"},
			tgbotapi.KeyboardButton{Text: "When Next?"},
		}},
		ResizeKeyboard: true,
	}

	msg.ReplyMarkup = kbMarkup

	_, err = bot.Send(msg)
	if err != nil {
		panic(err)
	}

	err = sendReportAndSchedule(bot)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message != nil {
			if update.Message.Chat.ID != config.TelegramChatID() {
				continue
			}
			if update.Message.Text == "Check Now" {
				err = sendReport(bot)
				if err != nil {
					panic(err)
				}
				continue
			}
			if update.Message.Text == "Kill Bot" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Killing bot")
				_, err = bot.Send(msg)
				if err != nil {
					panic(err)
				}
				fmt.Println("killed by user")
				os.Exit(0)
				return
			}
			if update.Message.Text == "When Next?" {
				whenNextText := "Reports are not scheduled"
				if config.TelegramMessageSchedule() != nil {
					nextTime := config.TelegramMessageSchedule().Next(time.Now())
					untilNext := time.Until(nextTime)
					whenNextText = fmt.Sprintf("Next scheduled report in %s (on %s)",
						formatting.Duration(untilNext), nextTime.Format("Mon 02-Jan 15:04"))
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, whenNextText)
				_, err = bot.Send(msg)
				if err != nil {
					panic(err)
				}
				continue
			}

		} else if update.CallbackQuery != nil {
			fmt.Println("Callback for", update.CallbackQuery.Data)
			pluginId := update.CallbackQuery.Data
			if pluginId == registry.Void {
				continue
			}
			p, err := registry.GetPluginByLabelOrName(pluginId)
			if err != nil {
				fmt.Println("Err while processing callbackj: ", err)
				continue
			}
			if p == nil {
				continue
			}

			p.Eval()
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, p.Output())

			if _, err = bot.AnswerCallbackQuery(callback); err != nil {
				panic(err)
			}
		}
	}
}
