package fiber

import (
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

func (h *request) GetHeaderBytes(key string) []byte {
	return h.req.Header.Peek(key)
}

func (h *request) GetHeaderString(key string) string {
	return string(h.GetHeaderBytes(key))
}

func (h *request) SetHeader(key string, val string) {
	h.req.Header.Set(key, val)
}

func (h *request) AddHeader(key string, val string) {
	h.req.Header.Add(key, val)
}
