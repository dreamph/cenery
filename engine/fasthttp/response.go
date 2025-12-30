package fasthttp

import (
	"bytes"
	"sync/atomic"

	"github.com/dreamph/cenery"
	"github.com/valyala/fasthttp"
)

type response struct {
	resp    *fasthttp.Response
	resBody *bytes.Buffer
}

var captureResponseBody atomic.Bool

// EnableResponseCapture toggles response body capture (for logging/testing).
// When disabled, response bodies are not buffered to avoid extra allocs/writes.
func EnableResponseCapture(enabled bool) {
	captureResponseBody.Store(enabled)
}

func NewResponse(ctx *fasthttp.RequestCtx) cenery.Response {
	resp := &response{resp: &ctx.Response}
	if captureResponseBody.Load() {
		resp.resBody = &bytes.Buffer{}
	}
	return resp
}

func (h *response) Body() []byte {
	if h.resBody != nil && h.resBody.Len() > 0 {
		return h.resBody.Bytes()
	}
	body := h.resp.Body()
	if len(body) == 0 {
		return nil
	}
	dup := make([]byte, len(body))
	copy(dup, body)
	return dup
}

func (h *response) SetBody(data []byte) {
	h.resp.SetBody(data)
	if h.resBody != nil {
		h.resBody.Reset()
		_, _ = h.resBody.Write(data)
	}
}

func (h *response) GetHeader(key string) string {
	return string(h.resp.Header.Peek(key))
}

func (h *response) SetHeader(key string, val string) {
	h.resp.Header.Set(key, val)
}

func (h *response) AddHeader(key string, val string) {
	h.resp.Header.Add(key, val)
}
