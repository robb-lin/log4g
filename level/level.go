package level

import (
	"strings"
)

/**
Logging level
 */
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	levelStringArray = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

func (l Level) String() string {
	if l < 0 || int(l) > len(levelStringArray) {
		return "UNKNOWN"
	}
	return levelStringArray[int(l)]
}

func (l Level) IsGreaterOrEqual(level Level) bool {
	if(int(l) >= int(level)) {
		return true
	}
	return false
}

/**
如果查找不到则返回DEBUG
 */
func GenLevel(levelName string) Level {
	for index, levelStr := range levelStringArray {
		if strings.ToUpper(strings.TrimSpace(levelName)) == levelStr {
			return Level(index)
		}
	}
	return DEBUG
}