package slog

import (
	"fmt"
	"io"
	"io/ioutil"
	"goon3/lib/kscan/lib/chinese"
	"log"
	"os"
	"runtime"
	"strings"
)

var this logger
var splitStr = "> "

type logger struct {
	info     *log.Logger
	warning  *log.Logger
	debug    *log.Logger
	data     *log.Logger
	error    *log.Logger
	encoding string
	//fooLine *log.Logger
}

func Init(Debug bool, encoding string) {
	this.info = log.New(os.Stdout, "\r[+]", log.Ldate|log.Ltime)
	this.warning = log.New(os.Stdout, "\r[*]", log.Ldate|log.Ltime)
	this.error = log.New(io.MultiWriter(os.Stderr), "\rError:", 0)
	this.data = log.New(os.Stdout, "\r", 0)
	if Debug {
		this.debug = log.New(os.Stdout, "\r[-]", log.Ldate|log.Ltime)
	} else {
		this.debug = log.New(ioutil.Discard, "\r[-]", log.Ldate|log.Ltime)
	}

	this.encoding = encoding
	//this.fooline = log.New(os.Stdout, "[*]", 0)
	//infoFile,err:=os.OpenFile("/data/service_logs/info.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	//warnFile,err:=os.OpenFile("/data/service_logs/warn.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	//errFile,err:=os.OpenFile("/data/service_logs/errors.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	//
	//if infoFile!=nil || warnFile != nil || err!=nil{
	//	log.Fatalln("打开日志文件失败：",err)
	//}
	//Info = log.New(os.Stdout, "[*]", log.Ldate|log.Ltime)
	//Warning = log.New(os.Stdout, "[*]", log.Ldate|log.Ltime)
	//Error = log.New(io.MultiWriter(os.Stderr,errFile),"Error:",log.Ldate | log.Ltime | log.Lshortfile)
	//Info = log.New(io.MultiWriter(os.Stderr,infoFile),"Info:",log.Ldate | log.Ltime | log.Lshortfile)
	//Warning = log.New(io.MultiWriter(os.Stderr,warnFile),"Warning:",log.Ldate | log.Ltime | log.Lshortfile)
	//Error = log.New(io.MultiWriter(os.Stderr,errFile),"Error:",log.Ldate | log.Ltime | log.Lshortfile)
}

func (t *logger) Data(s string) {
	t.data.Print(s)
}

func (t *logger) Info(s string) {
	t.info.Print(splitStr, s)
}

func (t *logger) Error(s string) {
	t.error.Print(s)
	os.Exit(0)
}

func (t *logger) Warning(s string) {
	t.warning.Print(splitStr, s)
}

func (t *logger) Debug(s string) {
	if debugFilter(s) {
		return
	}
	_, file, line, _ := runtime.Caller(3)
	file = file[strings.LastIndex(file, "/")+1:]
	t.debug.Printf("%s%s(%d) %s", splitStr, file, line, s)
}

func (t *logger) DoPrint(logType string, logStr string) {
	if this.encoding == "gb2312" {
		logStr = chinese.ToGBK(logStr)
	} else {
		logStr = chinese.ToUTF8(logStr)
	}
	switch logType {
	case "Debug":
		t.Debug(logStr)
	case "Info":
		t.Info(logStr)
	case "Data":
		t.Data(logStr)
	case "Warning":
		t.Warning(logStr)
	case "Error":
		t.Error(logStr)
	}

}

//func (t *logger) Error(s string) {
//	_, file, line, _ := runtime.Caller(2)
//	file = file[strings.LastIndex(file, "/")+1:]
//	t.error.Printf("%s%s(%d) %s", splitStr, file, line, s)
//	os.Exit(0)
//}

//func (t *logger) Errorf(format string, v ...interface{}) {
//	_, file, line, _ := runtime.Caller(2)
//	file = file[strings.LastIndex(file, "/")+1:]
//	format = fmt.Sprintf("%s%s(%d) %s", splitStr, file, line, format)
//	t.error.Printf(format, v...)
//	os.Exit(0)
//}
func Error(s ...interface{}) {
	this.DoPrint("Error", fmt.Sprint(s...))
}

func Info(s ...interface{}) {
	this.DoPrint("Info", fmt.Sprint(s...))
}

func Infof(format string, v ...interface{}) {
	this.DoPrint("Info", fmt.Sprintf(format, v...))
}

func Warning(s ...interface{}) {
	this.DoPrint("Warning", fmt.Sprint(s...))
}

func Warningf(format string, v ...interface{}) {
	this.DoPrint("Warning", fmt.Sprintf(format, v...))
}

func Debug(s ...interface{}) {
	this.DoPrint("Debug", fmt.Sprint(s...))
}

func Debugf(format string, v ...interface{}) {
	this.DoPrint("Debug", fmt.Sprintf(format, v...))
}

func Data(v ...interface{}) {
	this.DoPrint("Data", fmt.Sprint(v...))
}

//func Error(s string) {
//	this.Error(s)
//}

//func Errorf(format string, v ...interface{}) {
//	this.Errorf(format, v...)
//}

func debugFilter(s string) bool {
	//Debug 过滤器
	if strings.Contains(s, "too many") { //发现存在线程过高错误
		Error("当前线程过高，请降低线程!或者请执行\"ulimit -n 50000\"命令放开操作系统限制,MAC系统可能还需要执行：\"launchctl limit maxfiles 50000 50000\"")
	}
	if strings.Contains(s, "STEP1:CONNECT") {
		return true
	}
	return false
}
