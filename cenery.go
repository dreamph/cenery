package cenery

import (
	"context"
	"io"
	"mime/multipart"
)

type FileData struct {
	FileData        []byte `json:"fileData"`
	FileName        string `json:"fileName"`
	FileSize        int64  `json:"fileSize"`
	FileContentType string `json:"fileContentType"`
}

type FileStream struct {
	File            multipart.File
	FileName        string
	FileSize        int64
	FileContentType string
}

type Ctx interface {
	Params(key string, defaultValue ...string) string
	QueryParam(key string, defaultValue ...string) string
	BodyParser(out any) error
	BodyParserStream(out any) error
	BodyStream() io.ReadCloser
	FormFile(fileKey string) *FileData
	FormFiles(fileKey string) *[]FileData

	FormFileStream(fileKey string) (*FileStream, error)
	FormFilesStream(fileKey string) ([]*FileStream, error)

	SendString(status int, data string) error
	Send(status int, data []byte) error
	SendJSON(status int, data any) error

	// Streaming response (no memory allocation)
	SendStream(status int, contentType string, reader io.Reader) error

	Request() Request
	Response() Response
	Next() error
}

type Request interface {
	Body() []byte
	SetBody(data []byte)

	BodyStream() io.ReadCloser

	GetHeader(key string) string
	SetHeader(key string, val string)
	AddHeader(key string, val string)
}

type Response interface {
	Body() []byte
	SetBody(data []byte)

	GetHeader(key string) string
	SetHeader(key string, val string)
	AddHeader(key string, val string)
}

type Handler = func(Ctx) error

type App interface {
	Get(path string, handlers ...Handler)
	Post(path string, handlers ...Handler)
	Put(path string, handlers ...Handler)
	Delete(path string, handlers ...Handler)
	Head(path string, handlers ...Handler)
	Options(path string, handlers ...Handler)
	Connect(path string, handlers ...Handler)
	Patch(path string, handlers ...Handler)
	Trace(path string, handlers ...Handler)
	Name() string
	Use(handlers ...Handler)
	Listen(addr string) error
	Shutdown(ctx context.Context) error
}
