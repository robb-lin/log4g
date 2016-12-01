package appender

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"robb-lin/log4g/layout"
	. "robb-lin/log4g/logrecord"
	"strings"
	"sync"
	"time"
)

type DailyRollingFileAppender struct {
	fileAppender
	datePattern      string
	patternHandler   dailyPatternHandler
	nextCreateFileAt time.Time
}

var drfaLock sync.Mutex

func NewDailyRollingFileAppenderWithBuffered(name string, layout layout.Layout, filename string, fileAppend bool, bufferedIo bool, bufferSize int, datePattern string) *DailyRollingFileAppender {
	if datePattern == "" {
		datePattern = DAILY_APPENDER_DEFAULT_DATE_PATTERN
	}
	patternHandler := genPatternHandler(datePattern)
	entity := &DailyRollingFileAppender{fileAppender{name: name, layout: layout, filename: filename,
		fileAppend: fileAppend, bufferedIo: bufferedIo, bufferSize: bufferSize}, datePattern, patternHandler, time.Now().AddDate(0, 0, 1)}

	drfaLock.Lock()
	entity.createWriter()
	entity.rollOver()
	drfaLock.Unlock()
	if bufferedIo && bufferSize > 0 {
		entity.writeBuf = bytes.NewBuffer(make([]byte, 0, bufferSize))
	}
	return entity
}

func NewDailyRollingFileAppender(name string, layout layout.Layout, filename string, fileAppend bool, datePattern string) *DailyRollingFileAppender {
	return NewDailyRollingFileAppenderWithBuffered(name, layout, filename, fileAppend, false, FILE_APPENDER_DEFAULT_BUFFER_SIZE, datePattern)
}

/**
  初始化写文件
*/
func (this *DailyRollingFileAppender) createWriter() {
	//目录判断, 不存在则创建
	fdDir := filepath.Dir(this.filename)
	var _, err = os.Stat(fdDir)
	if os.IsNotExist(err) {
		os.MkdirAll(fdDir, 0777)
	}
	//判断是否文件追加
	if !this.fileAppend {
		os.Remove(this.filename)
	}
	//打开文件，若不存在自动创建
	fd, err := os.OpenFile(this.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s open error\n", this.filename)
	} else {
		this.file = fd
	}
	//初始化
	fileinfo, err := os.Stat(this.filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s get fileinfo error\n", this.filename)
		this.writeCount = 0
	} else {
		this.writeCount = fileinfo.Size()
	}

	if this.bufferedIo && this.writeBuf == nil {
		bs := make([]byte, 0, this.bufferSize)
		this.writeBuf = bytes.NewBuffer(bs)
	}
}

func (this *DailyRollingFileAppender) rollOver() {
	needCreateNewFile, modTime := this.checkNeedCreateNewFile()
	if needCreateNewFile {
		targetFileName := this.filename + this.patternHandler.PatternFormateFunc(modTime)
		_, err := os.Stat(targetFileName)
		if err == nil || os.IsExist(err) {
			os.Remove(targetFileName)
		}
		if this.bufferedIo && this.writeBuf.Len() > 0 {
			this.file.WriteString(this.writeBuf.String())
		}

		this.file.Close()
		err = os.Rename(this.filename, targetFileName)
		this.createWriter()
	}
}

func (this *DailyRollingFileAppender) checkNeedCreateNewFile() (bool, time.Time) {
	needCreateNewFile := false
	fileinfo, err := os.Stat(this.filename)
	var modTime time.Time
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s get fileinfo error\n", this.filename)
	} else {
		modTime = fileinfo.ModTime()
		mody, modm, modd := modTime.Date()
		modh := modTime.Hour()
		curTime := time.Now()
		cury, curm, curd := curTime.Date()
		curh := curTime.Hour()
		rollingType := this.patternHandler.RollingType
		switch rollingType {
		case 1: //year
			if mody != cury {
				needCreateNewFile = true
			}
			this.nextCreateFileAt = time.Date(cury+1, 1, 1, 0, 0, 0, 0, time.Local)
		case 2:
			if mody != cury || modm != curm {
				needCreateNewFile = true
			}
			if curm == 12 {
				this.nextCreateFileAt = time.Date(cury+1, 1, 1, 0, 0, 0, 0, time.Local)
			} else {
				this.nextCreateFileAt = time.Date(cury, curm+1, 1, 0, 0, 0, 0, time.Local)
			}
		case 3:
			if mody != cury || modm != curm || modd != curd {
				needCreateNewFile = true
			}
			this.nextCreateFileAt = time.Date(cury, curm, curd, 0, 0, 0, 0, time.Local).Add(24 * time.Hour)
		case 4:
			if mody != cury || modm != curm || modd != curd || modh != curh {
				needCreateNewFile = true
			}
			this.nextCreateFileAt = time.Date(cury, curm, curd, curh, 0, 0, 0, time.Local).Add(time.Hour)
		}
	}

	return needCreateNewFile, modTime
}

