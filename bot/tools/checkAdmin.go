package tools

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetAdminUser(bot *tgbotapi.BotAPI, chatID int64) ([]tgbotapi.ChatMember, error) {
	return bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: chatID,
		},
	})
}

func IsUserAdmin(bot *tgbotapi.BotAPI, userID int64, chatID int64) bool {
	admins, err := GetAdminUser(bot, chatID)
	//logtool.SugLog.Info(admins)
	if err != nil {
		return false
	}

	for _, admin := range admins {
		if admin.User.ID == userID {
			return true
		}
	}

	return false
}
