package chi

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/dreamph/cenery"
	"github.com/go-chi/chi/v5"
)

type serverCtx struct {
	w    http.ResponseWriter
	r    *http.Request
	next http.Handler
	resp cenery.Response
}

func NewServerCtx(w http.ResponseWriter, r *http.Request, next http.Handler) cenery.Ctx {
	resp, writer := NewResponse(w)
	return &serverCtx{
		w:    writer,
		r:    r,
		next: next,
		resp: resp,
	}
}

func (s *serverCtx) Params(key string, defaultValue ...string) string {
	val := chi.URLParam(s.r, key)
	if len(defaultValue) == 1 {
		if val == "" {
			val = defaultValue[0]
		}
	}
	return val
}

func (s *serverCtx) QueryParam(key string, defaultValue ...string) string {
	val := s.r.URL.Query().Get(key)
	if len(defaultValue) == 1 {
		if val == "" {
			val = defaultValue[0]
		}
	}
	return val
}

func (s *serverCtx) BodyParser(out any) error {
	if s.r.Body == nil {
		return errors.New("request body can't be empty")
	}
	data, err := io.ReadAll(s.r.Body)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}

	s.r.Body = io.NopCloser(bytes.NewBuffer(data))
	return jsonUnmarshal(data, out)
}

func (s *serverCtx) FormFile(fileKey string) *cenery.FileData {
	return FormFile(s.r, fileKey)
}

func (s *serverCtx) FormFiles(fileKey string) *[]cenery.FileData {
	return FormFiles(s.r, fileKey)
}

func (s *serverCtx) FormFileStream(fileKey string) (*cenery.FileStream, error) {
	return FormFileStream(s.r, fileKey)
}

func (s *serverCtx) FormFilesStream(fileKey string) ([]*cenery.FileStream, error) {
	return FormFilesStream(s.r, fileKey)
}

var enableSendBufferPooling atomic.Bool

// EnableSendBufferPooling toggles pooling for Send() to reuse buffers
// when constructing response payloads in hot paths.
func EnableSendBufferPooling(enabled bool) {
	enableSendBufferPooling.Store(enabled)
}

var sendBufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 32*1024))
	},
}

func (s *serverCtx) SendString(status int, data string) error {
	s.w.WriteHeader(status)
	_, err := io.WriteString(s.w, data)
	return err
}

func (s *serverCtx) Send(status int, data []byte) error {
	s.w.WriteHeader(status)

	if !enableSendBufferPooling.Load() || len(data) == 0 {
		if len(data) == 0 {
			return nil
		}
		_, err := s.w.Write(data)
		return err
	}

	buf := sendBufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	_, _ = buf.Write(data)
	_, err := buf.WriteTo(s.w)
	sendBufferPool.Put(buf)
	return err
}

func (s *serverCtx) SendJSON(status int, data any) error {
	payload, err := jsonMarshal(data)
	if err != nil {
		return err
	}
	s.w.Header().Set("Content-Type", "application/json")
	s.w.WriteHeader(status)
	_, err = s.w.Write(payload)
	return err
}

func (s *serverCtx) SendStream(status int, contentType string, reader io.Reader) error {
	s.w.Header().Set("Content-Type", contentType)
	s.w.WriteHeader(status)
	_, err := io.Copy(s.w, reader)
	return err
}

func (s *serverCtx) Request() cenery.Request {
	return NewRequest(s.r)
}

func (s *serverCtx) Response() cenery.Response {
	return s.resp
}

func (s *serverCtx) Next() error {
	if s.next == nil {
		return nil
	}
	s.next.ServeHTTP(s.w, s.r)
	return nil
}
