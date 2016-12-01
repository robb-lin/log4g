package logrecord


import (
	"time"
	. "robb-lin/log4g/level"
)

type LogRecord struct {
	Level Level   // The log Level
	Created time.Time  // The time at which the log message was created (nanoseconds)
	Caller string     // The message caller
	Message string   // The log message
	File string   // 输出日志的文件
	Line int  // 输出日志的行号
}