// Package log is a wrapper around standard go logger that adds the debug output
// functions.
package dlog

import (
	"context"
	"fmt"
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

type key int

var loggerKey key

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
		l.Output(2, fmt.Sprint(v...))
	}
}

func (l *Logger) Debugln(v ...interface{}) {
	if l.Logger == nil {
		l.Logger = defaultLogger()
	}
	if l.debug {
		l.Output(2, fmt.Sprintln(v...))
	}
}

func (l *Logger) Debugf(format string, a ...interface{}) {
	if l.Logger == nil {
		l.Logger = defaultLogger()
	}
	if l.debug {
		l.Output(2, fmt.Sprintf(format, a...))
	}
}

// NewContext returns a new Context that has logger attached.
func NewContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext returns the Logger value stored in ctx, if any.  If no Logger
// is present, it returns the standard logger instance.
func FromContext(ctx context.Context) *Logger {
	l, ok := ctx.Value(loggerKey).(*Logger)
	if l == nil || !ok {
		l = std
	}
	return l
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
	std.Output(2, fmt.Sprintf(format, v...))
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	std.Output(2, fmt.Sprintln(v...))
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	std.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	std.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	std.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to Print() followed by a call to panic().
func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Output(2, s)
	panic(s)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Output(2, s)
	panic(s)
}

// Panicln is equivalent to Println() followed by a call to panic().
func (l *Logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.Output(2, s)
	panic(s)
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
