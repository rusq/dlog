// Package log is a wrapper around standard go logger that adds the debug output
// functions.
package dlog

import (
	"bytes"
	"context"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestLogger_SetDebug(t *testing.T) {
	type fields struct {
		Logger *log.Logger
		debug  bool
		//mu     sync.Mutex
	}
	type args struct {
		b bool
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantDebug bool
		wantFlags int
	}{
		{"set debug", fields{Logger: defaultLogger(), debug: false}, args{true}, true, log.LstdFlags + log.Lshortfile},
		{"reset debug", fields{Logger: defaultLogger(), debug: true}, args{false}, false, log.LstdFlags},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				Logger: tt.fields.Logger,
				debug:  tt.fields.debug,
			}
			l.SetDebug(tt.args.b)
			if l.debug != tt.wantDebug {
				t.Errorf("want debug: %v, got debug: %v", tt.wantDebug, l.debug)
			}
			if flags := l.Flags(); flags != tt.wantFlags {
				t.Errorf("want flags: %v, got flags: %v", tt.wantFlags, flags)
			}
		})
	}
}

func TestLogger_Debug(t *testing.T) {
	type fields struct {
		Logger *log.Logger
		debug  bool
	}
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantOutputRe string
	}{
		{"debug is on",
			fields{debug: true},
			args{v: []interface{}{"message1 ", "message2"}},
			`^.*message1\s+message2`,
		},
		{"debug is off",
			fields{debug: false},
			args{v: []interface{}{"message1 ", "message2"}},
			`^$`,
		},
		{"debug is on, prefix is set",
			fields{Logger: log.New(os.Stderr, "testxxx: ", log.LstdFlags), debug: true},
			args{v: []interface{}{"message1 ", "message2"}},
			`^testxxx: .*message1\s+message2$`,
		},
		{"debug is off, prefix is set",
			fields{Logger: log.New(os.Stderr, "testxxx: ", log.LstdFlags), debug: false},
			args{v: []interface{}{"message1 ", "message2"}},
			`^$`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				Logger: tt.fields.Logger,
				debug:  tt.fields.debug,
			}
			if l.Logger == nil {
				l.Logger = defaultLogger()
			}
			re, err := regexp.Compile(tt.wantOutputRe)
			if err != nil {
				t.Fatal(err)
			}
			for i, fn := range []func(arg ...interface{}){l.Debug, l.Debugln} {
				var buf bytes.Buffer
				l.SetOutput(&buf)
				fn(tt.args.v...)

				if !re.Match(bytes.TrimSpace(buf.Bytes())) {
					t.Errorf("output for fn: %d: mismatch: wantRE: %q, got: %q", i, tt.wantOutputRe, buf.String())
				}
			}
		})
	}
}

func TestLogger_Debugf(t *testing.T) {
	type fields struct {
		Logger *log.Logger
		debug  bool
	}
	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantOutputRe string
	}{
		{"debug is on",
			fields{debug: true},
			args{format: "%s%s", v: []interface{}{"message1 ", "message2"}},
			`^.*dlog_test\.go:.*message1\s+message2`,
		},
		{"debug is off",
			fields{debug: false},
			args{format: "%s%s", v: []interface{}{"message1 ", "message2"}},
			`^$`,
		},
		{"debug is on, prefix is set",
			fields{Logger: log.New(os.Stderr, "testxxx: ", log.LstdFlags), debug: true},
			args{format: "%s%s", v: []interface{}{"message1 ", "message2"}},
			`^testxxx: .*dlog_test\.go:.*message1\s+message2$`,
		},
		{"debug is off, prefix is set",
			fields{Logger: log.New(os.Stderr, "testxxx: ", log.LstdFlags), debug: false},
			args{format: "%s%s", v: []interface{}{"message1 ", "message2"}},
			`^$`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				Logger: tt.fields.Logger,
			}
			if l.Logger == nil {
				l.Logger = defaultLogger()
			}
			re, err := regexp.Compile(tt.wantOutputRe)
			if err != nil {
				t.Fatal(err)
			}
			l.SetDebug(tt.fields.debug)

			var buf bytes.Buffer
			l.SetOutput(&buf)
			l.Debugf(tt.args.format, tt.args.v...)

			if !re.Match(bytes.TrimSpace(buf.Bytes())) {
				t.Errorf("output mismatch: wantRE: %q, got: %q", tt.wantOutputRe, buf.String())
			}
		})
	}
}

func TestNewContext(t *testing.T) {
	var buf strings.Builder
	l := New(&buf, ">", log.LstdFlags, true)
	ctx := NewContext(context.Background(), l)

	lFromCtx, ok := ctx.Value(loggerKey).(*Logger)
	if lFromCtx == nil || !ok {
		t.Fatal("failed to get the logger from context")
	}

	const testVal = "TestNewContext"
	lFromCtx.Println(testVal)
	if !strings.Contains(buf.String(), testVal) {
		t.Fatalf("invalid logger: test value %v, not found in output %v", testVal, buf.String())
	}
}

func TestFromContext(t *testing.T) {
	custom := New(os.Stdout, ">", log.LstdFlags, true)
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *Logger
	}{
		{"standard logger", args{context.Background()}, std},
		{"custom logger", args{NewContext(context.Background(), custom)}, custom},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
