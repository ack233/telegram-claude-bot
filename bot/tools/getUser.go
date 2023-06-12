package tools

import (
	"tebot/pkgs/logtool"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetUser(bot *tgbotapi.BotAPI, chatID int64, userID int64) (tgbotapi.ChatMember, error) {
	return bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	})

}

func CheckisGroupUser(bot *tgbotapi.BotAPI, chatID int64, userID int64) bool {
	u, err := GetUser(bot, chatID, userID)
	if err == nil {
		return u.Status == "member" || u.Status == "administrator" || u.Status == "creator"
	}
	logtool.Errorerror(err)
	return false
}
