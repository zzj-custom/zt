package iMysql

import (
	"errors"
	"fmt"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"net/url"
	"sync"
	"time"
	"zt/backend/pkg/iLogger"
)

var (
	dbOnce    sync.Once
	clients   map[string]*gorm.DB
	err       error
	logWriter iLogger.Writer = iLogger.DefaultLogger
)

func WithLogger(l iLogger.Writer) {
	logWriter = l
}

func Conns() map[string]*gorm.DB {
	return clients
}

func Conn(key string) (*gorm.DB, error) {
	if clients == nil {
		return nil, errors.New("db connections not initialized")
	}
	conn, ok := clients[key]
	if !ok {
		return nil, errors.New("no such db connection defined")
	}
	return conn, nil
}

func NewClient(dbs map[string]*Database) (map[string]*gorm.DB, error) {
	conns := map[string]*gorm.DB{}
	for k, v := range dbs {
		var dsn = buildDSN(v)
		dialector := gormMySQL.New(gormMySQL.Config{
			DSN:                       dsn,
			SkipInitializeWithVersion: false,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			DontSupportForShareClause: true,
		})
		cfg := &gorm.Config{
			DisableAutomaticPing: false,
		}
		if v.TablePrefix != "" {
			cfg.NamingStrategy = schema.NamingStrategy{TablePrefix: v.TablePrefix, SingularTable: v.Singular}
		}

		config := logger.Config{
			SlowThreshold:             0,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  0,
		}
		if v.UseLog {
			var logLvl logger.LogLevel
			if v.LogLevel > 0 && v.LogLevel <= 4 {
				logLvl = logger.LogLevel(v.LogLevel)
			} else {
				logLvl = logger.Error
			}

			if v.Slowlog.String() != "" {
				config.SlowThreshold = v.Slowlog.Duration
			}
			if logLvl > 0 {
				config.LogLevel = logLvl
			}
		}
		l := logger.New(logWriter, config)
		cfg.Logger = l

		conn, err := gorm.Open(dialector, cfg)
		if err != nil {
			panic(fmt.Sprintf("Failed to create connection for db: %s, error: %+v", k, err))
		}
		sqlDB, err := conn.DB()
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to db: %s, error: %+v", k, err))
		}
		sqlDB.SetMaxIdleConns(v.MaxIdleConn)
		sqlDB.SetMaxOpenConns(v.MaxOpenConn)
		sqlDB.SetConnMaxIdleTime(time.Duration(v.ConnMaxIdleTime) * time.Second)
		sqlDB.SetConnMaxLifetime(time.Duration(v.ConnMaxLifetime) * time.Second)

		conns[k] = conn
	}

	return conns, nil
}

func Client(dbs map[string]*Database) (map[string]*gorm.DB, error) {
	dbOnce.Do(func() {
		clients, err = NewClient(dbs)
	})
	return clients, err
}

// Release 释放数据库连接池
func Release() {
	for _, conn := range clients {
		db, err := conn.DB()
		if err != nil {
			continue
		}
		_ = db.Close()
	}
}

func buildDSN(v *Database) string {
	var dsn string
	if v.DSN != "" {
		dsn = v.DSN
	} else {
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?parseTime=true",
			url.QueryEscape(v.Username),
			url.QueryEscape(v.Password),
			v.Host,
			v.Database,
		)
	}
	return dsn
}
