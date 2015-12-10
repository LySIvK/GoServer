package loger

import (
	"fmt"
	"gopath/github.com/issue9/term/colors"
	"log"
	"os"
	"sync"
	"time"
)

const (
	LogDebug = 1 ///调试
	LogInfo  = 2 ///信息
	LogWarn  = 3 ///警告
	LogError = 4 ///错误
	LogFatal = 5 ///致命错误
)

var loger *log.Logger   //! 日志组件
var std *log.Logger     //! 终端日志组件
var LogMinLevel int     //! 最小日志等级
var terminalOutput bool //! 是否同时终端输出
var currentDay int      //! 当前天
var logPath string      //! 当前日志文件路径
var fileName string     //! 日志文件名字

var outputLock sync.RWMutex

func Output(logType int, format string, v ...interface{}) { ///信息
	if logType < LogMinLevel {
		return ///过滤此等级
	}
	outputLock.Lock()
	//colorDic := map[string]string{"[Debug]": CLR_G, "[Info]": CLR_B, "[Warn]": CLR_Y, "[Error]": CLR_R, "[Fatal]": CLR_P}
	prefixStr := "unknow"
	c := colors.New(colors.Stdout, colors.White, colors.Black)
	c.SetColor(colors.White, colors.Black)
	switch logType { //! 根据消息类型自由搭配颜色
	case LogDebug:
		prefixStr = "[Debug]"
		c.SetColor(colors.Green, colors.Black)
		c = colors.New(colors.Stdout, colors.Green, colors.Black)
	case LogInfo:
		prefixStr = "[Info]"
		c.SetColor(colors.Cyan, colors.Black)
	case LogWarn:
		prefixStr = "[Warn]"
		c.SetColor(colors.Yellow, colors.Black)
	case LogError:
		prefixStr = "[Error]"
		c.SetColor(colors.Red, colors.Black)
	case LogFatal:
		prefixStr = "[Fatal]"
		c.SetColor(colors.Red, colors.White)
	}
	s := fmt.Sprintf(prefixStr+format, v...)
	out := fmt.Sprintf(format, v...)
	ChangDay() //! 跨天则生成新文件
	loger.Output(3, s)
	if true == terminalOutput {
		c.Printf(prefixStr)
		std.Output(3, out)
	}

	outputLock.Unlock()
}

func Print(format string, v ...interface{}) { //! Print
	fmt.Printf(format, v...)
}

func Debug(format string, v ...interface{}) { //! 调试
	Output(LogDebug, format, v...)
}

func Info(format string, v ...interface{}) { //! 信息
	Output(LogInfo, format, v...)
}

func Warn(format string, v ...interface{}) { //! 警告
	Output(LogWarn, format, v...)
}

func Error(format string, v ...interface{}) { //! 错误
	Output(LogError, format, v...)
}

func Fatal(format string, v ...interface{}) { //! 致命错误,使用会造成服务器退出,慎用!!
	Output(LogFatal, format, v...)
	os.Exit(1)
}

func ChangDay() { //! 跨天则生成新文件
	now := time.Now()
	currentDay := now.Day()
	timeformat := fileName + "-20060102.log"
	fileName := logPath + "/" + now.Format(timeformat)
	_, err := os.Stat(fileName)
	if currentDay == currentDay && nil == err && loger != nil {
		return
	}
	currentDay = currentDay
	logfile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0)
	if err != nil {
		log.Println("open log file fail!", err.Error())
	}
	loger = log.New(logfile, "", log.Ltime|log.Lmicroseconds|log.Lshortfile)
}

//! 创建日志管理器  日志路径  日志警示最小等级  终端同步输出
func InitLoger(Path string, MinLevel int, Output bool, Name string) {
	os.Mkdir(Path, 7777) //! 创建log目录,如果存在则忽略  log.Ldate|
	logPath = Path
	LogMinLevel = MinLevel
	terminalOutput = Output
	fileName = Name
	std = log.New(os.Stderr, "", log.Ltime|log.Lmicroseconds|log.Lshortfile)
}
