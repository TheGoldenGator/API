package middleware

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buf        bytes.Buffer
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{ResponseWriter: w}
}

func (w *LogResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LogResponseWriter) Write(body []byte) (int, error) {
	w.buf.Write(body)
	return w.ResponseWriter.Write(body)
}

type LogMiddleware struct {
	logger *log.Logger
}

func NewLogMiddleware(logger *log.Logger) *LogMiddleware {
	return &LogMiddleware{logger: logger}
}

func (m *LogMiddleware) Func() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			logRespWriter := NewLogResponseWriter(w)
			next.ServeHTTP(logRespWriter, r)

			//requestQ := "INSERT INTO requests (timestamp, processing_time, agent, body, http_version, method, remote_address, url, response) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);"
			//ct := utils

			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			//bodyStr := buf.String()

			// TODO: When serving an HTML file it tries to log the whole body as a request. Find a way to ignore certain paths from sending too much data.
			//database.SQLQuery(requestQ, ct, time.Since(startTime).String(), r.UserAgent(), bodyStr, r.Proto, r.Method, r.RemoteAddr, r.RequestURI, logRespWriter.buf.String())

			m.logger.Printf(
				"origin=%s duration=%s status=%d",
				r.RemoteAddr,
				time.Since(startTime).String(),
				logRespWriter.statusCode)
		})
	}
}
