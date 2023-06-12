package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendEditmsg(bot *tgbotapi.BotAPI, chatID int64, messageID int, text string) (*tgbotapi.Message, error) {
	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	s, r := bot.Send(edit)
	return &s, r

}

func SendFirstmsg(bot *tgbotapi.BotAPI, chatID int64, messageID int, text string) (*tgbotapi.Message, error) {
	message := tgbotapi.NewMessage(chatID, text)
	message.ReplyToMessageID = messageID // Assuming msg is the Message struct containing the @ event
	s, r := bot.Send(message)
	return &s, r
}
