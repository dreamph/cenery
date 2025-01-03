package fiber

import (
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

func (s *serverCtx) BodyParser(out interface{}) error {
	return s.ctx.BodyParser(out)
}

func (s *serverCtx) FormFile(fileKey string) *cenery.FileData {
	return FormFile(s.ctx, fileKey)
}

func (s *serverCtx) FormFiles(fileKey string) *[]cenery.FileData {
	return FormFiles(s.ctx, fileKey)
}

func (s *serverCtx) SendString(status int, data string) error {
	return s.ctx.Status(status).SendString(data)
}

func (s *serverCtx) Send(status int, data []byte) error {
	return s.ctx.Status(status).Send(data)
}

func (s *serverCtx) SendJSON(status int, data interface{}) error {
	return s.ctx.Status(status).JSON(data)
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
