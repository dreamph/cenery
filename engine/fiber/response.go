package fiber

import (
	"github.com/dreamph/cenery"
	"github.com/valyala/fasthttp"
)

type response struct {
	resp *fasthttp.Response
}

func NewResponse(req *fasthttp.Response) cenery.Response {
	return &response{resp: req}
}

func (h *response) Body() []byte {
	return h.resp.Body()
}

func (h *response) SetBody(data []byte) {
	h.resp.SetBody(data)
}

func (h *response) GetHeader(key string) string {
	return string(h.resp.Header.Peek(key))
}

func (h *response) SetHeader(key string, val string) {
	h.resp.Header.Add(key, val)
}

func (h *response) AddHeader(key string, val string) {
	h.resp.Header.Add(key, val)
}
