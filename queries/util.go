package queries

import (
	"net/http"
	"time"
)

func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

func GetRFCTimestamp() string {
	now := time.Now()
	rfc := now.Format(time.RFC3339)
	return rfc
}
