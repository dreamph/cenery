package fasthttp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/dreamph/cenery"
	"github.com/valyala/fasthttp"
)

type serverCtx struct {
	ctx  *fasthttp.RequestCtx
	next fasthttp.RequestHandler
	resp cenery.Response
}

func NewServerCtx(ctx *fasthttp.RequestCtx, next fasthttp.RequestHandler) cenery.Ctx {
	resp := NewResponse(ctx)
	return &serverCtx{
		ctx:  ctx,
		next: next,
		resp: resp,
	}
}

func (s *serverCtx) Params(key string, defaultValue ...string) string {
	val := s.ctx.UserValue(key)
	str := ""
	if val != nil {
		switch v := val.(type) {
		case string:
			str = v
		case []byte:
			str = string(v)
		default:
			str = fmt.Sprint(val)
		}
	}
	if len(defaultValue) == 1 {
		if str == "" {
			str = defaultValue[0]
		}
	}
	return str
}

func (s *serverCtx) QueryParam(key string, defaultValue ...string) string {
	val := string(s.ctx.QueryArgs().Peek(key))
	if len(defaultValue) == 1 {
		if val == "" {
			val = defaultValue[0]
		}
	}
	return val
}

func (s *serverCtx) BodyParser(out any) error {
	ct := string(s.ctx.Request.Header.ContentType())
	if !strings.HasPrefix(ct, "application/json") && ct != "" {
		return errors.New("unsupported content type")
	}

	data := s.ctx.PostBody()
	if len(data) == 0 {
		return nil
	}
	return jsonUnmarshal(data, out)
}

func (s *serverCtx) FormFile(fileKey string) *cenery.FileData {
	return FormFile(s.ctx, fileKey)
}

func (s *serverCtx) FormFiles(fileKey string) *[]cenery.FileData {
	return FormFiles(s.ctx, fileKey)
}

func (s *serverCtx) FormFileStream(fileKey string) (*cenery.FileStream, error) {
	return FormFileStream(s.ctx, fileKey)
}

func (s *serverCtx) FormFilesStream(fileKey string) ([]*cenery.FileStream, error) {
	return FormFilesStream(s.ctx, fileKey)
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
	s.ctx.SetStatusCode(status)
	s.ctx.SetBodyString(data)
	return nil
}

func (s *serverCtx) Send(status int, data []byte) error {
	s.ctx.SetStatusCode(status)

	if !enableSendBufferPooling.Load() || len(data) == 0 {
		if len(data) == 0 {
			return nil
		}
		s.ctx.SetBody(data)
		return nil
	}

	buf := sendBufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	_, _ = buf.Write(data)
	s.ctx.SetBody(buf.Bytes())
	sendBufferPool.Put(buf)
	return nil
}

func (s *serverCtx) SendJSON(status int, data any) error {
	payload, err := jsonMarshal(data)
	if err != nil {
		return err
	}
	s.ctx.Response.Header.SetContentType("application/json")
	s.ctx.SetStatusCode(status)
	s.ctx.SetBody(payload)
	return nil
}

func (s *serverCtx) SendStream(status int, contentType string, reader io.Reader) error {
	s.ctx.Response.Header.Set("Content-Type", contentType)
	s.ctx.SetStatusCode(status)
	s.ctx.SetBodyStream(reader, -1)
	return nil
}

func (s *serverCtx) Request() cenery.Request {
	return NewRequest(&s.ctx.Request)
}

func (s *serverCtx) Response() cenery.Response {
	return s.resp
}

func (s *serverCtx) Next() error {
	if s.next != nil {
		s.next(s.ctx)
	}
	return nil
}
