package base

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"time"
)

var Logger *slog.Logger

func InitLog() {
	r := &lumberjack.Logger{
		Filename:   Config().Application.Log.Path + "/runtime.log",
		LocalTime:  true,
		MaxSize:    Config().Application.Log.Size,
		MaxAge:     Config().Application.Log.Age,
		MaxBackups: Config().Application.Log.Backups,
		Compress:   false,
	}
	Logger = slog.New(slog.NewJSONHandler(r, &slog.HandlerOptions{
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

func Println(c1, c2 int, s string) {
	fmt.Printf("\033[1;%v;%vm%s\033[0m\n", c1, c2, s)
}
