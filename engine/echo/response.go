package echo

import (
	"bytes"
	"github.com/dreamph/cenery"
	"github.com/labstack/echo/v4"
	"io"
)

type response struct {
	resp *echo.Response
}

func NewResponse(req *echo.Response) cenery.Response {
	return &response{resp: req}
}

func (h *response) Body() []byte {
	resBody := new(bytes.Buffer)
	_ = io.MultiWriter(h.resp.Writer, resBody)

	return resBody.Bytes()
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
