package log4g

import (
	"testing"
	"fmt"
	"os"
	"bytes"
	"os/signal"
	"time"
)

func TestAny(t *testing.T) {
	//t.Log([]interface{}{1, 2, 3}[:0])
	//t.Log(int((0 + 1) / 2))
	//layout.PrintLayout()
	//fmt.Println()
	//fmt.Println(DEBUG.String())
	//var cappdender appender.Appender = appender.ConsoleAppender{}
	//fmt.Println(cappdender.GetAppenderName())

	//pc,file,line,ok := runtime.Caller(2)
	//fmt.Println(pc)
	//fmt.Println(file)
	//fmt.Println(line)
	//fmt.Println(ok)
	//f := runtime.FuncForPC(pc)
	//fmt.Println(f.Name())
	//fmt.Println(GetLogger("bb"))

	filename := "D:/mydocs/t1231232.txt"
	//_, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "%s open error", filename)
	//}
	os.Remove(filename)
	fileinfo, err := os.Stat(filename)
	fmt.Printf("%+v, %v\n", fileinfo, err)
}

func cleanup() {
	fmt.Println("cleanup")
}

func TestLogger(t *testing.T) {
	/*
	name string
	level Level
	appender appender.Appender
	additivity bool
	*/
	//logger := Logger{name: "BB", level: level.INFO, appender: appender.NewConsoleAppender("BB", layout.NewPatternLayoutDefault()), additivity:false}
	logger, _ := GetLogger("test_logger")
	//fmt.Println(logger.additivity)
	//fmt.Printf("%+v", logger.appender)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	go func(){
		// Block until a signal is received.
		s := <-c
		fmt.Println("Got signal:", s)
	}()

	logger.Debug(12345)  // 不输出
	for i := 1; i < 10000; i++ {
		logger.Info(fmt.Sprintf("%s%d", "sdfsd", i))  // 输出
	}

	//logger.Info("sdfsdf22")  // 输出
	//
	//fmt.Println(level.GenLevel("iNFo"))
	//GetRootLogger().Info("bababa")
	//name string, layout layout.Layout, filename string, file_append bool, immediate_flush bool, max_file_size int64, max_backup_index int
	//fmt.Printf("%s - %+v", appenderConf.Name, appenders[appenderConf.Name]);
}

func TestBytesBuf(t *testing.T) {
	bs := make([]byte, 0, 50)
	byteBuf := bytes.NewBuffer(bs)
	for i := 1; i < 100; i++ {
		byteBuf.WriteString("波波")
		//for _, b := range []byte("波波") {
		//	bs = append(bs, b)
		//}
	}
	byteBuf.Reset()
	fmt.Println(byteBuf.Len())

	for i := 1; i < 100; i++ {
		byteBuf.WriteString("波波")
		//for _, b := range []byte("波波") {
		//	bs = append(bs, b)
		//}
	}

	fmt.Println(byteBuf.Len())
}


func TestFileInfoOp(t *testing.T) {
	tz, _ := time.LoadLocation("Asia/Shanghai")
	fmt.Println(tz.String());
	fmt.Println(time.Local);
	filename := "C:/Users/linrb/Desktop/commit/datasync-device-1.0.jar"
	fileinfo, _ := os.Stat(filename)
	fmt.Printf("%+v\n\n", fileinfo)
	fmt.Printf("%+v\n\n", fileinfo.ModTime().Format("2006-01-02 15:04:05"))
	fmt.Printf("%+v\n\n", fileinfo.ModTime().Truncate(24 * time.Hour))
	//datetime.UnixNano() / int64(time.Millisecond) % 1000
	fmt.Printf("%+v\n\n", time.Now().Format("2006-01-02 15:04:05"))

	fmt.Printf("%+v\n\n", time.Now().Truncate(24 * time.Hour).Unix())

	fmt.Printf("%+v\n\n", time.Now().Unix())

	y, m, d := time.Now().Date()
	fmt.Println(y, m, d)
	//h := time.Now().Hour()
	hTime := time.Date(y + 1, 1, 1, 0, 0, 0, 0, time.Local)
	fmt.Println(time.Now().Before(hTime))
	fmt.Printf("%+v\n\n", hTime)

}
