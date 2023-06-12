package handlers

import (
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"tebot/bot/tools"
	"tebot/pkgs/base"
	"tebot/pkgs/config"
	"tebot/pkgs/limitmap"
	"tebot/pkgs/logtool"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	maxChatQueueSize = 20
)

var (
	semChatQueues chan struct{}

	chatQueues = make(map[int64]*ChatQueue)
)

func (b *Bot) Run() {
	semChatQueues = make(chan struct{}, b.Config.MaxConcurrentChats)

	updates := b.Api.ListenForWebhook("/" + b.Api.Token)
	go func() {
		err := http.ListenAndServeTLS(b.Config.ListenAddr, b.Config.CertFile, b.Config.CaFile, nil)
		if err != nil {
			logtool.SugLog.Fatalf("Error starting server: %v", err)
		}
	}()

	var mutex = &sync.Mutex{}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		go func(update tgbotapi.Update) {

			chatID := update.Message.Chat.ID
			mutex.Lock()

			chatQueue, exists := chatQueues[chatID]

			if !exists {
				chatQueue = &ChatQueue{
					Queue:           make(chan *tgbotapi.Update, maxChatQueueSize),
					IsRunning:       0,
					MessageCleanTag: true,
					ChatLimitMap:    limitmap.NewLimitMap(chatID, config.Claudeconfig.MaxConversationids),
					IsGroup:         tools.CheckisGroup(&update),
				}
				chatQueues[chatID] = chatQueue
			}
			mutex.Unlock()

			// If the chatQueue is full, this will block until there's space.
			go func() {
				logtool.SugLog.Debugf("send queue")

				chatQueue.Queue <- &update
			}()

			// Atomically check if IsRunning is 0 and set it to 1
			if atomic.CompareAndSwapInt32(&chatQueue.IsRunning, 0, 1) {
				logtool.SugLog.Debug(atomic.LoadInt32(&chatQueue.IsRunning))

				go b.handleMessage(chatQueue)
			}
			logtool.SugLog.Debug(atomic.LoadInt32(&chatQueue.IsRunning))
		}(update)

	}


}

func (b *Bot) handleMessage(chatQueue *ChatQueue) {
	defer func() {
		logtool.SugLog.Debug("stop")

		atomic.StoreInt32(&chatQueue.IsRunning, 0)

	}() // Update the flag when the goroutine ends.
	timeout := time.After(6 * time.Minute) // Set a timeout of 5 minutes.
	logtool.SugLog.Debug("start1")

	for {
		select {
		case task, ok := <-chatQueue.Queue:
			if !ok {
				return // If the queue is closed, exit the goroutine.
			}
			semChatQueues <- struct{}{}
			// Process the message here...
			//fmt.Println(task.Message.Text)

			logtool.SugLog.Debug("start2")

			b.handleUpdates(MessageQueue{
				ChatQueue: chatQueue,
				Update:    task,
				BotAPI:    b.Api,
			})
			<-semChatQueues
			timeout = time.After(6 * time.Minute) // Reset the timeout after receiving a message.
		case <-timeout:
			return // If there's no new message in 5 minutes, exit the goroutine.
		}
	}
}

func (b *Bot) handleUpdates(msgq MessageQueue) {

	switch {
	case msgq.Message.Text != "":
		//信息入库
		InsertDbMsg(base.Dbc, msgq.Update)

		//添加拦截函数
		msgq.use(msgq.checText()).
			use(msgq.checkIgnoreMsgTag(b.Config.IgnoreMsgTime))
		//设置关键词回复
		if f, ok := msgq.handletext(); ok {
			msgq.run(f)

		} else if !msgq.IsGroup {
			msgq.run(func() {
				PrivateChannelMsg(b.Api, msgq.Update, msgq.ChatLimitMap)
			})

		} else if msgq.Message.ReplyToMessage != nil && msgq.Message.ReplyToMessage.From.UserName == b.Api.Self.UserName {
			msgq.use(msgq.replacGroupText()).
				run(func() {
					ReplyEventMsg(b.Api, msgq.Update, msgq.ChatLimitMap)
				})

		} else if msgq.Message.Entities != nil && strings.Contains(msgq.Message.Text, "@"+b.Api.Self.UserName) {
			msgq.use(msgq.replacGroupText()).
				run(func() {
					MentionEventMsg(b.Api, msgq.Update, msgq.ChatLimitMap)
				})

		}
		//logtool.SugLog.Info(msgq.Message.Text)
	case msgq.Message.NewChatMembers != nil:
		msgq.use(msgq.checkisGroup()).
			run(func() {
				HandleNewMembers(b.Api, msgq.Update)
			})
	case msgq.Message.LeftChatMember != nil:
		msgq.use(msgq.checkisGroup()).
			run(func() {
				HandleLeftMembers(b.Api, msgq.Update)
			})
	}

	//logtool.SugLog.Info("收到的完整更新: %+v\n", update.Message)
}
