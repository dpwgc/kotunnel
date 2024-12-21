package base

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"time"
)

var closePrint = false
var Logger *slog.Logger

func InitLog(config LogConfig) {

	if len(config.Path) <= 0 {
		config.Path = "./logs"
	}
	if config.Size <= 0 {
		config.Size = 1
	}
	if config.Age <= 0 {
		config.Age = 30
	}
	if config.Backups <= 0 {
		config.Backups = 1000
	}

	closePrint = config.ClosePrint
	r := &lumberjack.Logger{
		Filename:   config.Path + "/runtime.log",
		LocalTime:  true,
		MaxSize:    config.Size,
		MaxAge:     config.Age,
		MaxBackups: config.Backups,
		Compress:   false,
	}
	Logger = slog.New(slog.NewTextHandler(r, &slog.HandlerOptions{
		AddSource: true, // 输出日志语句的位置信息
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey { // 格式化 key 为 "time" 的属性值
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.DateTime))
				}
			}
			return a
		},
	}))
}

const (
	Red    = 31
	Yellow = 33
	Blue   = 36
	Green  = 32
)

func Error(s string) {
	log(Red, s)
}
func Info(s string) {
	log(Blue, s)
}

func Success(s string) {
	log(Green, s)
}

func Warn(s string) {
	log(Yellow, s)
}

func log(color int, s string) {
	if !closePrint {
		fmt.Printf("\033[1;%v;%vm<%s> %s\033[0m\n", color, 40, time.Now().Format("2006-01-02 15:04:05"), s)
	}
	if Logger != nil {
		if color == Red {
			Logger.Error(s)
		} else if color == Yellow {
			Logger.Warn(s)
		} else {
			Logger.Info(s)
		}
	}
}
