package layout

import (
	. "robb-lin/log4g/logrecord"
)

type Layout interface {
	Log4gFormat(record LogRecord) string
}
