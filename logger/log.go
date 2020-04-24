package logger

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	countLogger *log.Logger
	loggerList  []*log.Logger
	//debugLogger *log.Logger
	//errorLogger *log.Logger
	//fatalLogger *log.Logger
	level   LogLevel
	isDebug bool
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	ERROR
	FATAL
)

//type TimeTask struct{}
//
//func (task *TimeTask) Run() {
//	create()
//	clean()
//}
//func clean() {
//	debugFilePrefix := getCurrentDirectory() + "/logfile/"
//	date := time.Now().AddDate(0, 0, -2).Format("20060102")
//	debugFile := debugFilePrefix + date + "_debug.txt"
//	if exist, err := isExists(debugFile); err == nil && exist {
//		err := os.Remove(debugFile)
//		if err != nil {
//			log.Println("debugFile clean failed", err)
//		}
//	}
//}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		//log.Fatal(err)
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//func create() {
//	logFilePrefix := getCurrentDirectory() + "/logfile/"
//	name := time.Now().Format("20060102")
//
//	CountLogFile, err := initLogFile(logFilePrefix+"count/", name+"count.txt")
//	if err != nil {
//		log.Println("countLogFile created failed", err)
//	}
//
//	InfoLogFile, err := initLogFile(logFilePrefix, "info_"+name+".txt")
//	if err != nil {
//		log.Println("infoLogFile created failed", err)
//	}
//
//	DebugLogFile, err := initLogFile(logFilePrefix, "debug_"+name+".txt")
//	if err != nil {
//		log.Println("debugLogFile created failed")
//	}
//
//	ErrorLogFile, err := initLogFile(logFilePrefix, "error_"+name+".txt")
//	if err != nil {
//		log.Println("errorLogFile created failed")
//	}
//
//	FatalLogFile, err := initLogFile(logFilePrefix, "fatal_"+name+".txt")
//	if err != nil {
//		log.Println("fatalLogFile created failed")
//	}
//
//	countLogger = CountLogFile
//	if !isDebug {
//		InfoLog := log.New(InfoLogFile, "[Info]", log.LstdFlags)
//		infoLogger = InfoLog
//		DebugLog := log.New(DebugLogFile, "[Debug]", log.LstdFlags)
//		debugLogger = DebugLog
//		ErrorLog := log.New(ErrorLogFile, "[Error]", log.LstdFlags)
//		errorLogger = ErrorLog
//		FatalLog := log.New(FatalLogFile, "[Fatal]", log.LstdFlags)
//		fatalLogger = FatalLog
//	} else {
//		InfoLogToStdout := log.New(os.Stdout, "[Info]", log.LstdFlags)
//		infoLogger = InfoLogToStdout
//		DebugLogToStdout := log.New(os.Stdout, "[Debug]", log.LstdFlags)
//		debugLogger = DebugLogToStdout
//		ErrorLogToStdout := log.New(os.Stdout, "[Error]", log.LstdFlags)
//		errorLogger = ErrorLogToStdout
//		FatalLogToStdout := log.New(os.Stdout, "[Fatal]", log.LstdFlags)
//		fatalLogger = FatalLogToStdout
//	}
//}

func InitLogConfig(logLevel LogLevel, debug bool) {
	InitLogConfigWithPrefix(logLevel, debug, "")
}

func InitLogConfigWithPrefix(logLevel LogLevel, debug bool, filePrefix string) {
	level = logLevel
	isDebug = debug
	loggerList = make([]*log.Logger, 4)
	var logFilePrefix string
	if len(filePrefix) > 0 {
		logFilePrefix = getCurrentDirectory() + "/logfile/" + filePrefix + "_"
	} else {
		logFilePrefix = getCurrentDirectory() + "/logfile/"
	}
	//name := time.Now().Format("20060102")

	InfoLogFile := &lumberjack.Logger{
		Filename:   logFilePrefix + "info.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     1,     //days
		Compress:   false, // disabled by default
		LocalTime:  true,
	}

	DebugLogFile := &lumberjack.Logger{
		Filename:   logFilePrefix + "debug.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     1,     //days
		Compress:   false, // disabled by default
		LocalTime:  true,
	}

	ErrorLogFile := &lumberjack.Logger{
		Filename:  logFilePrefix + "error.log",
		MaxSize:   500,   // megabytes
		MaxAge:    1,     //days
		Compress:  false, // disabled by default
		LocalTime: true,
	}

	FatalLogFile := &lumberjack.Logger{
		Filename:  logFilePrefix + "fatal.log",
		MaxSize:   500,   // megabytes
		MaxAge:    1,     //days
		Compress:  false, // disabled by default
		LocalTime: true,
	}

	CountLogFile := &lumberjack.Logger{
		Filename:  logFilePrefix + "count/" + "count.log",
		MaxSize:   20,
		LocalTime: true,
	}
	//CountLogFile, err := initLogFile(logFilePrefix+"count/", name+"count.txt")
	//if err != nil {
	//	log.Println("countLogFile created failed", err)
	//}

	countLogger = log.New(CountLogFile, "", -1)
	if !isDebug {
		loggerList[DEBUG] = log.New(DebugLogFile, "[Debug]", log.LstdFlags)
		//debugLogger = DebugLog
		loggerList[INFO] = log.New(InfoLogFile, "[Info]", log.LstdFlags)
		//infoLogger = InfoLog
		loggerList[ERROR] = log.New(ErrorLogFile, "[Error]", log.LstdFlags)
		//errorLogger = ErrorLog
		loggerList[FATAL] = log.New(FatalLogFile, "[Fatal]", log.LstdFlags)
		//fatalLogger = FatalLog
	} else {
		loggerList[DEBUG] = log.New(os.Stdout, "[Debug]", log.LstdFlags)
		//debugLogger = DebugLogToStdout
		loggerList[INFO] = log.New(os.Stdout, "[Info]", log.LstdFlags)
		//infoLogger = InfoLogToStdout
		loggerList[ERROR] = log.New(os.Stdout, "[Error]", log.LstdFlags)
		//errorLogger = ErrorLogToStdout
		loggerList[FATAL] = log.New(os.Stdout, "[Fatal]", log.LstdFlags)
		//fatalLogger = FatalLogToStdout
	}
}

