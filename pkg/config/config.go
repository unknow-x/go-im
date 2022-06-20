// Package config
/**
  @author:kk
  @data:2021/6/18
  @note
**/
package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

func Initialize(conf interface{}) {
	configFile := "im"
	configType := "yaml"

	// 读取默认配置
	v1 := viper.New()
	v1.SetConfigName(configFile)
	v1.SetConfigType(configType)
	v1.AddConfigPath(os.Getenv("GO_IM_CONFIG_PATH"))
	v1.AddConfigPath("config")
	err := v1.ReadInConfig()
	if err != nil {
		zap.L().Fatal("load default config error", zap.Error(err))
	}
	configs := v1.AllSettings()
	// 将默认配置全部以默认值写入
	for k, v := range configs {
		viper.SetDefault(k, v)
	}

	// 设置环境变量前缀，用以区分 Go 的系统环境变量
	viper.SetEnvPrefix("goimenv")
	viper.AutomaticEnv()

	// 读取环境覆盖配置
	mode := viper.GetString("mode")
	if mode != "" {
		configFile += "_" + mode
	}
	viper.SetConfigName(configFile)
	viper.SetConfigType(configType)
	viper.AddConfigPath(".")
	viper.AddConfigPath(os.Getenv("GO_IM_CONFIG_PATH"))
	viper.AddConfigPath("config")
	hDir, err := homedir.Dir()
	if err != nil {
		zap.L().Fatal("fetch user homedir error", zap.Error(err))
	}
	viper.AddConfigPath(hDir)
	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			zap.L().Fatal("load env config error", zap.Error(err))
		}
	}

	// 配置写入结构体
	if conf != nil {
		if err := viper.Unmarshal(&conf); err != nil {
			zap.L().Fatal("unmarshal conf failed", zap.Error(err))
		}

		// 监控配置文件变化
		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			fmt.Println("夭寿啦~配置文件被人修改啦...")
			if err := viper.Unmarshal(&conf); err != nil {
				zap.L().Panic("re unmarshal conf failed", zap.Error(err))
			}
		})
	}
}
