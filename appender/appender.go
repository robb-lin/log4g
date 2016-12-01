package appender

import (
	. "robb-lin/log4g/logrecord"
)

type Appender interface {
	GetAppenderName() string
	DoAppend(record LogRecord)
}
