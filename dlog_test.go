// Package log is a wrapper around standard go logger that adds the debug output
// functions.
package dlog

import (
	"log"
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
