// Package log is a wrapper around standard go logger that adds the debug output
// functions.
package dlog

import (
	"io"
	"log"
	"os"
	"sync"
)

type Logger struct {
	*log.Logger
	debug bool
	mu    sync.Mutex
}

var std *Logger

func init() {
	isDebug := (os.Getenv("DEBUG") != "")
	flags := log.LstdFlags
	if isDebug {
		flags |= log.Lshortfile
	}
	std = New(os.Stderr, "", flags, isDebug)
}

func New(out io.Writer, prefix string, flag int, debug bool) *Logger {
	return &Logger{Logger: log.New(out, prefix, flag), debug: debug}
}

func (l *Logger) Debug(v ...interface{}) {
	if l.Logger == nil {
		l.Logger = defaultLogger()
	}
	if l.debug {
		l.Logger.Print(v...)
	}
}

func (l *Logger) Debugln(v ...interface{}) {
	if l.Logger == nil {
		l.Logger = defaultLogger()
	}
	if l.debug {
		l.Logger.Println(v...)
	}
}

func (l *Logger) Debugf(format string, a ...interface{}) {
	if l.Logger == nil {
		l.Logger = defaultLogger()
	}
	if l.debug {
		l.Logger.Printf(format, a...)
	}
}

// SetOutput sets the output destination for the standard logger.
func SetOutput(w io.Writer) {
	std.Logger.SetOutput(w)
}

// Flags returns the output flags for the standard logger.
// The flag bits are Ldate, Ltime, and so on.
func Flags() int {
	return std.Flags()
}

// SetFlags sets the output flags for the standard logger.
// The flag bits are Ldate, Ltime, and so on.
func SetFlags(flag int) {
	std.SetFlags(flag)
}

// Prefix returns the output prefix for the standard logger.
func Prefix() string {
	return std.Prefix()
}

// SetPrefix sets the output prefix for the standard logger.
func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

// SetDebug sets/resets the debugging output.
func SetDebug(b bool) {
	std.SetDebug(b)
}

// SetDebug sets/resets the debugging output.
func (l *Logger) SetDebug(b bool) {
	if l.Logger == nil {
		l.Logger = defaultLogger()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.debug = b
	if b {
		l.SetFlags(l.Flags() | log.Lshortfile)
	} else {
		l.SetFlags(l.Flags() &^ (1 << log.Lshortfile))
	}
}

func defaultLogger() *log.Logger {
	return log.New(os.Stderr, "", log.LstdFlags)
}

// Writer returns the output destination for the standard logger.
func Writer() io.Writer {
	return std.Writer()
}

// These functions write to the standard logger.

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	std.Print(v...)
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	std.Printf(format, v...)
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	std.Println(v...)
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	std.Fatal(v...)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	std.Fatalf(format, v...)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	std.Fatalln(v...)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	std.Panic(v...)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	std.Panicf(format, v...)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...interface{}) {
	std.Panicln(v...)
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is the count of the number of
// frames to skip when computing the file name and line number
// if Llongfile or Lshortfile is set; a value of 1 will print the details
// for the caller of Output.
func Output(calldepth int, s string) error {
	return std.Output(calldepth+1, s) // +1 for this frame.
}

func Debug(v ...interface{}) {
	std.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}

func Debugln(v ...interface{}) {
	std.Debugln(v...)
}
