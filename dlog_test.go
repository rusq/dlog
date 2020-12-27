// Package log is a wrapper around standard go logger that adds the debug output
// functions.
package dlog

import "testing"

func TestSetDebug(t *testing.T) {
	type args struct {
		b bool
	}
	tests := []struct {
		name      string
		args      args
		wantDebug bool
	}{
		{"yes debug", args{true}, true},
		{"no debug", args{false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prevVal := std.debug
			defer func() { std.debug = prevVal }()
			SetDebug(tt.args.b)
			if std.debug != tt.wantDebug {
				t.Errorf("want debug: %v, got debug: %v", tt.wantDebug, std.debug)
			}
		})
	}
}
