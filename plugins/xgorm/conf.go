package xgorm

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataBaseConf struct {
	Driver          string        `json:"driver"`
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	DBName          string        `json:"dbName"`
	Username        string        `json:"username"`
	Password        string        `json:"password"`
	MaxIdleConn     int           `json:",optional,maxIdleConn"`
	MaxOpenConn     int           `json:",optional,maxOpenConn"`
	ConnMaxIdleTime time.Duration `json:",optional,connMaxIdleTime"`
	ConnMaxLifeTime time.Duration `json:",optional,connMaxLifeTime"`
	Parameters      string        `json:",optional,parameters"`
	SSLMode         string        `json:",optional,sslMode"`
}

func (c DataBaseConf) GetDSN() string {
	switch c.Driver {
	case "mysql":
		return c.MysqlDSN()
	case "postgres":
		return c.PostgresDSN()
	default:
		return "mysql"
	}
}

// get db by driver
func (c DataBaseConf) GetDB() (*gorm.DB, error) {
	dsn := c.GetDSN()
	switch c.Driver {
	case "mysql":
		return newMySQLDialect(dsn)
	case "postgres":
		return newPostgresDialect(dsn)
	default:
		return newMySQLDialect(dsn)
	}
}

// MysqlDSN returns mysql DSN.
func (c DataBaseConf) MysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", c.Username, c.Password, c.Host, c.Port, c.DBName, c.Parameters)
}

// PostgresDSN returns Postgres DSN.
func (c DataBaseConf) PostgresDSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", c.Username, c.Password, c.Host, c.Port, c.DBName,
		c.SSLMode)
}

// newMySQLDialect build mysql dialect
func newMySQLDialect(dsn string) (*gorm.DB, error) {
	mysqlConfig := mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 191,
	}
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	return db, err
}

// newPostgresDialect build postgres dialect
func newPostgresDialect(dsn string) (*gorm.DB, error) {
	pgsqlConfig := postgres.Config{
		DSN:                  dsn, // DSN data source name
		PreferSimpleProtocol: false,
	}
	db, err := gorm.Open(postgres.New(pgsqlConfig), &gorm.Config{})
	return db, err
}
