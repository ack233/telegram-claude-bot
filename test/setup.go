package test

import (
	"tebot/pkgs/config"
	"tebot/pkgs/initfunc"
	"tebot/pkgs/logtool"
)

func Init() {

	//初始化配置文件
	config.Init()

	//初始化日志系统
	loglevel := config.ViperConfig.GetString("loglevel")
	logfile := config.ViperConfig.GetString("logfile")
	logtool.InitEvent(loglevel, logfile)

	//初始化功能函数
	initfunc.InitFun()

	//初始化数据库
	//dsn := fmt.Sprintf(`%s`, config.ViperConfig.GetString("database"))
	//base.InitDatabase(dsn)
	//base.InitTable()
}
