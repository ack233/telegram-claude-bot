package handlers

import (
	"strings"
	"tebot/claude"
	"tebot/pkgs/logtool"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendEventMsgWithCallback(handlerFunc func() CallbackFunc, isGroup bool, limitMap Limitmap) {
	callbackFunc := handlerFunc()
	callbackFunc(isGroup, limitMap)
}

func MentionEventMsg(bot *tgbotapi.BotAPI, update *tgbotapi.Update, limitMap Limitmap) {
	for _, entity := range update.Message.Entities {
		if entity.Type == "mention" {
			logtool.SugLog.Infof("%s收到艾特事件...", update.Message.Chat.Title)

			sendEventMsgWithCallback(
				func() CallbackFunc {
					return SendEventMsg(bot, update, "")
				}, true, limitMap,
			)
			return
		}
	}
}

func ReplyEventMsg(bot *tgbotapi.BotAPI, update *tgbotapi.Update, limitMap Limitmap) {
	logtool.SugLog.Infof("%s收到引用事件...", update.Message.Chat.Title)

	conversationID := limitMap.Get(update.Message.ReplyToMessage.MessageID)

	sendEventMsgWithCallback(
		func() CallbackFunc {
			return SendEventMsg(bot, update, conversationID)
		}, true, limitMap,
	)

}

func PrivateChannelMsg(bot *tgbotapi.BotAPI, update *tgbotapi.Update, limitMap Limitmap) {
	chatid := update.Message.Chat.ID
	logtool.SugLog.Infof("用户%s%s发起聊天...", update.Message.From.FirstName, update.Message.From.LastName)

	conversationID := limitMap.Get(int(chatid))

	sendEventMsgWithCallback(
		func() CallbackFunc {
			return SendEventMsg(bot, update, conversationID)
		}, false, limitMap,
	)

}

func checkRetryErr(err error) bool {
	return strings.Contains(err.Error(), "Too Many Requests: retry after")
}

func SendEventMsg(bot *tgbotapi.BotAPI, update *tgbotapi.Update, conversationID string) (callbackfunc func(bool, Limitmap)) {
	var lastMessage *tgbotapi.Message
	ChatID := update.Message.Chat.ID
	MessageID := update.Message.MessageID

	callbackfunc = func(b bool, l Limitmap) {
		logtool.SugLog.Warn("回调函数异常")
	}

	var err error
	var eventdatas eventdataStruct
	text := update.Message.Text

	bot.Send(tgbotapi.NewChatAction(ChatID, tgbotapi.ChatTyping))

	ticker := time.NewTicker(time.Millisecond * 1600) // Set your desired interval
	defer ticker.Stop()

	eventChan := claude.GetMessFromClaude(text, conversationID)
	for {
		eventClaude, ok := <-eventChan

		if ok {
			if err = eventClaude.DecodeErr; err != nil {
				logtool.SugLog.Error(err)
				return
			}
			neweventdate := eventClaude.Message.Content.Parts[0]
			if eventdatas.eventdata == neweventdate {
				continue
			}
			//eventdata = regexp.
			//	MustCompile(`(?i)(claude)|(Anthropic)`).
			//	ReplaceAllString(eventdata, "小烧杯")
			eventdatas = eventdataStruct{
				eventdata:      neweventdate,
				isSend:         false,
				ConversationID: eventClaude.ConversationID,
			}
		}
		logtool.SugLog.Debug(eventdatas.isSend, ok)

		select {
		case <-ticker.C:
			logtool.SugLog.Debug("send1")

			if eventdatas.eventdata == "" {
				if !ok {
					return
				}
				break
			}
			//eventdatas.afterReplacedata = regexp.
			//	MustCompile(`(?i)(claude)|(Anthropic)|(克劳德)`).
			//	ReplaceAllString(eventdatas.eventdata, "your ai")
			eventdatas.afterReplacedata = eventdatas.eventdata
			if !eventdatas.isSend {
				if lastMessage != nil {
					_, err = sendEditmsg(bot, ChatID, lastMessage.MessageID, eventdatas.afterReplacedata)
					if err != nil {
						if !ok {
							if strings.Contains(err.Error(), "Bad Request: message is not modified") {
								logtool.SugLog.Warn("最后一次消息未发生修改")
								err = nil
							} else if !checkRetryErr(err) {
								for {
									var lastSendmsg *tgbotapi.Message
									lastSendmsg, err = SendFirstmsg(bot, ChatID, MessageID, eventdatas.afterReplacedata)
									if err != nil {
										if !checkRetryErr(err) {
											logtool.SugLog.Infof("事件处理失败: %v", err.Error())
											return
										}
									} else {
										if lastMessage != nil {
											msgToDelete := tgbotapi.NewDeleteMessage(ChatID, lastMessage.MessageID)
											bot.Send(msgToDelete)
											logtool.SugLog.Warn("删除旧的不完整消息")
										}
										lastMessage = lastSendmsg
										break
									}
								}
							}
						}
					}
				} else {
					lastMessage, err = SendFirstmsg(bot, ChatID, MessageID, eventdatas.afterReplacedata)
				}
			}

			if err == nil {
				if ok {
					eventdatas.isSend = true
				} else {
					logtool.SugLog.Info("事件处理完成")
					callbackfunc = func(isGroup bool, limitmap Limitmap) {
						if isGroup {
							limitmap.Add(lastMessage.MessageID, eventdatas.ConversationID)
						} else {
							limitmap.Add(int(ChatID), eventdatas.ConversationID)
						}
					}
					return callbackfunc
				}
			} else {
				logtool.SugLog.Error(err)
			}

		default:
			logtool.SugLog.Debug("default")

			if !ok {
				logtool.SugLog.Debug("default2")

				ticker.Stop()
				ticker = time.NewTicker(time.Millisecond)
				time.Sleep(2 * time.Millisecond)
			}
		}

	}
}
