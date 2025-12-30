package gin

import (
	"bytes"
	"io"
	"sync/atomic"

	"github.com/dreamph/cenery"
	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	writer io.Writer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func (w *responseBodyWriter) WriteString(s string) (int, error) {
	return io.WriteString(w.writer, s)
}

type response struct {
	resp    gin.ResponseWriter
	resBody *bytes.Buffer
}

var captureResponseBody atomic.Bool

// EnableResponseCapture toggles response body capture (for logging/testing).
// When disabled, response bodies are not buffered to avoid extra allocs/writes.
func EnableResponseCapture(enabled bool) {
	captureResponseBody.Store(enabled)
}

func NewResponse(ctx *gin.Context) cenery.Response {
	resp := ctx.Writer
	var resBodyBuffer *bytes.Buffer
	if captureResponseBody.Load() {
		resBodyBuffer = &bytes.Buffer{}
		writer := &responseBodyWriter{
			ResponseWriter: resp,
			writer:         io.MultiWriter(resp, resBodyBuffer),
		}
		ctx.Writer = writer
		resp = writer
	}

	return &response{
		resp:    resp,
		resBody: resBodyBuffer,
	}
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
