package handlers

import (
	"tebot/pkgs/base"
	"tebot/pkgs/logtool"
	"tebot/bot/tools"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func InsertDbMsg(db *base.Dbr, update *tgbotapi.Update) {
	groupId := update.Message.Chat.ID
	ChatTitle := update.Message.Chat.Title
	FromID := update.Message.From.ID
	FromFirstName := update.Message.From.FirstName
	FromLastName := update.Message.From.LastName
	MessageText := update.Message.Text
	MessageID := update.Message.MessageID
	MessageTime := update.Message.Time()

	go func() {

		logtool.SugLog.Infof(
			"会话: %s, 会话id: %d, 用户id: %d, 用户姓名: %s%s, 消息体: %s, 日期: %s",
			ChatTitle,
			groupId,
			FromID,
			FromFirstName,
			FromLastName,
			MessageText,
			tools.FormatTime(MessageTime),
		)

		text := base.Communication_context{
			Userid:    FromID,
			Chatname:  ChatTitle,
			Chatid:    groupId,
			Username:  FromFirstName + FromLastName,
			Contentid: MessageID,
			Content:   MessageText,
			Time:      MessageTime,
		}
		db.Create(&text)
	}()
}
