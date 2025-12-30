package gin

import (
	"bytes"
	"io"
	"net/http"

	"github.com/dreamph/cenery"
)

type request struct {
	req *http.Request
}

func NewRequest(req *http.Request) cenery.Request {
	return &request{req: req}
}

func (h *request) Body() []byte {
	if h.req.Body != nil {
		data, err := io.ReadAll(h.req.Body)
		if err != nil {
			// Reset body even on error to prevent reading again
			h.SetBody(nil)
			return nil
		}
		h.SetBody(data) // Reset for re-reading
		return data
	}
	return nil
}

func (h *request) SetBody(data []byte) {
	h.req.Body = io.NopCloser(bytes.NewReader(data))
}

func (h *request) BodyStream() io.ReadCloser {
	return h.req.Body
}

func (h *request) GetHeader(key string) string {
	return h.req.Header.Get(key)
}

func (h *request) SetHeader(key string, val string) {
	h.req.Header.Set(key, val)
}

func (h *request) AddHeader(key string, val string) {
	h.req.Header.Add(key, val)
}
