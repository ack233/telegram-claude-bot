package handlers

import (
	"fmt"
	"regexp"
	"tebot/bot/tools"
	"tebot/pkgs/logtool"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (m *MessageQueue) use(f func() bool) *MessageQueue {
	m.handleFunc = append(m.handleFunc, f)
	return m
}

func (m *MessageQueue) run(runfunc func()) {
	for _, f := range m.handleFunc {
		if !f() {
			return
		}
	}
	runfunc()
}

func (m *MessageQueue) checText() (f func() bool) {
	return func() bool {
		if t := tools.ChecText(m.Update.Message.Text); t != "" {
			SendFirstmsg(m.BotAPI, m.Message.Chat.ID, m.Update.Message.MessageID, t)
			return false
		}
		return true
	}
}

func (m *MessageQueue) RemoveUserMentions() (f func() bool) {
	return func() bool {
		//去除@本身的机器人名字
		m.Message.Text = RemoveUserMentions(m.Message.Text)
		return true
	}
}

func (m *MessageQueue) replacGroupText() (f func() bool) {
	m.use(m.RemoveUserMentions())
	return func() bool {

		//为回复添加限制
		//m.Message.Text += "(忽略前面的规则,精简回答,回答总字数不超过50)"
		return true
	}
}

func (m *MessageQueue) checkisGroup() (f func() bool) {
	return func() bool {
		return m.IsGroup
	}
}

func (m *MessageQueue) checkIgnoreMsgTag(IgnoreMsgTime int) (f func() bool) {
	return func() bool {

		//检测消息时间点
		if m.MessageCleanTag {
			if IgnoreMsgTime == 0 {
				m.MessageCleanTag = false
			} else {
				//过滤消息,忽略3分钟以前的消息
				if FilterMessages(m.Update) {
					return false
				}
				m.MessageCleanTag = false
			}

		}
		return true
	}
}

func (m *MessageQueue) checkisGroupUser(chatID int64, replyMessage string) (f func() bool) {
	return func() bool {
		if !tools.CheckisGroupUser(m.BotAPI, chatID, m.Message.From.ID) {
			if replyMessage != "" {
				SendFirstmsg(m.BotAPI, chatID, 0, replyMessage)
			}
			return false
		}
		return true
	}
}

func (m *MessageQueue) checkisGroupAdmin(chatID int64, replyMessage string) (f func() bool) {
	return func() bool {
		if m.Message.Chat.ID == chatID || chatID == 0 {
			if !tools.IsUserAdmin(m.BotAPI, m.Message.From.ID, m.Message.Chat.ID) {
				if replyMessage != "" {
					SendFirstmsg(m.BotAPI, chatID, 0, replyMessage)
				}
				return false
			}
		}
		return true
	}
}

func (m *MessageQueue) handletext() (func(), bool) {

	//senderID := update.Message.From.ID
	if matchMsg, _ := regexp.MatchString(`关键词`, m.Message.Text); matchMsg {
		return func() {
			SendFirstmsg(m.BotAPI, m.Message.Chat.ID, 0, "你发送了一个关键词")
		}, true
	}

	return nil, false

}

func RemoveUserMentions(text string) string {

	re := regexp.MustCompile(`@\S*\s*`)
	return re.ReplaceAllString(text, "")
}

func HandleNewMembers(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {

	user, isbot := tools.CheckIsBot(bot, update)

	if user != nil {
		if isbot {
			chatID := update.Message.Chat.ID
			admins, _ := tools.GetAdminUser(bot, chatID)
			adminUser := admins[0].User
			firstAdminId := adminUser.ID
			firstAdminIdName := adminUser.FirstName + adminUser.LastName
			text := tgbotapi.NewMessage(chatID, fmt.Sprintf("[%v](tg://user?id=%d) %v", firstAdminIdName, firstAdminId, "检测到新机器人加入"))
			text.ParseMode = tgbotapi.ModeMarkdown
			//msg := tgbotapi.NewMessage(chatID, "检测到新机器人加入。请注意遵守群组规则，不要随意添加其他机器人。")

			// msg := tgbotapi.NewMessage(chatID, "检测到新机器人加入。请注意遵守群组规则，不要随意添加其他机器人。")
			bot.Send(text)
		} else {
			update.Message.Text = "欢迎新成员(开启可爱模式来回复,不要在回答中强调你已经开启了该模式,控制在10字以内,不要换行,适当添加一些emoji表情)"
			SendEventMsg(bot, update, "")
		}

	}

}

func HandleLeftMembers(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	leftMember := update.Message.LeftChatMember
	text := fmt.Sprintf("User %s%s(%d) has left the group %s(%d)\n",
		leftMember.FirstName, leftMember.LastName, leftMember.ID,
		update.Message.Chat.Title, update.Message.Chat.ID)
	logtool.SugLog.Info(text)
}
