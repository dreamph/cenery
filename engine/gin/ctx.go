package gin

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"sync"
	"sync/atomic"

	"github.com/dreamph/cenery"
	"github.com/gin-gonic/gin"
)

type serverCtx struct {
	ctx  *gin.Context
	resp cenery.Response
}

func NewServerCtx(ctx *gin.Context) cenery.Ctx {
	resp := NewResponse(ctx)
	return &serverCtx{
		ctx:  ctx,
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
	val := s.ctx.Query(key)
	if len(defaultValue) == 1 {
		if val == "" {
			val = defaultValue[0]
		}
	}
	return val
}

func (s *serverCtx) BodyParser(out any) error {
	return s.ctx.ShouldBind(out)
}

func (s *serverCtx) BodyParserStream(out any) error {
	if s.ctx.Request.Body == nil {
		return errors.New("request body can't be empty")
	}
	dec := json.NewDecoder(s.ctx.Request.Body)
	if err := dec.Decode(out); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return nil
}

func (s *serverCtx) BodyStream() io.ReadCloser {
	return s.ctx.Request.Body
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
	s.ctx.Status(status)
	_, err := s.ctx.Writer.WriteString(data)
	return err
}

func (s *serverCtx) Send(status int, data []byte) error {
	s.ctx.Status(status)
	if !enableSendBufferPooling.Load() || len(data) == 0 {
		if len(data) == 0 {
			return nil
		}
		_, err := s.ctx.Writer.Write(data)
		return err
	}

	buf := sendBufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	_, _ = buf.Write(data)
	_, err := buf.WriteTo(s.ctx.Writer)
	sendBufferPool.Put(buf)
	return err
}

func (s *serverCtx) SendJSON(status int, data any) error {
	s.ctx.JSON(status, data)
	return nil
}

func (s *serverCtx) SendStream(status int, contentType string, reader io.Reader) error {
	s.ctx.Header("Content-Type", contentType)
	s.ctx.Status(status)
	_, err := io.Copy(s.ctx.Writer, reader)
	return err
}

func (s *serverCtx) Request() cenery.Request {
	return NewRequest(s.ctx.Request)
}

func (s *serverCtx) Response() cenery.Response {
	return s.resp
}

func (s *serverCtx) Next() error {
	s.ctx.Next()
	if len(s.ctx.Errors) > 0 {
		return s.ctx.Errors.Last()
	}
	return nil
}
