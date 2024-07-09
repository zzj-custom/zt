package bootstrapper

import (
	"github.com/sirupsen/logrus"
	"os"
	"zt/backend/cmd/bootstrapper/config"
	"zt/backend/pkg/iMysql"
	"zt/backend/pkg/iRedis"
)

func Release() {
	releaseResource()
}

// 初始化资源
func initResource() {
	// 初始化数据库
	initDatabase()
	// 初始化redis
	initRedis()

	// 重置log输出到标准输出
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   false,
	})
}

// 释放资源
func releaseResource() {
	// 释放redis
	iRedis.Release()
	// 释放数据库
	iMysql.Release()
}

// 初始化数据库
func initDatabase() {
	c := config.Get()
	_, err := iMysql.Client(c.Database)
	if err != nil {
		panic(err)
	}
}

// 初始化redis
func initRedis() {
	c := config.Get()
	err := iRedis.InitMultiPools(c.Redis)
	if err != nil {
		panic(err)
	}
}
