package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/wudidayan/go-blog/pkg/file"
	"github.com/wudidayan/go-blog/pkg/setting"
)

var (
	fd         *os.File
	logger     *log.Logger
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG int = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func Setup() {
	var err error
	filePath := getLogFilePath()
	fileName := getLogFileName()
	fd, err = openLogFile(fileName, filePath)
	if err != nil {
		log.Fatalln(err)
	}

	logger = log.New(fd, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	Info("logging Setup() ok!")
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v)
}

func setPrefix(level int) {
	var logPrefix string
	_, file, line, ok := runtime.Caller(2)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], file, line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}

func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))
}

func openLogFile(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := file.CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = file.IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := file.Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}
