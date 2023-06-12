package base

import (
	"errors"
	"net"
	"strings"
	"tebot/pkgs/logtool"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Dbc *Dbr

const maxAttempts = 5
const delayBetweenAttemptsSeconds = 5

type Dbr struct {
	*gorm.DB
}

func (d *Dbr) Create(value interface{}) (db *gorm.DB) {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db = d.DB.Create(value)
		if db.Error == nil || errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return
		}
		time.Sleep(delayBetweenAttemptsSeconds * time.Second)
		logtool.SugLog.Errorf("db create failed  %v, attempt %d of %d. Retrying in %vs...\n", db.Error, attempt, maxAttempts, delayBetweenAttemptsSeconds)
	}
	return
}

type retryPlugin struct {
	maxAttempts          int
	delayBetweenAttempts time.Duration
}

func (p *retryPlugin) Name() string {
	return "retryPlugin"
}

func isConnectionError(err error) (string, bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		//08000: 连接异常
		//08001: SQL 客户端无法建立连接
		//08003: 连接不存在
		//08004: SQL 服务器拒绝连接
		//08006: 连接故障
		//08007: 在事务结束后连接被丢弃
		//53300: too many clients already
		//57P01: 数据库服务正在关闭或者已经关闭
		connectionErrorCodes := []string{"08000", "08001", "08003", "08004", "08006", "08007", "53300", "57P01"}
		for _, code := range connectionErrorCodes {
			if pgErr.Code == code {
				return code, true
			}
		}
	} else if netErr, ok := err.(net.Error); ok {
		// handle network error
		return netErr.Error(), true
	} else if strings.Contains(err.Error(), "connection refused") {
		// handle connection refused error
		return "connection_refused", true
	} else if strings.Contains(err.Error(), "bad connection") {
		return "", false
	}
	return "", false
}

func (p *retryPlugin) Initialize(db *gorm.DB) (err error) {
	// Wrap the original Create function with retry logic
	originalCreate := db.Callback().Create().Get("gorm:create")
	db.Callback().Create().Replace("gorm:create", func(scope *gorm.DB) {
		for attempt := 1; attempt <= p.maxAttempts; attempt++ {

			originalCreate(scope)
			if scope.Error == nil || errors.Is(scope.Error, gorm.ErrRecordNotFound) {
				break
			}
			if code, b := isConnectionError(scope.Error); b {
				logtool.SugLog.Errorf("Insert failed state is %s, attempt %d of %d. Retrying in %v...\n", code, attempt, p.maxAttempts, p.delayBetweenAttempts)
				time.Sleep(p.delayBetweenAttempts)
				//scope = scope.Session(&gorm.Session{NewDB: true})
				scope.Error = nil
			} else {
				logtool.SugLog.Errorf("Insert failed, %v", scope.Error)
				break
			}
		}
	})
	return nil
}

func InitDatabase(dsn string) *Dbr {
	for {
		var err error
		dbr, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logtool.GromLog})
		if err != nil {
			logtool.SugLog.Error("数据库连接失败")
			time.Sleep(time.Second * 5)
			continue
		}

		retryPluginInstance := &retryPlugin{
			maxAttempts:          maxAttempts,
			delayBetweenAttempts: delayBetweenAttemptsSeconds * time.Second,
		}
		dbr.Use(retryPluginInstance)
		sqlDB, err := dbr.DB()
		logtool.Fatalerror(err)

		// SetMaxIdleConns 设置连接池中的空闲连接的最大数量
		sqlDB.SetMaxIdleConns(10)

		// SetMaxOpenConns 设置数据库的最大连接数量
		sqlDB.SetMaxOpenConns(105)

		// SetConnMaxLifetime 设置连接的最大可复用时间
		sqlDB.SetConnMaxLifetime(time.Hour)
		logtool.SugLog.Info("数据库建立连接")
		Dbc = &Dbr{dbr}
		return Dbc
	}
}

func PingDatabase(db *gorm.DB) bool {
	if mdb, err := db.DB(); err != nil {
		if err := mdb.Ping(); err != nil {
			logtool.SugLog.Error("数据连接不可用")
			return false
		}
	}
	return true
}

func CheckTableIfNotExists(db *gorm.DB, dst interface{}) bool {
	Migrator := db.Migrator()
	return Migrator.HasTable(dst)
}
