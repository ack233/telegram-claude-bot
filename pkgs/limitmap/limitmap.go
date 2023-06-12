package limitmap

import (
	"container/list"
	"fmt"
	"sync"
	"tebot/pkgs/base"
	"time"
)

const maxMapSize = 200

type LimitMap struct {
	chatid             int64
	maxMapSize         int
	maxConversationids int
	data               map[int]string
	count              map[string]int
	order              *list.List
	strToKeyList       map[string]*list.List
	sync.Mutex
}

func NewLimitMap(chatID int64, maxConversationids int) *LimitMap {
	lm := &LimitMap{
		chatid:             chatID,
		maxMapSize:         maxMapSize,
		maxConversationids: maxConversationids,
		data:               make(map[int]string),
		count:              make(map[string]int),
		order:              list.New(),
		strToKeyList:       make(map[string]*list.List),
	}
	go func() {

		// Initialize LimitMap with the latest maxSize data from the database
		var items []base.LimitItem
		base.Dbc.Where("chatid = ?", chatID).Order("time desc").Limit(maxConversationids).Find(&items)
		for _, item := range items {
			lm.AddToMemory(item.Contentid, item.Conversationid)
		}
	}()
	return lm
}

func (lm *LimitMap) AddToMemory(key int, value string) {
	lm.Lock()
	defer lm.Unlock()
	if len(lm.data) >= lm.maxMapSize {
		// remove the first item when reach the maxSize
		e := lm.order.Front()
		valueToRemove := lm.data[e.Value.(int)]
		lm.strToKeyList[valueToRemove].Remove(lm.strToKeyList[valueToRemove].Front())
		lm.count[valueToRemove]--
		delete(lm.data, e.Value.(int))
		lm.order.Remove(e)
	}
	lm.data[key] = value
	lm.order.PushBack(key)
	lm.count[value]++
	if _, ok := lm.strToKeyList[value]; !ok {
		lm.strToKeyList[value] = list.New()
	}
	lm.strToKeyList[value].PushBack(key)
}

func (lm *LimitMap) AddDb(key int, value string) {
	base.Dbc.Create(&base.LimitItem{
		Chatid:         lm.chatid,
		Contentid:      key,
		Conversationid: value,
		//精确到毫秒
		Time: time.Now().Truncate(time.Millisecond),
	})
}

func (lm *LimitMap) Add(key int, value string) {
	// Check if value already exists 100 times
	if lm.count[value] >= lm.maxConversationids {
		go func() {
			base.Dbc.Where("chatid = ? AND conversationid = ?", lm.chatid, value).Delete(&base.LimitItem{})
		}()
		// Remove all keys associated with this value
		for e := lm.strToKeyList[value].Front(); e != nil; e = e.Next() {
			keyToRemove := e.Value.(int)
			lm.order.Remove(lm.order.Front())
			delete(lm.data, keyToRemove)

		}
		lm.strToKeyList[value].Init() // Clear the list
		lm.count[value] = 0
		return
	}

	lm.AddToMemory(key, value)

	go lm.AddDb(key, value)
}

func (lm *LimitMap) Display() {
	for e := lm.order.Front(); e != nil; e = e.Next() {
		key := e.Value.(int)
		fmt.Printf("Key: %d, Value: %s\n", key, lm.data[key])
	}
}

func (lm *LimitMap) Get(key int) string {
	// Try to get from memory
	value, exists := lm.data[key]
	if exists {
		return value
	}

	// Try to get from database
	var item base.LimitItem
	result := base.Dbc.First(&item, "chatid = ? and Contentid = ?", lm.chatid, key)
	if result.Error != nil {
		return ""
	}
	return item.Conversationid
}
