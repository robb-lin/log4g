package appender

import (
	. "robb-lin/log4g/logrecord"
	"fmt"
	"robb-lin/log4g/layout"
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
	fmt.Print(this.layout.Log4gFormat(record))
}