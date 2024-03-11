package echo

import (
	"github.com/dreamph/cenery"
	"github.com/labstack/echo/v4"
)

type serverCtx struct {
	ctx echo.Context
}

func NewServerCtx(ctx echo.Context) cenery.Ctx {
	return &serverCtx{ctx: ctx}
}

func (s *serverCtx) Params(key string, defaultValue ...string) string {
	val := s.ctx.Param(key)
	if len(defaultValue) == 1 {
		if val != "" {
			val = defaultValue[0]
		}
	}
	return val
}

func (s *serverCtx) QueryParam(key string, defaultValue ...string) string {
	val := s.ctx.QueryParam(key)
	if len(defaultValue) == 1 {
		if val != "" {
			val = defaultValue[0]
		}
	}
	return val
}

func (s *serverCtx) BodyParser(out interface{}) error {
	return s.ctx.Bind(out)
}

func (s *serverCtx) FormFile(fileKey string) *cenery.FileData {
	return FormFile(s.ctx, fileKey)
}

func (s *serverCtx) FormFiles(fileKey string) *[]cenery.FileData {
	return FormFiles(s.ctx, fileKey)
}

func (s *serverCtx) SendString(status int, data string) error {
	return s.ctx.String(status, data)
}

func (s *serverCtx) Send(status int, data []byte) error {
	return s.SendString(status, string(data))
}

func (s *serverCtx) SendJSON(status int, data interface{}) error {
	return s.ctx.JSON(status, data)
}

func (s *serverCtx) Request() cenery.Request {
	return NewRequest(s.ctx.Request())
}

func (s *serverCtx) Response() cenery.Response {
	return NewResponse(s.ctx.Response())
}
