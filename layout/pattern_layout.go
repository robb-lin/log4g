package layout

import (
	. "robb-lin/log4g/logrecord"
	"time"
	"strings"
	"fmt"
)

const (
	DEFAULT_CONVERSION_PATTERN = "%d{dd HH:mm:ss,SSS} %p %c:%l - %m%n"
)

type PatternLayout struct {
	pattern string
	convertFunc func(record LogRecord) string
}

func NewPatternLayoutDefault() *PatternLayout {
	return NewPatternLayoutCunstomized(DEFAULT_CONVERSION_PATTERN)
}

func NewPatternLayoutCunstomized(pattern string) *PatternLayout {
	return &PatternLayout{pattern:pattern, convertFunc: createConvertFunc(pattern)}
}

func (this *PatternLayout) getConversionPattern() string {
	return this.pattern
}

/**
       a、%c 显示logger调用者
       b、%d 显示日志记录时间 %d{yyyy/MM/dd HH:mm:ss,SSS}	2005/10/12 22:23:30,117
       c、%f 显示调用logger的文件
       d、%l 显示调用logger的代码行
       e、%m 显示输出信息
       f、%n 换行符
       g、%p 日志级别
       i、%%  显示一个百分号
 */
func (this PatternLayout) Log4gFormat(record LogRecord) string {
	return this.convertFunc(record)
}


func createConvertFunc(curPattern string) func(record LogRecord) string {
	//[%d{dd HH:mm:ss,SSS} %p] %c:%l - %m%n
	var targetList = []string{}
	var dT string
	var dateFormatFunc func(time time.Time) string
	if strings.Count(curPattern, "%d") > 0 {
		dindex := strings.Index(curPattern, "%d")
		dT = curPattern[dindex:]
		dT = dT[0:(strings.Index(dT, "}") + 1)]
		dateFormatFunc = loggerDateFormat(dT[3:strings.Index(dT, "}")])
		targetList = append(targetList, "%d")
	}
	if strings.Contains(curPattern, "%c") {
		targetList = append(targetList, "%c")
	}
	if strings.Contains(curPattern, "%f") {
		targetList = append(targetList, "%f")
	}
	if strings.Contains(curPattern, "%l") {
		targetList = append(targetList, "%l")
	}
	if strings.Contains(curPattern, "%m") {
		targetList = append(targetList, "%m")
	}
	if strings.Contains(curPattern, "%n") {
		targetList = append(targetList, "%n")
	}
	if strings.Contains(curPattern, "%p") {
		targetList = append(targetList, "%p")
	}
	if strings.Contains(curPattern, "%%") {
		targetList = append(targetList, "%%")
	}

	convertFunc := func(record LogRecord) string {
		res := curPattern + ""
		for _, t := range targetList {
			switch t {
			case "%d":
				res = strings.Replace(res, dT, dateFormatFunc(record.Created), 1)
			case "%c":
				res = strings.Replace(res, t, record.Caller, 1)
			case "%f":
				res = strings.Replace(res, t, record.File, 1)
			case "%l":
				res = strings.Replace(res, t, fmt.Sprintf("%d", record.Line) , 1)
			case "%m":
				res = strings.Replace(res, t, record.Message, 1)
			case "%n":
				res = strings.Replace(res, t, "\n", 1)
			case "%p":
				res = strings.Replace(res, t, record.Level.String(), 1)
			case "%%":
				res = strings.Replace(res, t, "%", 1)
			}
		}
		return res
	}

	return  convertFunc
}

func loggerDateFormat(format string) func(time time.Time) string {
	//"yyyy/MM/dd HH:mm:ss,SSS"
	//"2006/01/02 15:04:05"
	yCount := strings.Count(format, "y")
	if yCount > 2 {
		format = strings.Replace(format, strings.Repeat("y", yCount), "2006", 1)
	} else if yCount <=2 && yCount > 0 {
		format = strings.Replace(format, strings.Repeat("y", yCount), "06", 1)
	}
	MCount := strings.Count(format, "M")
	if MCount > 0 {
		format = strings.Replace(format, strings.Repeat("M", MCount), "01", 1)
	}
	dCount := strings.Count(format, "d")
	if dCount > 0 {
		format = strings.Replace(format, strings.Repeat("d", dCount), "02", 1)
	}
	HCount := strings.Count(format, "H")
	if HCount > 0 {
		format = strings.Replace(format, strings.Repeat("H", HCount), "15", 1)
	}
	mCount := strings.Count(format, "m")
	if mCount > 0 {
		format = strings.Replace(format, strings.Repeat("m", mCount), "04", 1)
	}
	sCount := strings.Count(format, "s")
	if sCount > 0 {
		format = strings.Replace(format, strings.Repeat("s", sCount), "05", 1)
	}
	SCount := strings.Count(format, "S")
	if SCount > 0 {
		format = strings.Replace(format, strings.Repeat("S", SCount), "==MILLISECOND_TARGET==", 1)
	}

	return func(datetime time.Time) string {
		res := datetime.Format(format)
		datetimeMilliSecond := datetime.UnixNano() / int64(time.Millisecond) % 1000
		return strings.Replace(res, "==MILLISECOND_TARGET==", fmt.Sprintf("%d", datetimeMilliSecond), 1)
	}
}
