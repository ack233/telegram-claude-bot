package tools

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CheckIsBot(bot *tgbotapi.BotAPI, update *tgbotapi.Update) (*tgbotapi.User, bool) {

	for _, newUser := range update.Message.NewChatMembers {
		if newUser.IsBot && newUser.UserName != bot.Self.UserName {
			return &newUser, true
		} else {
			return &newUser, false
		}
	}
	return nil, false
}
