package gocache

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWrapHandler(t *testing.T) {
	type testResponse struct {
		Message string    `json:"message"`
		SentAt  time.Time `json:"sentAt"`
	}

	handlerHits := 0

	// set router and request handler
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// increment handler hits
		handlerHits++

		// set cookies
		http.SetCookie(w, &http.Cookie{
			Name:  "Test-Cookie",
			Value: "test-cookie",
		})
		// set status code
		w.WriteHeader(http.StatusTeapot)
		// set header
		w.Header().Set("X-Test", "test")
		// write body
		json.NewEncoder(w).Encode(testResponse{
			Message: "test message",
			SentAt:  time.Now(),
		})
	})

	// make request identifier
	identifier := NewRequestIdentifier(1*time.Minute, func(r *http.Request) string {
		// only if request method is GET
		if r.Method != http.MethodGet {
			return ""
		}

		id := r.Method
		id += r.RequestURI
		cookies := r.Cookies()
		for _, cookie := range cookies {
			id += cookie.String()
		}

		return id
	})

	// init cache
	cachedHandler, err := CacheHTTP(router, identifier, CacheConfig{
		ID: "HTTPResponse",
	})
	if err != nil {
		t.Error(err)
		return
	}

	// test request
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	resrec1 := httptest.NewRecorder()
	cachedHandler.ServeHTTP(resrec1, request)

	resrec2 := httptest.NewRecorder()
	cachedHandler.ServeHTTP(resrec2, request)

	if handlerHits > 1 {
		t.Error("handler should have only been hit once but was hit", handlerHits, "times instead")
		return
	}
}
