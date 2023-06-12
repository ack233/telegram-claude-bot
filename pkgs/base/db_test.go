package base

import (
	"fmt"
	"sync"
	"tebot/pkgs/config"
	"tebot/pkgs/initfunc"
	"tebot/pkgs/logtool"
	"testing"
	"time"

	"gorm.io/gorm"
)

var wg sync.WaitGroup

func TestInitDatabase(t *testing.T) {

	tests := []struct {
		name string
		want *gorm.DB
	}{
		{}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//初始化配置文件
			config.Init()

			//初始化日志系统
			loglevel := config.ViperConfig.GetString("loglevel")
			logfile := config.ViperConfig.GetString("logfile")
			logtool.InitEvent(loglevel, logfile)

			//初始化功能函数
			initfunc.InitFun()
			logtool.InitEvent("info", "/home/gowork/tebot/bot.log")
			dsn := fmt.Sprintf(`%s`, config.ViperConfig.GetString("database"))
			InitDatabase(dsn)

			//go PrintCurrentConnections()
			//logtool.Errorerror(err)
			InitTable()

			//删除LimitItem表
			//var tab LimitItem
			//err := Dbc.Exec(fmt.Sprintf(`TRUNCATE TABLE "%s"`, tab.TableName())).Error
			//logtool.Errorerror(err)

			//批量插入数据
			//insertTest


			//查询LimitItem表行数
			var tab2 LimitItem
			var count int64
			Dbc.Model(tab2).Count(&count)
			fmt.Println(count) // 输出行
		})
	}
}

func insertTest(key int, value string) {
	for i := 1; i <= 10000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			logtool.SugLog.Info("*************", key, value)
			err := Dbc.Create(&LimitItem{
				Chatid:         100,
				Contentid:      i,
				Conversationid: fmt.Sprintf("Value %d", i%10),
				//精确到毫秒
				Time: time.Now().Truncate(time.Millisecond),
			}).Error
			logtool.Errorerror(err)
		}(i)
	}
	wg.Wait()
}

func PrintCurrentConnections() {
	ticker := time.NewTicker(3 * time.Millisecond)

	for range ticker.C {
		var count int
		result := Dbc.Raw("SELECT COUNT(*) FROM pg_stat_activity").Scan(&count)
		if result.Error != nil {
			fmt.Println("Error getting connections:", result.Error)
		} else {
			fmt.Println("Current connections:", count)
		}

	}
}
