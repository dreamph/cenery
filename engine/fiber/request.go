package fiber

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
	return h.req.Body()
}

func (h *request) SetBody(data []byte) {
	h.req.SetBody(data)
}

func (h *request) BodyStream() io.ReadCloser {
	if stream := h.req.BodyStream(); stream != nil {
		return io.NopCloser(stream)
	}
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
