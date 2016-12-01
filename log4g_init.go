package log4g

import (
	"strconv"
	"robb-lin/log4g/appender"
	"robb-lin/log4g/level"
	"os"
	"fmt"
	"io/ioutil"
	"runtime"
	"path/filepath"
	"encoding/xml"
	"robb-lin/log4g/layout"
	"strings"
)

func init() {
	_,logManagerFilePath,_,_ := runtime.Caller(0)
	// ===== 加载配置文件 =====
	filename := os.Getenv("LOG4G_XML")
	existXmlConf := true
	if filename == "" {
		filename = filepath.Join(filepath.Dir(logManagerFilePath), DEFAULT_XML_CONFIGURATION_FILE)
		_, err := os.Stat(filename)
		if err != nil && !os.IsExist(err){
			filename = filepath.Join(filepath.Dir(os.Args[0]), DEFAULT_XML_CONFIGURATION_FILE)
			_, err = os.Stat(filename)
			if err != nil && !os.IsExist(err) {
				existXmlConf = false
			}
		}
	}
	if !existXmlConf {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Not found %q for reading\n", DEFAULT_XML_CONFIGURATION_FILE)
		os.Exit(1)
	}
	fd, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not open %q for reading: %s\n", filename, err)
		os.Exit(1)
	}

	contents, err := ioutil.ReadAll(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not read %q: %s\n", filename, err)
		os.Exit(1)
	}

	// ===== xml配置文件解析 =====
	xc := new(xmlLoggerConfig)
	if err := xml.Unmarshal(contents, xc); err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not parse XML configuration in %q: %s\n", filename, err)
		os.Exit(1)
	}
	// ===========  解析appender  ==============
	appenders := make(map[string]appender.Appender)
	// 默认appender
	appenders["default"] = appender.NewConsoleAppender("DefaultConsoleAppender", layout.NewPatternLayoutDefault())
	// 解析配置
	for _, appenderConf := range xc.Appenders {
		layoutCls := appenderConf.Layout.Class       // appender的layout配置类
		layoutParams := appenderConf.Layout.Params   // appender的layout配置类的参数
		var appenderLayout layout.Layout
		// 构造layout
		switch layoutCls {
		case "PatternLayout":
			for _, layoutParam := range layoutParams {
				if layoutParam.Name == "ConversionPattern" {
					appenderLayout = layout.NewPatternLayoutCunstomized(layoutParam.Value)
				}
			}
		case "SimpleLayout":
			appenderLayout = layout.NewSimpleLayout()
		}
		if appenderLayout == nil {
			appenderLayout = layout.NewPatternLayoutDefault()
		}

		// 构造appender
		switch appenderConf.Class {
		case "ConsoleAppender":
			appenders[appenderConf.Name] = appender.NewConsoleAppender(appenderConf.Name, appenderLayout)
		case "RollingFileAppender":
			paramKvs := make(map[string]string)
			for _, param := range appenderConf.Params {
				if param.Name != "" && param.Value != "" {
					paramKvs[param.Name] = param.Value
				}
			}
			filename, found := paramKvs["File"]
			if !found {
				fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: RollingFileAppender the param name for 'File' Could not Empty\n")
				//os.Exit(1)
				break
			}
			fileAppend := true
			if fileAppendStr, found := paramKvs["Append"]; found {
				fileAppend, _ = strconv.ParseBool(fileAppendStr)
			}
			bufferedIO := false
			if bufferedIOStr, found := paramKvs["BufferedIO"]; found {
				bufferedIO, _ = strconv.ParseBool(bufferedIOStr)
			}
			bufferSize := appender.FILE_APPENDER_DEFAULT_BUFFER_SIZE
			if bufferSizeStr, found := paramKvs["BufferSize"]; found {
				buffSizeTmp, _ := strconv.ParseInt(bufferSizeStr, 10, 32)
				bufferSize = int(buffSizeTmp)
			}
			maxFileSize := int64(102400)  // 100KB
			if maxFileSizeStr, found := paramKvs["MaxFileSize"]; found {
				maxFileSizeStr = strings.ToUpper(maxFileSizeStr)
				if strings.HasSuffix(maxFileSizeStr, "KB") {
					filesizeNum, _ := strconv.ParseInt(maxFileSizeStr[:strings.LastIndex(maxFileSizeStr, "KB")], 10, 32)
					maxFileSize = filesizeNum * 1024
				} else if strings.HasSuffix(maxFileSizeStr, "MB") {
					filesizeNum, _ := strconv.ParseInt(maxFileSizeStr[:strings.LastIndex(maxFileSizeStr, "MB")], 10, 32)
					maxFileSize = filesizeNum * 1024 * 1024
				} else if strings.HasSuffix(maxFileSizeStr, "GB") {
					filesizeNum, _ := strconv.ParseInt(maxFileSizeStr[:strings.LastIndex(maxFileSizeStr, "GB")], 10, 32)
					maxFileSize = filesizeNum * 1024 * 1024 * 1024
				}
			}
			maxBackupIndex := 10
			if maxBackupIndexStr, found := paramKvs["MaxBackupIndex"]; found {
				maxBackupIndexTmp, _ := strconv.ParseInt(maxBackupIndexStr, 10, 32)
				maxBackupIndex = int(maxBackupIndexTmp)
			}
			//filename, fileAppend [, buffered_io, buffer_size], maxFileSize, maxBackupIndex
			appenders[appenderConf.Name] = appender.NewRollingFileAppenderWithBuffered(appenderConf.Name, appenderLayout,
				filename, fileAppend, bufferedIO, bufferSize, maxFileSize, maxBackupIndex)
		case "DailyRollingFileAppender":
			paramKvs := make(map[string]string)
			for _, param := range appenderConf.Params {
				if param.Name != "" && param.Value != "" {
					paramKvs[param.Name] = param.Value
				}
			}
			filename, found := paramKvs["File"]
			if !found {
				fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: DailyRollingFileAppender the param name for 'File' Could not Empty\n")
				//os.Exit(1)
				break
			}
			fileAppend := true
			if fileAppendStr, found := paramKvs["Append"]; found {
				fileAppend, _ = strconv.ParseBool(fileAppendStr)
			}
			bufferedIO := false
			if bufferedIOStr, found := paramKvs["BufferedIO"]; found {
				bufferedIO, _ = strconv.ParseBool(bufferedIOStr)
			}
			bufferSize := appender.FILE_APPENDER_DEFAULT_BUFFER_SIZE
			if bufferSizeStr, found := paramKvs["BufferSize"]; found {
				buffSizeTmp, _ := strconv.ParseInt(bufferSizeStr, 10, 32)
				bufferSize = int(buffSizeTmp)
			}

			datePattern, found := paramKvs["DatePattern"]
			if !found {
				datePattern = appender.DAILY_APPENDER_DEFAULT_DATE_PATTERN
			}

			//filename, fileAppend [, buffered_io, buffer_size], datePattern
			appenders[appenderConf.Name] = appender.NewDailyRollingFileAppenderWithBuffered(appenderConf.Name, appenderLayout,
				filename, fileAppend, bufferedIO, bufferSize, datePattern)
		}
	}

	// ======= 解析root  ======
	rootLoggerLevel := level.GenLevel(xc.RootLogger.Priority)
	var rootLoggerAppender appender.Appender = appenders["default"]
	if rla, found := appenders[xc.RootLogger.AppenderRef]; found {
		rootLoggerAppender = rla
	}

	rootLogger := Logger{name: "root", level: rootLoggerLevel, appender: rootLoggerAppender, additivity: false}
	ROOT_LOGGER = rootLogger

	// ======= 解析loggers  ======
	if len(xc.Loggers) > 0 {
		for _, loggerConf := range xc.Loggers {
			loggerLevel := level.GenLevel(loggerConf.Priority)
			var loggerAppender appender.Appender = appenders["default"]
			if rla, found := appenders[loggerConf.AppenderRef]; found {
				loggerAppender = rla
			}
			var additivity bool = true
			if parseRes, err := strconv.ParseBool(loggerConf.Additivity); err == nil {
				additivity = parseRes
			}

			logger := Logger{name: loggerConf.Name, level: loggerLevel, appender: loggerAppender, additivity: additivity}
			loggerManager[logger.name] = logger
		}
	}
}
