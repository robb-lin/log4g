package appender

import (
	"robb-lin/log4g/layout"
	. "robb-lin/log4g/logrecord"
	"path/filepath"
	"os"
	"fmt"
	"sync"
	"bytes"
)

type RollingFileAppender struct {
	fileAppender
	maxFileSize int64
	maxBackupIndex int
}

var rfaLock sync.Mutex

func NewRollingFileAppenderWithBuffered(name string, layout layout.Layout, filename string, fileAppend bool, bufferedIo bool, bufferSize int, maxFileSize int64, maxBackupIndex int) *RollingFileAppender {
	entity := &RollingFileAppender{fileAppender{name: name, layout:layout, filename:filename,
						fileAppend: fileAppend, bufferedIo:bufferedIo, bufferSize:bufferSize}, maxFileSize, maxBackupIndex}

	rfaLock.Lock()
	entity.createWriter()
	if entity.writeCount > entity.maxFileSize {
		entity.rollOver()
	}
	rfaLock.Unlock()
	if bufferedIo && bufferSize > 0 {
		entity.writeBuf = bytes.NewBuffer(make([]byte, 0, bufferSize))
	}
	return entity
}

func NewRollingFileAppender(name string, layout layout.Layout, filename string, fileAppend bool, maxFileSize int64, maxBackupIndex int) *RollingFileAppender {
	return NewRollingFileAppenderWithBuffered(name, layout, filename, fileAppend, false, FILE_APPENDER_DEFAULT_BUFFER_SIZE, maxFileSize, maxBackupIndex)
}

/**
    初始化写文件
 */
func (this *RollingFileAppender) createWriter() {
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

/**
    文件滚动处理
 */
func (this *RollingFileAppender) rollOver() {
	if this.writeCount < this.maxFileSize {
		return
	}

	renameSucceeded := true

	backupFileName := this.filename + "." + fmt.Sprintf("%d", this.maxBackupIndex)
	if this.maxBackupIndex > 0 {
		_, err := os.Stat(backupFileName)
		if err == nil || os.IsExist(err) {
			renameSucceeded = (os.Remove(backupFileName) == nil)
		}

		for i := (this.maxBackupIndex - 1); (i >= 1) && renameSucceeded; i-- {
			backupFileName = this.filename + "." + fmt.Sprintf("%d", i)
			_, err := os.Stat(backupFileName)
			if err == nil || os.IsExist(err) {
				targetFileName := (this.filename + "." + fmt.Sprintf("%d", (i + 1)))
				renameSucceeded = (os.Rename(backupFileName, targetFileName) == nil)
			}
		}

		if renameSucceeded {
			targetFileName := this.filename + ".1"
			this.file.Close()
			renameSucceeded = (os.Rename(this.filename, targetFileName) == nil)
		}
	}

	if renameSucceeded {
		this.createWriter()
	}
}

func (this *RollingFileAppender) GetMaxBackupIndex() int {
	return this.maxBackupIndex
}

func (this *RollingFileAppender) GetMaxFileSize() int64 {
	return this.maxFileSize
}

func (this *RollingFileAppender) GetAppenderName() string {
	if this.name == "" {
		return "RollingFileAppender"
	}
	return this.name
}

func (this *RollingFileAppender) DoAppend(record LogRecord) {
	targetMsg := this.layout.Log4gFormat(record)
	rfaLock.Lock()
	if this.file != nil {
		if this.bufferedIo {
			if (this.bufferSize -  this.writeBuf.Len()) >= len(targetMsg) {
				this.writeBuf.WriteString(targetMsg)
				targetMsg = ""
			} else {
				tmp := targetMsg
				targetMsg = this.writeBuf.String()
				this.writeBuf.Reset();
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

	if this.filename != "" && this.file != nil && this.writeCount > this.maxFileSize {
		this.rollOver()
	}
	rfaLock.Unlock()
}

/**
  刷新缓存
 */
func (this *RollingFileAppender) Flush() {
	if this.bufferedIo && this.writeBuf.Len() > 0 {
		rfaLock.Lock()
		targetMsg := this.writeBuf.String()
		if this.filename != "" && this.file != nil {
			ret, err := this.file.WriteString(targetMsg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "flush %s write log error: %v\n", this.filename, err)
			}
			this.writeCount += int64(ret)
			if this.writeCount > this.maxFileSize {
				this.rollOver()
			}
		}
		rfaLock.Unlock()
	}
}