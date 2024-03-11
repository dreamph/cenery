package echo

import (
	"bytes"
	"cenery"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
)

type response struct {
	resp *echo.Response
}

func NewResponse(req *echo.Response) cenery.Response {

	e := echo.New()
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	}))

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

func (h *response) GetHeaderBytes(key string) []byte {
	return []byte(h.resp.Header().Get(key))
}

func (h *response) GetHeaderString(key string) string {
	return h.resp.Header().Get(key)
}

func (h *response) SetHeader(key string, val string) {
	h.resp.Header().Set(key, val)
}

func (h *response) AddHeader(key string, val string) {
	h.resp.Header().Add(key, val)
}
