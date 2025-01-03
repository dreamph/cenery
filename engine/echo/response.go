package echo

import (
	"bytes"
	"io"
	"net/http"

	"github.com/dreamph/cenery"
	"github.com/labstack/echo/v4"
)

type responseBodyWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	// Write to the original ResponseWriter
	return w.Writer.Write(b)
}

type response struct {
	resp    *echo.Response
	resBody *bytes.Buffer
}

func NewResponse(resp *echo.Response) cenery.Response {
	resBodyBuffer := &bytes.Buffer{}
	writer := &responseBodyWriter{
		Writer:         io.MultiWriter(resp.Writer, resBodyBuffer),
		ResponseWriter: resp.Writer,
	}
	resp.Writer = writer

	return &response{
		resp:    resp,
		resBody: resBodyBuffer,
	}
}

func (h *response) Body() []byte {
	/*resBody := new(bytes.Buffer)
	_ = io.MultiWriter(h.resp.Writer, resBody)

	io.ReadAll(c.Request().Body)*/
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
