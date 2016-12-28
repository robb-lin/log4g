package appender

import (
	. "robb-lin/log4g/logrecord"
	"fmt"
	"robb-lin/log4g/layout"
	"robb-lin/log4g/level"
)

const (
	color_red = uint8(iota + 91)
	//color_green
	//color_yellow
	color_blue = uint8(iota + 93)
	color_magenta //洋红
)

type ConsoleAppender struct {
	name	string
	layout	layout.Layout
}

func NewConsoleAppender(name string, layout layout.Layout) *ConsoleAppender {
	return &ConsoleAppender{name: name, layout: layout}
}

func (this *ConsoleAppender) GetAppenderName() string {
	if this.name == "" {
		return "ConsoleAppender"
	}
	return this.name
}

func (this *ConsoleAppender) GetLayout() layout.Layout {
	return this.layout
}


func (this *ConsoleAppender) DoAppend(record LogRecord) {
	//do something
	msg := this.layout.Log4gFormat(record)
	switch record.Level {
	case level.INFO:
		msg = fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_blue, msg)
	case level.WARN:
		msg = fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_magenta, msg)
	case level.ERROR:
		msg = fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_red, msg)
	case level.FATAL:
		msg = fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_red, msg)
	}
	fmt.Print(msg)
}