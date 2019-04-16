package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var (
	info  = log.New(os.Stdout, "\033[36mINFO\033[00m  ", log.LstdFlags)
	warn  = log.New(os.Stdout, "\033[33mWARN\033[00m  ", log.LstdFlags)
	error = log.New(os.Stderr, "\033[31mERROR\033[00m ", log.LstdFlags)
)

const (
	infoCall int = 1 + iota
	warnCall
	errorCall
)

// Info calls Output to print to the info logger.
// Arguments are handled in the manner of fmt.Print.
func Info(v ...interface{}) {
	v = append([]interface{}{Caller(infoCall)}, v...)
	info.Print(v...)
}

// Infof calls Output to print to the info logger.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, v ...interface{}) {
	v = append([]interface{}{Caller(infoCall)}, v...)
	format = addCallerFormat(format)
	info.Printf(format, v...)
}

// Warn calls Output to print to the warn logger.
// Arguments are handled in the manner of fmt.Print.
func Warn(v ...interface{}) {
	v = append([]interface{}{Caller(warnCall)}, v...)
	warn.Print(v...)
}

// Warnf calls Output to print to the warn logger.
// Arguments are handled in the manner of fmt.Printf.
func Warnf(format string, v ...interface{}) {
	v = append([]interface{}{Caller(warnCall)}, v...)
	format = addCallerFormat(format)
	warn.Printf(format, v...)
}

// Error calls Output to print to the error logger.
// Arguments are handled in the manner of fmt.Print.
func Error(v ...interface{}) {
	v = append([]interface{}{Caller(errorCall)}, v...)
	error.Print(v...)
}

// Errorf calls Output to print to the error logger.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	v = append([]interface{}{Caller(errorCall)}, v...)
	format = addCallerFormat(format)
	error.Printf(format, v...)
}

// ErrorWithCaller calls Output to print to the error logger.
// Arguments are handled in the manner of fmt.Print.
func ErrorWithCaller(caller string, v ...interface{}) {
	v = append([]interface{}{caller}, v...)
	error.Print(v...)
}

// ErrorfWithCaller calls Output to print to the error logger.
// Arguments are handled in the manner of fmt.Printf.
func ErrorfWithCaller(caller string, format string, v ...interface{}) {
	v = append([]interface{}{caller}, v...)
	format = addCallerFormat(format)
	error.Printf(format, v...)
}

// Panic is equivalent to Error() followed by a call to panic().
func Panic(v ...interface{}) {
	v = append([]interface{}{Caller(errorCall)}, v...)
	error.Print(v...)
	panic(fmt.Sprint(v...))
}

// Panicf is equivalent to Errorf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	v = append([]interface{}{Caller(errorCall)}, v...)
	format = addCallerFormat(format)
	error.Printf(format, v...)
	panic(fmt.Sprintf(format, v...))
}

func Caller(c int) string {
	_, f, l, _ := runtime.Caller(2)
	switch c {
	case infoCall:
		return fmt.Sprintf("%v:%v: \n  \033[36m>>\033[00m  ", f, l)
	case warnCall:
		return fmt.Sprintf("%v:%v: \n  \033[33m>>\033[00m  ", f, l)
	case errorCall:
		return fmt.Sprintf("%v:%v: \n  \033[31m>>\033[00m  ", f, l)
	}
	return fmt.Sprintf("%v:%v: \n  \033[36m>>\033[00m  ", f, l)
}

func addCallerFormat(format string) string {
	return "%v" + format
}
