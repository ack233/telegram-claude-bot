package handlers

import (
	"tebot/pkgs/config"
	"tebot/pkgs/limitmap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	Api    *tgbotapi.BotAPI
	Config config.BotconfigStruct
}

type ChatQueue struct {
	Queue           chan *tgbotapi.Update
	IsRunning       int32 // 0 means not running, 1 means running.
	MessageCleanTag bool
	ChatLimitMap    *limitmap.LimitMap
	IsGroup         bool
}

type MessageQueue struct {
	*ChatQueue
	*tgbotapi.Update
	*tgbotapi.BotAPI
	handleFunc []func() bool
}

type Limitmap interface {
	Add(int, string)
	Get(int) string
}

type eventdataStruct struct {
	eventdata        string
	afterReplacedata string
	isSend           bool
	ConversationID   string
}

type CallbackFunc func(bool, Limitmap)
