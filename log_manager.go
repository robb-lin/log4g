package log4g

import (
	"robb-lin/log4g/appender"
	"reflect"
)

var loggerManager map[string]Logger = make(map[string]Logger)
var ROOT_LOGGER Logger

func GetLoggerManager() map[string]Logger {
	return loggerManager
}

func GetRootLogger() Logger {
	return ROOT_LOGGER
}

/**
  根据日志名称获取日志对象，如果找不到，则返回rootLogger
 */
func GetLogger(name string) (Logger, bool) {
	if log, found := loggerManager[name]; found {
		return log, found
	} else {
		return ROOT_LOGGER, found
	}
}

/**
   针对开启buff功能的logger全量刷新有缓存的logger
 */
func FlushAll() {
	//rootLogger
	loggers := []Logger{ROOT_LOGGER}
	//其他Logger
	for _, log := range loggerManager {
		loggers = append(loggers, log)
	}
	for _, log := range loggers {
		flushHandler(&log)
	}
}

/**
  针对开启buff功能的logger刷新数据
 */
func Flush(logname string) {
	if log, found := loggerManager[logname]; found {
		flushHandler(&log)
	}
}

func FlushRootLogger() {
	flushHandler(&ROOT_LOGGER)
}

func flushHandler(log *Logger) {
	appenderValue := reflect.ValueOf(log.getAppender()).Interface()
	switch appenderValue.(type) {
	case *appender.RollingFileAppender :
		rfa := appenderValue.(*appender.RollingFileAppender)
		rfa.Flush()
	case *appender.DailyRollingFileAppender:
		drfa := appenderValue.(*appender.DailyRollingFileAppender)
		drfa.Flush()
	}
}