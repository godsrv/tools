package cache

import (
	"errors"
	"time"
)

var (
	// ErrEmptyHost is an error that indicates no redis host is set.
	ErrEmptyHost = errors.New("empty redis host")
	// ErrEmptyType is an error that indicates no redis type is set.
	ErrEmptyType = errors.New("empty redis type")
	// ErrEmptyKey is an error that indicates no redis key is set.
	ErrEmptyKey = errors.New("empty redis key")
	// ErrEmptyPort is an error that no port is set.
	ErrEmptyPort = errors.New("empty redis port")
)

type (
	// A RedisConf is a redis config.
	RedisConf struct {
		Host        string        `json:"host"`
		DB          int           `json:"db,default=1"`
		Password    string        `json:"password,optional"`
		Port        int           `json:"port,optional"`
		Tls         bool          `json:"tls,optional"`
		PingTimeout time.Duration `json:"ping_timeout,default=1s"`
	}
)

// Validate validates the RedisConf.
func (rc RedisConf) Validate() error {

	if len(rc.Host) == 0 {
		return ErrEmptyHost
	}

	if rc.Port == 0 {
		return ErrEmptyPort
	}

	return nil
}
