package models

import (
	"fmt"
	"gen/config"
	"gen/log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	conns map[string]*gorm.DB
)

type SQLService struct {
	Cfg *config.AppConfig `inject:""`

	conns map[string]*gorm.DB
}

func DBS(dbName ...string) *gorm.DB {
	if len(dbName) > 0 {
		if conn, ok := conns[dbName[0]]; ok {
			return conn
		}
	}
	return DB
}

// InitDB 初始化数据库
func InitDB(cfg *config.AppConfig) error {
	conns = make(map[string]*gorm.DB)
	for _, v := range cfg.DBConfig {
		log.Logger.Info(v.Dialect)
		conn, err := openConn(v.Dialect, v.Dsn, v.MaxIdleConn, v.MaxOpenConn)
		if err != nil {
			return fmt.Errorf("open connection failed, error: %s", err.Error())
		}
		conns[v.Name] = conn
		if v.Name == "default" {
			DB = conn
		}
	}
	return nil
}

func openConn(dialect, dsn string, idle, open int) (*gorm.DB, error) {
	newLogger := logger.New(Writer{}, logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true})
	var openDB *gorm.DB
	var err error
	switch dialect {
	case "mysql":
		openDB, err = gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{Logger: newLogger})
		if err != nil {
			return nil, err
		}
		db, err := openDB.DB()
		if err != nil {
			return nil, err
		}
		db.SetMaxIdleConns(idle)
		db.SetMaxOpenConns(open)
	case "sqlite3":
		openDB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: newLogger})
		if err != nil {
			return nil, err
		}
		db, err := openDB.DB()
		if err != nil {
			return nil, err
		}
		db.SetMaxIdleConns(idle)
		db.SetMaxOpenConns(open)
		// db.SingularTable(true)

	}
	return openDB, nil
}

// ConnectDbSqlite3 连接sqlite3数据库
// func ConnectDbSqlite3(host string) *gorm.DB {
// 	dns := fmt.Sprintf(
// 		"%s",
// 		host,
// 	)
// 	db, err := gorm.Open("sqlite3", dns)
// 	if err != nil {
// 		log.Fatalf("models.ConnectDbSqlite3 err: %v", err)
// 	}
// 	db.SingularTable(true)
// 	return db
// }

// Writer 记录SQL日志
type Writer struct{}

func (w Writer) Printf(format string, args ...interface{}) {
	log.Logger.Debug(fmt.Sprintf(format, args...))
}
