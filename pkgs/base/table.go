package base

import (
	"time"
)

type Socialquotes struct {
	ID      uint   `gorm:"primaryKey"`
	Content string `gorm:"column:Content"`
}

func (*Socialquotes) TableName() string {
	return "socialquotes"
}

type Communication_context struct {
	ID        uint `gorm:"primaryKey"`
	Userid    int64
	Username  string
	Chatid    int64
	Chatname  string
	Contentid int
	Content   string
	Time      time.Time
}

func (*Communication_context) TableName() string {
	return "communication_context"
}

type LimitItem struct {
	Chatid         int64
	Contentid      int
	Conversationid string
	Time           time.Time
}

func (*LimitItem) TableName() string {
	return "limitItem"
}

func createTableIfNotExists(tables ...interface{}) error {
	err := Dbc.AutoMigrate(tables...)
	return err
}

func InitTable() {
	createTableIfNotExists(
		&Communication_context{},
		&Socialquotes{},
		&LimitItem{},
	)
}
