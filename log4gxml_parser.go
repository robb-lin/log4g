package log4g

type xmlLoggerConfig struct {
	Appenders  	[]xmlAppender 	`xml:"appender"`  // 阿！！！ xml与tag之间一定不能有空格
	Loggers 	[]xmlLogger	`xml:"logger"`
	RootLogger      xmlRootLogger   `xml:"root"`
}

type xmlAppender struct {
	Name		string		`xml:"name,attr"`
	Class		string		`xml:"class,attr"`
	Layout		xmlLayout    	`xml:"layout"`
	Params 		[]xmlParam	`xml:"param"`
}

type xmlLogger struct {
	Name     	string		`xml:"name,attr"`
	Additivity    	string		`xml:"additivity,attr"`
	Priority	string		`xml:"priority"`
	AppenderRef	string		`xml:"appender-ref"`
}

type xmlRootLogger struct {
	Priority	string		`xml:"priority"`
	AppenderRef	string		`xml:"appender-ref"`
}

type xmlParam struct {
	Name		string		`xml:"name,attr"`
	Value		string		`xml:"value,attr"`
}

type xmlLayout struct {
	Class    	string		`xml:"class,attr"`
	Params 		[]xmlParam	`xml:"param"`
}