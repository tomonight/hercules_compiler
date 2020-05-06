package log

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

//定义日志等级
const (
	PanicLevel uint32 = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

//定义日志最大行数
const (
	LogMaxLine = 15000
)

var (
	gStdLog      *logrus.Logger
	gFileLog     *logrus.Logger
	gLogInit     bool
	gPrint       bool
	gLogFileName string
	startDate    string
	curLogLine   int
	curIndex     int
)

func init() {
	gStdLog = nil
	gFileLog = nil
	gLogInit = false
	gPrint = false
	startDate = time.Now().Format("20060102")
	curLogLine = 0
	curIndex = 0
}

func getLogFileWriter(fileName string) (io.Writer, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

//@desc get log file line count
func getLogFileLineCount(logFile string) int {
	file, err := os.Open(logFile)
	if err != nil {
		return 0
	}
	defer file.Close()
	lineCount := 0
	buf := bufio.NewReader(file)
	for {
		_, _, err := buf.ReadLine()
		if err != nil || io.EOF == err {
			break
		}
		lineCount++
	}
	return lineCount
}

//@desc 获取当前日志文件路径
func getLocalLogFilePath(logFile string) string {
	var lastLogFile string
	logFileList := []string{}

	//log foramt logfile-20180522-1.log logfile-20180522-2.log
	dir := path.Dir(logFile)
	ext := path.Ext(logFile)
	basePath := strings.Split(path.Base(logFile), ext)[0]
	dateStr := time.Now().Format("20060102")

	logFilePre := basePath + "_" + dateStr + "_"
	regfileName := logFilePre + "*" + ext

	filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			ok, err := filepath.Match(regfileName, info.Name())
			if ok {
				if !info.IsDir() {
					logFileList = append(logFileList, info.Name())
				}
				return nil
			}
			return err
		})
	index := 1
	if len(logFileList) > 0 {
		for _, tmpFile := range logFileList {
			tmpStr := strings.Replace(tmpFile, logFilePre, "", -1)
			tmpStr = strings.Replace(tmpStr, ext, "", -1)
			curIndex, err := strconv.Atoi(tmpStr)
			if err == nil {
				if curIndex > index {
					index = curIndex
				}
			}
		}
	}
	lastLogFile = path.Join(dir, logFilePre+strconv.Itoa(index)+ext)

	logFileLineCount := getLogFileLineCount(lastLogFile)
	if logFileLineCount >= LogMaxLine {
		index++
		lastLogFile = path.Join(dir, logFilePre+strconv.Itoa(index)+ext)
		curLogLine = 0
		curIndex = index
	} else {
		curLogLine = logFileLineCount
		curIndex = index
	}

	return lastLogFile
}

//InitLog 初始化日志
func InitLog(logFile string, print bool, level uint32) error {
	if !gLogInit {
		if level < uint32(logrus.PanicLevel) || level > uint32(logrus.DebugLevel) {
			return errors.New("unrecognized log level")
		}
		if print {
			gStdLog = logrus.New()
			gStdLog.Formatter = new(logrus.TextFormatter)
			gStdLog.Level = logrus.Level(level)
			gStdLog.Out = os.Stdout
			gPrint = true
		}

		lastPath := getLocalLogFilePath(logFile)

		gLogFileName = logFile

		file, err := getLogFileWriter(lastPath)
		if err != nil {
			return err
		}

		gFileLog = logrus.New()
		gFileLog.Formatter = new(logrus.TextFormatter)
		gFileLog.Level = logrus.Level(level)
		gFileLog.Out = file
		gLogInit = true
	}
	return nil
}

//ReLoad 重新初始化日志
func ReLoad(logFile string, print bool, level uint32) error {
	gLogInit = false
	gPrint = false
	return InitLog(logFile, print, level)
}

