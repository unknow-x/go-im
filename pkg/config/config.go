// Package config
/**
  @author:kk
  @data:2021/6/18
  @note
**/
package config

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"os"
)

var Viper *viper.Viper

type StrMap map[string]interface{}

func init() {
	// 读取默认配置
	Viper = viper.New()
	Viper.SetConfigName("im")
	Viper.SetConfigType("yaml")
	Viper.AddConfigPath(os.Getenv("GO_IM_ENV_PATH"))
	Viper.AddConfigPath("config")
	err := Viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	configs := Viper.AllSettings()
	// 将default中的配置全部以默认配置写入
	for k, v := range configs {
		viper.SetDefault(k, v)
	}

	// 读取.env获取环境配置
	v2 := viper.New()
	v2.SetConfigName(".env")
	v2.SetConfigType("env")
	if err != nil {
		v2.Set("env", "production")
	}

	env := v2.GetString("env")
	// 根据配置的env读取相应的配置信息
	if env != "" {
		viper.SetConfigName(env)
		viper.SetConfigType("yaml")
		viper.AddConfigPath("config")
		err = viper.ReadInConfig()
		if err != nil {
			return
		}
	}
	// 设置环境变量前缀，用以区分 Go 的系统环境变量
	Viper.SetEnvPrefix("appenv")
	// Viper.Get() 时，优先读取环境变量
	Viper.AutomaticEnv()

	fmt.Println(Viper.AllSettings())
	os.Exit(0)
}

// Env 读取环境变量，支持默认值
func Env(envName string, defaultValue ...interface{}) interface{} {
	if len(defaultValue) > 0 {
		return Get(envName, defaultValue[0])
	}
	return Get(envName)
}

// Add 新增配置项
func Add(name string, configuration map[string]interface{}) {
	Viper.Set(name, configuration)
}

// Get 获取配置项，允许使用点式获取，如：core.name
func Get(path string, defaultValue ...interface{}) interface{} {
	// 不存在的情况
	if !Viper.IsSet(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return Viper.Get(path)
}

// GetString 获取 String 类型的配置信息
func GetString(path string, defaultValue ...interface{}) string {
	return cast.ToString(Get(path, defaultValue...))
}

// GetInt 获取 Int 类型的配置信息
func GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(Get(path, defaultValue...))
}

// GetInt64 获取 Int64 类型的配置信息
func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(Get(path, defaultValue...))
}

// GetInt32 获取 Int64 类型的配置信息
func GetInt32(path string, defaultValue ...interface{}) int32 {
	return cast.ToInt32(Get(path, defaultValue...))
}

// GetUint 获取 Uint 类型的配置信息
func GetUint(path string, defaultValue ...interface{}) uint {
	return cast.ToUint(Get(path, defaultValue...))
}

// GetBool 获取 Bool 类型的配置信息
func GetBool(path string, defaultValue ...interface{}) bool {
	return cast.ToBool(Get(path, defaultValue...))

}
