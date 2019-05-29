package logger

import (
	"github.com/kataras/golog"
)

type Logger struct {
	golog *golog.Logger
}

func NewLogger() *Logger {
	var logger Logger
	logger.golog = golog.New().SetTimeFormat("2006/01/02 15:04:05").SetLevel("debug")
	return &logger
}

func (l *Logger) Info(vl interface{}) {
	//time:=time.Now().Format("2006/01/02 15:04:05")
	l.golog.Info("[GAR] ", vl)
}