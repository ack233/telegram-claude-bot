package handlers

import (
	"tebot/bot/tools"
	"tebot/pkgs/logtool"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func FilterMessages(update *tgbotapi.Update) bool {
	messageTime := update.Message.Time()
	threshold := time.Now().Add(-3 * time.Minute)
	if messageTime.Before(threshold) {
		// Ignore the message
		logtool.SugLog.Warnf("Ignoring message: time is %s", tools.FormatTime(messageTime))
		return true
	}
	return false

}
