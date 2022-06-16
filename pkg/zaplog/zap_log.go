// Package zaplog
/**
  @author:kk
  @data:2021/9/4
  @note
**/
package zaplog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"im_app/pkg/config"
	"im_app/pkg/helpler"
)

var (
	Zap   *zap.Logger
	sugar *zap.SugaredLogger
)

func InitZapLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{
		fmt.Sprintf("%slog_%s.zaplog", config.GetString("core.log_address"), helpler.GetNowFormatTodayTime()),
		"stdout",
	}

	// 创建logger实例
	Zap, _ := cfg.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.WarnLevel),
	)
	zap.ReplaceGlobals(Zap) // 替换zap包中全局的logger实例
	sugar = Zap.Sugar()
}

func Fatal(format string, logs ...interface{}) {
	sugar.Fatalf(format, logs...)
}

func Error(format string, logs ...interface{}) {
	sugar.Errorf(format, logs...)
}

func Info(args ...interface{}) {
	sugar.Info(args)
}