func initLogFile(dir string, name string) (*os.File, error) {
	if exist, err := isExists(dir); err == nil {
		if exist {
			return os.OpenFile(dir+name, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
		} else {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return nil, err
			}
			return os.OpenFile(dir+name, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
		}
	} else {
		return nil, err
	}
}

func isExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

//Info
func Infof(format string, v ...interface{}) {
	if level == DEBUG || level == INFO {
		file, line := getFileLine()
		v = append(v, file)
		v = append(v, line)
		loggerList[INFO].Printf(format+" %s %v\n", v...)
	}
}

func Info(v ...interface{}) {
	if level == DEBUG || level == INFO {
		file, line := getFileLine()
		v = append(v, file)
		v = append(v, line)
		loggerList[INFO].Print(v...)
	}
}

func Infoln(v ...interface{}) {
	if level == DEBUG || level == INFO {
		file, line := getFileLine()
		v = append(v, file)
		v = append(v, line)
		loggerList[INFO].Println(v...)
	}
}

//Debug
func Debugf(format string, v ...interface{}) {
	if level == DEBUG {
		file, line := getFileLine()
		v = append(v, file)
		v = append(v, line)
		loggerList[DEBUG].Printf(format+" %s %v\n", v...)
	}
}

func Debug(v ...interface{}) {
	if level == DEBUG {
		file, line := getFileLine()
		v = append(v, file)
		v = append(v, line)
		loggerList[DEBUG].Print(v...)
	}

}

func Debugln(v ...interface{}) {
	if level == DEBUG {
		file, line := getFileLine()
		v = append(v, file)
		v = append(v, line)
		loggerList[DEBUG].Println(v...)
	}
}

//Error
func Errorf(format string, v ...interface{}) {
	if level == DEBUG || level == INFO || level == ERROR {
		file, line := getFileLine()
		v = append(v, file)
		v = append(v, line)
		loggerList[ERROR].Printf(format+" %s %v\n", v...)
	}
}

func Error(v ...interface{}) {
	if level == DEBUG || level == INFO || level == ERROR {
		file, line := getFileLine()
		v = append(v, file)
		v = append(v, line)
		loggerList[ERROR].Print(v...)
	}
}

func Errorln(v ...interface{}) {
	if level == DEBUG || level == INFO || level == ERROR {
		file, line := getFileLine()
		v = append(v, file)
		v = append(v, line)
		loggerList[ERROR].Println(v...)
	}
}

//Fatal
func Fatalf(format string, v ...interface{}) {
	if level == DEBUG || level == INFO || level == ERROR || level == FATAL {
		file, line := getFileLine()
		v = append(v, file)
		v = append(v, line)
		loggerList[FATAL].Fatalf(format+" %s %v\n", v...)
	}
}

func Fatal(v ...interface{}) {
	if level == DEBUG || level == INFO || level == ERROR || level == FATAL {
		file, line := getFileLine()
		v = append(v, file)
		v = append(v, line)
		loggerList[FATAL].Fatal(v...)
	}
}

func Fatalln(v ...interface{}) {
	if level == DEBUG || level == INFO || level == ERROR || level == FATAL {
		file, line := getFileLine()
		v = append(v, file)
		v = append(v, line)
		loggerList[FATAL].Fatalln(v...)
	}
}

func Count(methods []byte) {
	if countLogger != nil {
		countLogger.Println(string(methods))
	}
}

func getFileLine() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		fileSlice := strings.Split(file, "/")
		realFile := fileSlice[len(fileSlice)-1]
		return realFile, line
	} else {
		return "找不到路径", 0
	}
}

func GetLogger(level LogLevel) *log.Logger {
	return loggerList[level]
}
