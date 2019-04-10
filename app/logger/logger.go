package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
)

var (
	info  = log.New(os.Stdout, "\033[36mINFO\033[00m  ", log.LstdFlags)
	warn  = log.New(os.Stdout, "\033[33mWARN\033[00m  ", log.LstdFlags)
	error = log.New(os.Stderr, "\033[31mERROR\033[00m ", log.LstdFlags)
)

// Info calls Output to print to the info logger.
// Arguments are handled in the manner of fmt.Print.
func Info(v ...interface{}) {
	v = append([]interface{}{Caller()}, v...)
	info.Print(v...)
}

// Infof calls Output to print to the info logger.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, v ...interface{}) {
	v = append([]interface{}{Caller()}, v...)
	format = addCallerFormat(format)
	info.Printf(format, v...)
}

// Warn calls Output to print to the warn logger.
// Arguments are handled in the manner of fmt.Print.
func Warn(v ...interface{}) {
	v = append([]interface{}{Caller()}, v...)
	warn.Print(v...)
}

// Warnf calls Output to print to the warn logger.
// Arguments are handled in the manner of fmt.Printf.
func Warnf(format string, v ...interface{}) {
	v = append([]interface{}{Caller()}, v...)
	format = addCallerFormat(format)
	warn.Printf(format, v...)
}

// Error calls Output to print to the error logger.
// Arguments are handled in the manner of fmt.Print.
func Error(v ...interface{}) {
	v = append([]interface{}{Caller()}, v...)
	error.Print(v...)
}

// Errorf calls Output to print to the error logger.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	v = append([]interface{}{Caller()}, v...)
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
	v = append([]interface{}{Caller()}, v...)
	error.Print(v...)
	panic(fmt.Sprint(v...))
}

// Panicf is equivalent to Errorf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	v = append([]interface{}{Caller()}, v...)
	format = addCallerFormat(format)
	error.Printf(format, v...)
	panic(fmt.Sprintf(format, v...))
}

func Caller() string {
	_, f, l, _ := runtime.Caller(2)
	splited := strings.Split(f, "/")
	var pathSegs []string
	var foundMono bool
	for _, s := range splited {
		if foundMono {
			pathSegs = append(pathSegs, s)
		}
		if s == "src" {
			foundMono = true
		}
	}
	// fmt.Println()
	// f = path.Base(f)

	return fmt.Sprintf("%v:%v: ", path.Join(pathSegs...), l)
}

func addCallerFormat(format string) string {
	return "%v" + format
}
