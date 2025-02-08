package middleware

import (
	"testing"
)

func TestTraceId(t *testing.T) {
	// todo: add test case
}

func Test_traceIdKey_String(t *testing.T) {
	tests := []struct {
		name string
		k    traceIdKey
		want string
	}{
		{
			name: "traceId header must be 'X-Trace-Id'",
			k:    TraceIdHeader,
			want: "X-Trace-Id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
