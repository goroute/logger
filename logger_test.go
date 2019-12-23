package logger

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goroute/route"
	"github.com/stretchr/testify/assert"
)

func TestLoggerWithTextFormat(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		method   string
		path     string
		expected string
	}{
		{
			name:     "Test info status code",
			status:   http.StatusOK,
			method:   http.MethodGet,
			path:     "/ok",
			expected: "200 method=GET path=/ok latency=",
		},
		{
			name:     "Test gray status code",
			status:   http.StatusContinue,
			method:   http.MethodGet,
			path:     "/",
			expected: "100 method=GET path=/ latency=",
		},
		{
			name:     "Test warn status code",
			status:   http.StatusNotFound,
			method:   http.MethodPost,
			path:     "/",
			expected: "404 method=POST path=/ latency=",
		},
		{
			name:     "Test err status code",
			status:   http.StatusInternalServerError,
			method:   http.MethodDelete,
			path:     "/err",
			expected: "500 method=DELETE path=/err latency=",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			buf := new(bytes.Buffer)
			mux := route.NewServeMux()
			req := httptest.NewRequest(test.method, test.path, nil)
			rec := httptest.NewRecorder()
			c := mux.NewContext(req, rec)
			h := func(c route.Context) error {
				return c.String(test.status, "test")
			}
			mw := New(Output(buf), Format(FormatTypeText))

			mw(c, h)

			str := buf.String()
			assert.Contains(t, str, test.expected)
		})
	}
}

func TestLoggerWithJSONFormat(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		method   string
		path     string
		expected string
	}{
		{
			name:     "Test info status code",
			status:   http.StatusOK,
			method:   http.MethodGet,
			path:     "/ok",
			expected: `{"status":200,"method":"GET","path":"/ok","latency":"`,
		},
		{
			name:     "Test gray status code",
			status:   http.StatusContinue,
			method:   http.MethodGet,
			path:     "/",
			expected: `{"status":100,"method":"GET","path":"/","latency":"`,
		},
		{
			name:     "Test warn status code",
			status:   http.StatusNotFound,
			method:   http.MethodPost,
			path:     "/",
			expected: `{"status":404,"method":"POST","path":"/","latency":"`,
		},
		{
			name:     "Test err status code",
			status:   http.StatusInternalServerError,
			method:   http.MethodDelete,
			path:     "/err",
			expected: `{"status":500,"method":"DELETE","path":"/err","latency":"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			buf := new(bytes.Buffer)
			mux := route.NewServeMux()
			req := httptest.NewRequest(test.method, test.path, nil)
			rec := httptest.NewRecorder()
			c := mux.NewContext(req, rec)
			h := func(c route.Context) error {
				return c.String(test.status, "test")
			}
			mw := New(Output(buf), Format(FormatTypeJSON))

			mw(c, h)

			str := buf.String()
			fmt.Println("res", str)
			assert.Contains(t, str, test.expected)
		})
	}
}