func (this *DailyRollingFileAppender) GetDatePattern() string {
	return this.datePattern
}

func (this *DailyRollingFileAppender) GetAppenderName() string {
	if this.name == "" {
		return "DailyRollingFileAppender"
	}
	return this.name
}

func (this *DailyRollingFileAppender) DoAppend(record LogRecord) {
	targetMsg := this.layout.Log4gFormat(record)
	rfaLock.Lock()
	if this.nextCreateFileAt.Before(time.Now()) {
		this.rollOver()
	}

	if this.file != nil {
		if this.bufferedIo {
			if (this.bufferSize - this.writeBuf.Len()) >= len(targetMsg) {
				this.writeBuf.WriteString(targetMsg)
				targetMsg = ""
			} else {
				tmp := targetMsg
				targetMsg = this.writeBuf.String()
				this.writeBuf.Reset()
				this.writeBuf.WriteString(tmp)
			}
		}
		if targetMsg != "" {
			ret, err := this.file.WriteString(targetMsg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s write log error: %v\n", this.filename, err)
			}
			this.writeCount += int64(ret)
		}
	} else {
		fmt.Fprintf(os.Stderr, "%s write log error, this.file is nil\n", this.filename)
	}
	rfaLock.Unlock()
}

/**
  刷新缓存
*/
func (this *DailyRollingFileAppender) Flush() {
	if this.bufferedIo && this.writeBuf.Len() > 0 {
		rfaLock.Lock()
		targetMsg := this.writeBuf.String()
		if this.filename != "" && this.file != nil {
			ret, err := this.file.WriteString(targetMsg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "flush %s write log error\n", this.filename)
			}
			this.writeCount += int64(ret)
			if this.nextCreateFileAt.Before(time.Now()) {
				this.rollOver()
			}
		}
		rfaLock.Unlock()
	}
}

type dailyPatternHandler struct {
	DatePattern        string
	RollingType        int // 1: 按年  2：按月 3：按天 4：按时段
	PatternFormateFunc func(time time.Time) string
}

func genPatternHandler(datePattern string) dailyPatternHandler {
	//"yyyy/MM/dd HH:mm:ss,SSS"
	//"2006/01/02 15:04:05"
	storePattern := datePattern
	handler := dailyPatternHandler{storePattern, 1, nil}
	yCount := strings.Count(datePattern, "y")
	if yCount > 0 {
		handler.RollingType = 1
		if yCount > 2 {
			datePattern = strings.Replace(datePattern, strings.Repeat("y", yCount), "2006", 1)
		} else {
			datePattern = strings.Replace(datePattern, strings.Repeat("y", yCount), "06", 1)
		}
	}
	MCount := strings.Count(datePattern, "M")
	if MCount > 0 {
		handler.RollingType = 2
		datePattern = strings.Replace(datePattern, strings.Repeat("M", MCount), "01", 1)
	}
	dCount := strings.Count(datePattern, "d")
	if dCount > 0 {
		handler.RollingType = 3
		datePattern = strings.Replace(datePattern, strings.Repeat("d", dCount), "02", 1)
	}
	HCount := strings.Count(datePattern, "H")
	if HCount > 0 {
		handler.RollingType = 4
		datePattern = strings.Replace(datePattern, strings.Repeat("H", HCount), "15", 1)
	}

	handler.PatternFormateFunc = func(datetime time.Time) string {
		res := datetime.Format(datePattern)
		return res
	}
	return handler
}
