package config

type BotconfigStruct struct {
	BotToken           string `yaml:"botToken" validate:"required" `
	WebhookUrl         string `yaml:"webhookUrl" validate:"required"`
	ListenAddr         string `yaml:"listenAddr" validate:"required"`
	CertFile           string `yaml:"certFile" validate:"required"`
	CaFile             string `yaml:"caFile" validate:"required"`
	IgnoreMsgTime      int    `yaml:"ignoreMsgTime"`
	MaxConcurrentChats int    `yaml:"maxConcurrentChats"`
}

type ClaudeconfigStruct struct {
	Api                string `yaml:"api" validate:"required" `
	Token              string `yaml:"token" validate:"required" `
	ChannelID          string `yaml:"channelID" validate:"required" `
	MaxConversationids int    `yaml:"maxConversationids" validate:"required" `
	ReplyModel         string `yaml:"replyModel"`
}

type ForbiddenWordStruct struct {
	Politics []string `yaml:"politics"`
	Attack   []string `yaml:"attack"`
}

type config struct {
	Botconfig     *BotconfigStruct     `yaml:"botconfig"`
	Claudeconfig  *ClaudeconfigStruct  `yaml:"claudeconfig"`
	ForbiddenWord *ForbiddenWordStruct `yaml:"forbiddenWord"`
}
