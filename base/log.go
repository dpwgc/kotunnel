package base

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"time"
)

var Logger *slog.Logger

func InitLog(opt LogOptions) {
	r := &lumberjack.Logger{
		Filename:   opt.Path + "/runtime.log",
		LocalTime:  true,
		MaxSize:    opt.Size,
		MaxAge:     opt.Age,
		MaxBackups: opt.Backups,
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
	Red   = 31
	Blue  = 36
	Green = 32
)

func Tips(color int, s string, seconds ...int) {
	fmt.Printf("\033[1;%v;%vm<%s> %s\033[0m\n", color, 40, time.Now().Format("2006-01-02 15:04:05"), s)
	if color == Red {
		Logger.Error(s)
	} else {
		Logger.Info(s)
	}
	if len(seconds) > 0 && seconds[0] > 0 {
		time.Sleep(time.Duration(seconds[0]) * time.Second)
	}
}
