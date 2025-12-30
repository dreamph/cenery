package chi

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"sync/atomic"

	"github.com/dreamph/cenery"
)

type responseBodyWriter struct {
	http.ResponseWriter
	writer io.Writer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func (w *responseBodyWriter) WriteString(s string) (int, error) {
	return io.WriteString(w.writer, s)
}

func (w *responseBodyWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *responseBodyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

func (w *responseBodyWriter) CloseNotify() <-chan bool {
	if cn, ok := w.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return make(chan bool)
}

func (w *responseBodyWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}

type response struct {
	resp    http.ResponseWriter
	resBody *bytes.Buffer
}

var captureResponseBody atomic.Bool

// EnableResponseCapture toggles response body capture (for logging/testing).
// When disabled, response bodies are not buffered to avoid extra allocs/writes.
func EnableResponseCapture(enabled bool) {
	captureResponseBody.Store(enabled)
}

func NewResponse(w http.ResponseWriter) (cenery.Response, http.ResponseWriter) {
	var resBodyBuffer *bytes.Buffer
	if captureResponseBody.Load() {
		resBodyBuffer = &bytes.Buffer{}
		writer := &responseBodyWriter{
			ResponseWriter: w,
			writer:         io.MultiWriter(w, resBodyBuffer),
		}
		w = writer
	}

	return &response{
		resp:    w,
		resBody: resBodyBuffer,
	}, w
}

func (h *response) Body() []byte {
	if h.resBody == nil {
		return nil
	}
	return h.resBody.Bytes()
}

func (h *response) SetBody(data []byte) {
	_, _ = h.resp.Write(data)
}

func (h *response) GetHeader(key string) string {
	return h.resp.Header().Get(key)
}

func (h *response) SetHeader(key string, val string) {
	h.resp.Header().Set(key, val)
}

func (h *response) AddHeader(key string, val string) {
	h.resp.Header().Add(key, val)
}
