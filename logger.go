package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/goroute/logger/terminal"
	"github.com/goroute/route"
)

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

// FormatType defines logs format mode.
type FormatType string

const (
	// FormatTypeText defines text format mode.
	FormatTypeText FormatType = "text"
	// FormatTypeJSON defines json format mode.
	FormatTypeJSON FormatType = "json"
)

// Options defines the options for logger middleware.
type Options struct {
	// Skipper defines a function to skip middleware.
	Skipper route.Skipper

	// Output defines logs output writter. Default is os.Stderr.
	Output io.Writer

	// Format defines logs format mode. Default is text.
	Format FormatType

	bufferPool *sync.Pool
	isTerminal bool
}

// Option defines option func.
type Option func(*Options)

// GetDefaultOptions returns default options.
func GetDefaultOptions() Options {
	return Options{
		Skipper: route.DefaultSkipper,
		Output:  os.Stderr,
		Format:  FormatTypeText,
	}
}

// Skipper sets skipper option.
func Skipper(skipper route.Skipper) Option {
	return func(o *Options) {
		o.Skipper = skipper
	}
}

// Format sets logs format option.
func Format(format FormatType) Option {
	return func(o *Options) {
		o.Format = format
	}
}

// Output sets logs output.
func Output(w io.Writer) Option {
	return func(o *Options) {
		o.Output = w
	}
}

// New returns a middleware which logs requests
func New(options ...Option) route.MiddlewareFunc {
	// Apply options.
	opts := GetDefaultOptions()
	for _, opt := range options {
		opt(&opts)
	}

	opts.bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	opts.isTerminal = terminal.IsTerminal(int(os.Stdout.Fd()))

	return func(c route.Context, next route.HandlerFunc) error {
		if opts.Skipper(c) {
			return next(c)
		}

		req := c.Request()
		res := c.Response()
		start := time.Now()
		err := next(c)
		if err != nil {
			c.Error(err)
		}
		stop := time.Now()
		latency := stop.Sub(start).String()

		// Get write buffer from sync pool.
		b := opts.bufferPool.Get().(*bytes.Buffer)
		b.Reset()
		defer opts.bufferPool.Put(b)

		var errMsg string
		var color int
		if err != nil && res.Status >= 500 {
			errMsg = fmt.Sprintf("%v", err)
			// 5xx
			color = colorRed
		} else {
			if res.Status >= 200 && res.Status < 300 {
				// 2xx
				color = colorBlue
			} else if res.Status >= 500 {
				// 5xx
				color = colorRed
			} else if res.Status >= 300 && res.Status < 500 {
				// 3xx and 4xx
				color = colorYellow
			} else {
				// 1xx
				color = colorGray
			}
		}

		if opts.Format == FormatTypeText {
			if opts.isTerminal {
				fmt.Fprintf(b, "\x1b[%dm%d\x1b[0m method=%s path=%s latency=%s", color, res.Status, req.Method, req.URL.Path, latency)
				if errMsg != "" {
					fmt.Fprintf(b, " \x1b[%dmerr=%s\x1b[0m", color, errMsg)
				}
			} else {
				fmt.Fprintf(b, "%d method=%s path=%s latency=%s %s", res.Status, req.Method, req.URL.Path, latency, errMsg)
				if errMsg != "" {
					fmt.Fprintf(b, " err=%s", errMsg)
				}
			}
			b.WriteByte('\n')
			opts.Output.Write(b.Bytes())
		} else {
			encoder := json.NewEncoder(b)
			encoder.Encode(struct {
				Status  int    `json:"status"`
				Method  string `json:"method"`
				Path    string `json:"path"`
				Latency string `json:"latency"`
				Error   string `json:"error,omitempty"`
			}{
				Status:  res.Status,
				Method:  req.Method,
				Path:    req.URL.Path,
				Latency: latency,
				Error:   errMsg,
			})
			opts.Output.Write(b.Bytes())
		}

		b = nil
		return nil
	}
}
