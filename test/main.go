package main

import (
	"fmt"
	"robb-lin/log4g"
	"os/signal"
	"os"
	"time"
	"robb-lin/log4g/level"
	"syscall"
	"path/filepath"
	"sync"
)

func main() {
	//testAny();
	testLogger()
	//testDailyLogger()
	//testSignal()
	//os.Rename("D:/data/gologs/dailyout.test", "D:/data/gologs/dailyout.test1")
	//fmt.Println(os.Remove("D:/data/gologs/bbbb"))
}

func testDailyLogger() {
	fmt.Println(time.Now())
	logger, _ := log4g.GetLogger("test_logger")
	for i := 0; i < 600; i++ {
		logger.Info(fmt.Sprintf("%s%d", "sdfsd", i))  // 输出
		time.Sleep(time.Second)
	}

	log4g.FlushAll()
	fmt.Println("End: ", time.Now())
}

func testRollingLogger() {
	var wg sync.WaitGroup
	fmt.Println(time.Now())
	logger, _ := log4g.GetLogger("test_logger")
	for t := 0; t < 50 ; t++ {
		wg.Add(1)
		go func() {
			//defer wg.Add(-1)
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				logger.Info(fmt.Sprintf("%s%d", "sdfsd", i))  // 输出
			}
		}()
	}

	wg.Wait()

	//logger, _ := log4g.GetLogger("test_logger")
	//logger.Info(fmt.Sprintf("%s%d", "sdfsd", 100))  // 输出

	log4g.FlushAll()
	fmt.Println("End: ", time.Now())
}

func testLogger() {
	logger := log4g.GetRootLogger()
	logger.Info(fmt.Sprintf("%s%d", "abscsd", 12345))
}

func testAny() {
	bb := "bbbb"
	fmt.Println(bb)
	bb = bb + "kkk"
	fmt.Println(bb)
	fmt.Println(level.DEBUG)
	fmt.Println("'.'yyyy-MM-dd'.log'")


	fmt.Println(filepath.Base(os.Args[0]))
}

func testSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Kill, os.Interrupt, syscall.SIGALRM, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for{
			s := <-c
			fmt.Println("get signal:", s)
			os.Exit(1)
		}
	}()
	time.Sleep(time.Hour)
}