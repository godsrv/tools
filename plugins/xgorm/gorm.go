package xgorm

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	Client *gorm.DB
)

// @author: lipper
// @function: NewGorm
// @description: 实例orm
// @return: *gorm.DB
func NewGorm(conf DataBaseConf) *gorm.DB {
	var err error
	// 获取dsn
	Client, err = conf.GetDB()
	if err != nil {
		panic(err)
	}

	if conf.Debug {
		Client = Client.Debug()
	}

	sqlDB, err := Client.DB()
	if err != nil {
		logrus.Panicf("get database/sql db err: %v", err)
	}

	if conf.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(conf.MaxIdleConn)
	}

	if conf.MaxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(conf.MaxOpenConn)
	}

	if conf.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(conf.ConnMaxIdleTime)
	}

	if conf.ConnMaxLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(conf.ConnMaxLifeTime)
	}

	return Client
}
