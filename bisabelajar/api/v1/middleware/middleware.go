package middlewarev1

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const (
	requestIDKey     contextKey = "requestID"
	requestTimestamp contextKey = "timestamp"
)

type RequestDetails struct {
	RequestID string
	Timestamp string
	URL       string
	Method    string
}

func RequestIDAndTimestampMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		timestamp := r.Header.Get("X-Timestamp")
		if timestamp == "" {
			timestamp = time.Now().Format(time.RFC3339)
		}

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		ctx = context.WithValue(ctx, requestTimestamp, timestamp)

		w.Header().Set("X-Request-ID", requestID)
		w.Header().Set("X-Timestamp", timestamp)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestDetails(r *http.Request) RequestDetails {
	ctx := r.Context()
	return RequestDetails{
		RequestID: ctx.Value(requestIDKey).(string),
		Timestamp: ctx.Value(requestTimestamp).(string),
		URL:       r.URL.String(),
		Method:    r.Method,
	}
}
