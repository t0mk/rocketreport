package telegram

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/t0mk/rocketreport/config"
	"github.com/t0mk/rocketreport/plugins"
	"github.com/t0mk/rocketreport/prices"
	"github.com/t0mk/rocketreport/zaplog"
)

func strIs(s string, patterns ...string) bool {
	for _, p := range patterns {
		if strings.EqualFold(s, p) {
			return true
		}
	}
	return false
}

type BotCommand struct {
	Key  string
	Desc string
	Func func(chatId string) *tgbotapi.MessageConfig
}

/*
type BotLogic struct {
	Plugins  *plugins.PluginSet
	Commands []BotCommand
}

func (bl *BotLogic) Help() *tgbotapi.MessageConfig {
	text := "Available commands:\n"
	for _, c := range bl.Commands {
		text += c.Key + " - " + c.Desc + "\n"
	}
	return &tgbotapi.MessageConfig{
		Text: text,
	}
}

func getLogic(ps plugins.PluginSet) *BotLogic {
	bl := &BotLogic{
		Plugins: &ps,
		Commands: []BotCommand{
			{
				Key:  "open",
				Desc: "Open the bot",
				Func: func(chatId) *tgbotapi.MessageConfig {
					return tgbotapi.NewMessage()
				},
			},
			{
				Key:  "help",
				Desc: "Show help",
				Func: func() *tgbotapi.MessageConfig {
					return bl.Help()
				},
			},
		},
	}
	return bl
}

func allPrefixes(s string) []string {
	prefixes := []string{}
	for i := 1; i < len(s); i++ {
		prefixes = append(prefixes, s[:i])
	}
	return prefixes
}

func (bl *BotLogic) getCommandFunc(command string) func() *tgbotapi.MessageConfig {
	for _, c := range bl.Commands {
		if strIs(command, allPrefixes()...) {
			return c.Func
		}
	}
	return nil
}
*/

func newMsg(chatId int64, text string) *tgbotapi.MessageConfig {
	ret := tgbotapi.NewMessage(chatId, text)
	return &ret
}

func MessageSubject() string {
	ethFiat, err := prices.PriEth(config.ChosenFiat())
	if err != nil {
		panic(err)
	}
	suff := fmt.Sprintf("%s/Ξ", config.ChosenFiat())
	ts := time.Now().Format("Mon 02-Jan 15:04")
	ethFiatStr := plugins.FloatSuffixFormatter(0, suff)(ethFiat)
	return fmt.Sprintf("%s - %s", ts, ethFiatStr)
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

	update := <-updates
	if update.Message != nil {
		chatId := update.Message.Chat.ID
		txt := fmt.Sprintf("Your  Chat ID is:\n%d", chatId)
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

func RunBot() {

	log := zaplog.New()
	bot := config.TelegramBot()

	log.Info("Authorized on account", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	ps := plugins.Selected

	for update := range updates {
		if update.Message != nil {
			//chatId := update.Message.Chat.ID
			//logic := getLogic(*ps, chatId)
			var msg *tgbotapi.MessageConfig

			received := update.Message.Text
			if strIs(received, "open", "o", "ope", "refresh", "r", "re", "ref") {
				msg = ps.TelegramFormat(update.Message.Chat.ID, MessageSubject())
			}
			if strIs(received, "help", "h", "he", "hel") {
				msg = newMsg(update.Message.Chat.ID, "help")
			}
			if msg != nil {
				_, err := bot.Send(*msg)
				if err != nil {
					panic(err)
				}
			}

		} else if update.CallbackQuery != nil {
			fmt.Println(update.CallbackQuery.Data)
			pluginId := update.CallbackQuery.Data
			if pluginId == plugins.Void {
				continue
			}
			p := plugins.GetPluginById(pluginId)
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
