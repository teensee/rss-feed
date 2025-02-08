package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type traceIdKey string

const TraceIdHeader traceIdKey = "X-Trace-Id"

func (k traceIdKey) String() string {
	return string(k)
}

func TraceId(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		traceId := r.Header.Get(string(TraceIdHeader))
		if traceId == "" {
			traceId = uuid.NewString()
		}

		ctx = context.WithValue(ctx, TraceIdHeader, traceId)
		w.Header().Add(string(TraceIdHeader), traceId)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