func doSplitLog() {
	if curLogLine < LogMaxLine {
		return
	}
	curIndex++
	//log foramt logfile-20180522-1.log logfile-20180522-2.log
	dir := path.Dir(gLogFileName)
	ext := path.Ext(gLogFileName)
	basePath := strings.Split(path.Base(gLogFileName), ext)[0]
	dateStr := time.Now().Format("20060102")
	logFilePre := basePath + "_" + dateStr + "_"
	lastPath := path.Join(dir, logFilePre+strconv.Itoa(curIndex)+ext)
	file, err := getLogFileWriter(lastPath)
	if err != nil {
		return
	}
	gFileLog.Out = file
	curLogLine = 0
}

//CallerPos  skip： 0是调用函数调文件位置
func CallerPos(skip int) string {
	skip++ // 过滤调自己调层次
	file, line := FileNameAndLineNum(skip)
	return fmt.Sprintf("%s:%d", file, line)
}

//FileNameAndLineNum skip： 0是调用函数调文件位置
func FileNameAndLineNum(skip int) (string, int) {
	skip++ // 过滤调自己调层次
	_, file, line, _ := runtime.Caller(skip)
	if file != "" {
		file = path.Base(file)
	}

	return file, line
}

//Debug 调试日志
func Debug(format string, args ...interface{}) {
	if gLogInit {
		doSplitLog()
		curLogLine++
		fileAndLine := CallerPos(1)
		gFileLog.Debugf(fmt.Sprintf("%s %s", fileAndLine, format), args...)
		if gPrint {
			gStdLog.Debugf(fmt.Sprintf("%s %s", fileAndLine, format), args...)
		}
	}
}

//Error 错误日志
func Error(format string, args ...interface{}) {
	if gLogInit {
		doSplitLog()
		curLogLine++

		fileAndLine := CallerPos(1)
		gFileLog.Errorf(fmt.Sprintf("%s %s", fileAndLine, format), args...)
		if gPrint {
			gStdLog.Errorf(fmt.Sprintf("%s %s", fileAndLine, format), args...)
		}
	}
}

//Info 输出信息
func Info(format string, args ...interface{}) {
	if gLogInit {
		doSplitLog()
		curLogLine++
		fileAndLine := CallerPos(1)
		gFileLog.Infof(fmt.Sprintf("%s %s", fileAndLine, format), args...)
		if gPrint {
			gStdLog.Infof(fmt.Sprintf("%s %s", fileAndLine, format), args...)
		}
	}
}

//Warn 警告日志
func Warn(format string, args ...interface{}) {
	if gLogInit {
		doSplitLog()
		curLogLine++
		fileAndLine := CallerPos(1)
		gFileLog.Warnf(fmt.Sprintf("%s %s", fileAndLine, format), args...)
		if gPrint {
			gStdLog.Warnf(fmt.Sprintf("%s %s", fileAndLine, format), args...)
		}
	}
}

//Fatal  ƒ毁灭性的日志
func Fatal(format string, args ...interface{}) {
	if gLogInit {
		doSplitLog()
		curLogLine++
		fileAndLine := CallerPos(1)
		gFileLog.Fatalf(fmt.Sprintf("%s %s", fileAndLine, format), args...)
		if gPrint {
			gStdLog.Fatalf(fmt.Sprintf("%s %s", fileAndLine, format), args...)
		}
	}
}

//Panic 恐慌日志
func Panic(format string, args ...interface{}) {
	if gLogInit {
		doSplitLog()
		curLogLine++
		fileAndLine := CallerPos(1)
		gFileLog.Panicf(fmt.Sprintf("%s %s", fileAndLine, format), args...)
		if gPrint {
			gStdLog.Panicf(fmt.Sprintf("%s %s", fileAndLine, format), args...)
		}
	}
}

func init() {
	//commented by xiong at 2019-06-07
	//日志初始化代码应该由业务系统main中使用，或者初始化为向屏幕打印日志
	//InitLog("deploy_service.log", true, DebugLevel)
}
