package log4g

import (
	"time"
	"runtime"
	"fmt"
	. "robb-lin/log4g/level"
	. "robb-lin/log4g/appender"
	"robb-lin/log4g/logrecord"
)

type Logger struct {
	name string
	level Level
	appender Appender
	additivity bool
}

//func NewLogger(name string, level Level, appender appender.Appender, additivity bool) *Logger {
//	return &Logger{name, level, appender, additivity}
//}

func (this Logger) getName() string {
	return this.name
}

func (this Logger) getLevel() Level {
	return this.level
}

func (this Logger) getAppender() Appender {
	return this.appender
}

func (this Logger) getAdditivity() bool {
	return this.additivity
}

func forcedLog(logger Logger, level Level, message interface{})  {
	//
	pc,file,line,_ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	record := logrecord.LogRecord{Level:level, Created:time.Now(), Caller: f.Name(), Message:fmt.Sprint(message), File:file, Line:line}
	//fmt.Printf("%+v\n", record)
	logger.appender.DoAppend(record)

	//rootLogger
	rootLogger := GetRootLogger()
	if rootLogger != logger && logger.additivity && level.IsGreaterOrEqual(rootLogger.level) {
		rootLogger.appender.DoAppend(record)
	}
}

func (this Logger) Debug(message interface{}) {
	curLevel := DEBUG
	if(curLevel.IsGreaterOrEqual(this.level)) {
		forcedLog(this, curLevel, message)
	}
}

func (this Logger) Info(message interface{}) {
	curLevel := INFO
	if(curLevel.IsGreaterOrEqual(this.level)) {
		forcedLog(this, curLevel, message)
	}
}

func (this Logger) Warn(message interface{}) {
	curLevel := WARN
	if(curLevel.IsGreaterOrEqual(this.level)) {
		forcedLog(this, curLevel, message)
	}
}

func (this Logger) Error(message interface{}) {
	curLevel := ERROR
	if(curLevel.IsGreaterOrEqual(this.level)) {
		forcedLog(this, curLevel, message)
	}
}

func (this Logger) Fatal(message interface{}) {
	curLevel := FATAL
	if(curLevel.IsGreaterOrEqual(this.level)) {
		forcedLog(this, curLevel, message)
	}
}