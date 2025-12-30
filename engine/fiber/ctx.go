package fiber

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/dreamph/cenery"
	"github.com/gofiber/fiber/v2"
)

type serverCtx struct {
	ctx *fiber.Ctx
}

func NewServerCtx(ctx *fiber.Ctx) cenery.Ctx {
	return &serverCtx{ctx: ctx}
}

func (s *serverCtx) Params(key string, defaultValue ...string) string {
	return s.ctx.Params(key, defaultValue...)
}

func (s *serverCtx) QueryParam(key string, defaultValue ...string) string {
	return s.ctx.Query(key, defaultValue...)
}

func (s *serverCtx) BodyParser(out any) error {
	return s.ctx.BodyParser(out)
}

func (s *serverCtx) BodyParserStream(out any) error {
	body := s.BodyStream()
	if body == nil {
		return errors.New("request body can't be empty")
	}
	dec := json.NewDecoder(body)
	if err := dec.Decode(out); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return nil
}

func (s *serverCtx) BodyStream() io.ReadCloser {
	return NewRequest(s.ctx.Request()).BodyStream()
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

func (s *serverCtx) SendString(status int, data string) error {
	return s.ctx.Status(status).SendString(data)
}

func (s *serverCtx) Send(status int, data []byte) error {
	return s.ctx.Status(status).Send(data)
}

func (s *serverCtx) SendJSON(status int, data any) error {
	return s.ctx.Status(status).JSON(data)
}

func (s *serverCtx) SendStream(status int, contentType string, reader io.Reader) error {
	s.ctx.Set("Content-Type", contentType)
	s.ctx.Status(status)
	return s.ctx.SendStream(reader)
}

func (s *serverCtx) Request() cenery.Request {
	return NewRequest(s.ctx.Request())
}

func (s *serverCtx) Response() cenery.Response {
	return NewResponse(s.ctx.Response())
}

func (s *serverCtx) Next() error {
	return s.ctx.Next()
}
