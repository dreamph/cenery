package fasthttp

import (
	"bytes"
	"io"

	"github.com/dreamph/cenery"
	"github.com/valyala/fasthttp"
)

type request struct {
	req *fasthttp.Request
}

func NewRequest(req *fasthttp.Request) cenery.Request {
	return &request{req: req}
}

func (h *request) Body() []byte {
	body := h.req.Body()
	if len(body) == 0 {
		return nil
	}
	dup := make([]byte, len(body))
	copy(dup, body)
	return dup
}

func (h *request) SetBody(data []byte) {
	h.req.SetBody(data)
}

func (h *request) BodyStream() io.ReadCloser {
	return io.NopCloser(bytes.NewReader(h.req.Body()))
}

func (h *request) GetHeader(key string) string {
	return string(h.req.Header.Peek(key))
}

func (h *request) SetHeader(key string, val string) {
	h.req.Header.Set(key, val)
}

func (h *request) AddHeader(key string, val string) {
	h.req.Header.Add(key, val)
}
