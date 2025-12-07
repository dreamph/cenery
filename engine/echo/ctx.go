package echo

import (
	"bytes"
	"io"
	"sync"
	"sync/atomic"

	"github.com/dreamph/cenery"
	"github.com/labstack/echo/v4"
)

type serverCtx struct {
	ctx  echo.Context
	next echo.HandlerFunc
	resp cenery.Response
}

func NewServerCtx(ctx echo.Context, next echo.HandlerFunc) cenery.Ctx {
	resp := NewResponse(ctx.Response())
	return &serverCtx{
		ctx:  ctx,
		next: next,
		resp: resp,
	}
}

func (s *serverCtx) Params(key string, defaultValue ...string) string {
	val := s.ctx.Param(key)
	if len(defaultValue) == 1 {
		if val == "" {
			val = defaultValue[0]
		}
	}
	return val
}

func (s *serverCtx) QueryParam(key string, defaultValue ...string) string {
	val := s.ctx.QueryParam(key)
	if len(defaultValue) == 1 {
		if val == "" {
			val = defaultValue[0]
		}
	}
	return val
}

func (s *serverCtx) BodyParser(out any) error {
	return s.ctx.Bind(out)
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
	return s.ctx.String(status, data)
}

func (s *serverCtx) Send(status int, data []byte) error {
	res := s.ctx.Response()
	res.WriteHeader(status)

	if !enableSendBufferPooling.Load() || len(data) == 0 {
		if len(data) == 0 {
			return nil
		}
		_, err := res.Writer.Write(data)
		return err
	}

	buf := sendBufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	_, _ = buf.Write(data)
	_, err := buf.WriteTo(res.Writer)
	sendBufferPool.Put(buf)
	return err
}

func (s *serverCtx) SendJSON(status int, data any) error {
	return s.ctx.JSON(status, data)
}

func (s *serverCtx) SendStream(status int, contentType string, reader io.Reader) error {
	s.ctx.Response().Header().Set("Content-Type", contentType)
	s.ctx.Response().WriteHeader(status)
	_, err := io.Copy(s.ctx.Response().Writer, reader)
	return err
}

func (s *serverCtx) Request() cenery.Request {
	return NewRequest(s.ctx.Request())
}

func (s *serverCtx) Response() cenery.Response {
	return s.resp
}

func (s *serverCtx) Next() error {
	if s.next != nil {
		return s.next(s.ctx)
	}
	return nil
}
