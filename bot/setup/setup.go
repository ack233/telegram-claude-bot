package setup

import (
	"tebot/bot/handlers"
	"tebot/pkgs/config"
	"tebot/pkgs/logtool"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Initbot(cf config.BotconfigStruct) (handlers.Bot, error) {

	bot, err := tgbotapi.NewBotAPI(cf.BotToken)
	logtool.Fatalerror(err)

	//bot.Debug = true

	logtool.SugLog.Infof("Authorized on account %s", bot.Self.UserName)
	wh, _ := tgbotapi.NewWebhookWithCert(cf.WebhookUrl+bot.Token, nil)

	_, err = bot.Request(wh)
	logtool.Fatalerror(err)

	info, err := bot.GetWebhookInfo()
	logtool.Fatalerror(err)

	if info.LastErrorDate != 0 {
		logtool.SugLog.Warnf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return handlers.Bot{Api: bot, Config: cf}, err
}

//
//func (b *Bot) InitApplication() {
//	// Initialize scheduled tasks
//	tasks.InitTasks(b.Api)
//}
