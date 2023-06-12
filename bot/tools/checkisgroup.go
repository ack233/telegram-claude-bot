package tools

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CheckisGroup(update *tgbotapi.Update) bool {
	return update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup"
}
