package viper

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// @author: lipper
// @function: NewViper
// @description: 解析配置文件
func NewViper(configFile *string, v any) {
	vp := viper.New()
	vp.SetConfigFile(*configFile)

	if err := vp.ReadInConfig(); err != nil {
		panic("Fatal error config file:" + err.Error())
	}
	if err := vp.Unmarshal(&v); err != nil {
		panic(err)
	}
	vp.OnConfigChange(func(e fsnotify.Event) {
		if err := vp.Unmarshal(&v); err != nil {
			panic(err)
		}
	})
}
