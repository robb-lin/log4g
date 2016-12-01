package appender

import (
	"robb-lin/log4g/layout"
	"os"
	"bytes"
)

const (
	FILE_APPENDER_DEFAULT_BUFFER_SIZE = 8192
	DAILY_APPENDER_DEFAULT_DATE_PATTERN = ".yyyy-MM-dd"
)

type fileAppender struct {
	name		string
	layout		layout.Layout
	filename	string
	fileAppend	bool
	bufferedIo	bool
	bufferSize	int	// defautl 8192
	file		*os.File
	writeCount	int64
	writeBuf	*bytes.Buffer
}

func (this fileAppender) GetBufferSize() int {
	return this.bufferSize
}

func (this fileAppender) GetBufferedIo() bool {
	return this.bufferedIo
}

func (this fileAppender) GetFileAppend() bool {
	return this.fileAppend
}

func (this fileAppender) GetFileName() string {
	return this.filename
}

func (this fileAppender) GetLayout() layout.Layout {
	return this.layout
}