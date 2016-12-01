package layout

import (
	. "robb-lin/log4g/logrecord"
	"fmt"
)

type SimpleLayout struct {

}

func NewSimpleLayout() *SimpleLayout {
	return &SimpleLayout{}
}

func (this SimpleLayout) Log4gFormat(record LogRecord) string {
	return fmt.Sprintf("[%s] %s %s %s:%d - %s \n" ,
			record.Created.Format("2006/01/02 15:04:05"),
			record.Level.String(),
			record.File,
			record.Caller,
			record.Line,
			record.Message)
}