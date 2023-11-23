package xgorm

import (
	"gorm.io/gorm"
)

// @author: lipper
// @function: NewGorm
// @description: 实例orm
// @return: *gorm.DB
func NewGorm(conf DataBaseConf) *gorm.DB {
	// 获取dsn
	db, err := conf.GetDB()
	if err != nil {
		panic(err)
	}
	db.InstanceSet("gorm:table_options", "ENGINE=InnoDB")
	sqlDB, _ := db.DB()

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

	return db
}
