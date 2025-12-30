package fasthttp

import (
	"bytes"
	"io"
	"mime/multipart"
	"net"
	"strings"
	"testing"

	"github.com/dreamph/cenery"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func newCtx(method, uri string, body []byte, contentType string) *fasthttp.RequestCtx {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	req.SetRequestURI(uri)
	if len(body) > 0 {
		req.SetBody(body)
	}
	if contentType != "" {
		req.Header.SetContentType(contentType)
	}

	var ctx fasthttp.RequestCtx
	ctx.Init(req, &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)}, nil)
	return &ctx
}

func TestParams(t *testing.T) {
	ctx := newCtx(fasthttp.MethodGet, "/users/123", nil, "")
	ctx.SetUserValue("id", "123")

	svc := NewServerCtx(ctx, nil)

	if got := svc.Params("id"); got != "123" {
		t.Errorf("Params() = %v, want %v", got, "123")
	}

	if got := svc.Params("id", "default"); got != "123" {
		t.Errorf("Params() with default = %v, want %v", got, "123")
	}

	if got := svc.Params("missing", "default"); got != "default" {
		t.Errorf("Params() with missing param = %v, want %v", got, "default")
	}
}

func TestQueryParam(t *testing.T) {
	ctx := newCtx(fasthttp.MethodGet, "/search?q=test&page=1", nil, "")

	svc := NewServerCtx(ctx, nil)

	if got := svc.QueryParam("q"); got != "test" {
		t.Errorf("QueryParam() = %v, want %v", got, "test")
	}

	if got := svc.QueryParam("q", "default"); got != "test" {
		t.Errorf("QueryParam() with default = %v, want %v", got, "test")
	}

	if got := svc.QueryParam("missing", "default"); got != "default" {
		t.Errorf("QueryParam() with missing param = %v, want %v", got, "default")
	}
}

func TestRouteParam(t *testing.T) {
	routerApp := router.New()
	a := New(routerApp).(*app)
	a.Get("/users/:id", func(c cenery.Ctx) error {
		if got := c.Params("id"); got != "123" {
			return c.SendString(fasthttp.StatusBadRequest, got)
		}
		return c.SendString(fasthttp.StatusOK, "ok")
	})

	ctx := newCtx(fasthttp.MethodGet, "/users/123", nil, "")
	handler := a.Handler()
	handler(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatalf("status = %v, want %v", ctx.Response.StatusCode(), fasthttp.StatusOK)
	}
	if string(ctx.Response.Body()) != "ok" {
		t.Fatalf("body = %v, want %v", string(ctx.Response.Body()), "ok")
	}
}

func TestBodyParser(t *testing.T) {
	jsonData := `{"name":"test","value":123}`
	ctx := newCtx(fasthttp.MethodPost, "/", []byte(jsonData), "application/json")

	svc := NewServerCtx(ctx, nil)

	var result struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	if err := svc.BodyParser(&result); err != nil {
		t.Errorf("BodyParser() error = %v", err)
	}

	if result.Name != "test" || result.Value != 123 {
		t.Errorf("BodyParser() parsed incorrectly: got %+v", result)
	}
}

func TestBodyParserStream(t *testing.T) {
	jsonData := `{"name":"test","value":123}`
	ctx := newCtx(fasthttp.MethodPost, "/", []byte(jsonData), "application/json")

	svc := NewServerCtx(ctx, nil)

	var result struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	if err := svc.BodyParserStream(&result); err != nil {
		t.Errorf("BodyParserStream() error = %v", err)
	}

	if result.Name != "test" || result.Value != 123 {
		t.Errorf("BodyParserStream() parsed incorrectly: got %+v", result)
	}
}

func TestBodyStream(t *testing.T) {
	jsonData := `{"name":"test"}`
	ctx := newCtx(fasthttp.MethodPost, "/", []byte(jsonData), "")

	svc := NewServerCtx(ctx, nil)

	body := svc.BodyStream()
	if body == nil {
		t.Fatalf("BodyStream() returned nil")
	}
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("BodyStream() read error: %v", err)
	}
	if string(data) != jsonData {
		t.Errorf("BodyStream() = %v, want %v", string(data), jsonData)
	}
}

func TestSendJSON(t *testing.T) {
	ctx := newCtx(fasthttp.MethodGet, "/", nil, "")

	svc := NewServerCtx(ctx, nil)

	data := map[string]any{
		"status":  "ok",
		"message": "success",
	}

	if err := svc.SendJSON(200, data); err != nil {
		t.Errorf("SendJSON() error = %v", err)
	}

	if ctx.Response.StatusCode() != 200 {
		t.Errorf("SendJSON() status = %v, want %v", ctx.Response.StatusCode(), 200)
	}

	if !strings.Contains(string(ctx.Response.Body()), `"status":"ok"`) {
		t.Errorf("SendJSON() body = %v, want to contain status:ok", string(ctx.Response.Body()))
	}
}

func TestSendStream(t *testing.T) {
	ctx := newCtx(fasthttp.MethodGet, "/", nil, "")

	svc := NewServerCtx(ctx, nil)

	data := strings.NewReader("Hello, World!")

	if err := svc.SendStream(200, "text/plain", data); err != nil {
		t.Errorf("SendStream() error = %v", err)
	}

	if ctx.Response.StatusCode() != 200 {
		t.Errorf("SendStream() status = %v, want %v", ctx.Response.StatusCode(), 200)
	}

	if string(ctx.Response.Header.ContentType()) != "text/plain" {
		t.Errorf("SendStream() Content-Type = %v, want %v", string(ctx.Response.Header.ContentType()), "text/plain")
	}

	stream := ctx.Response.BodyStream()
	if stream == nil {
		t.Fatalf("SendStream() expected body stream")
	}

	content, _ := io.ReadAll(stream)
	if string(content) != "Hello, World!" {
		t.Errorf("SendStream() body = %v, want %v", string(content), "Hello, World!")
	}
}

func TestFormFileStream(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	io.WriteString(part, "file content")
	writer.Close()

	ctx := newCtx(fasthttp.MethodPost, "/upload", body.Bytes(), writer.FormDataContentType())

	svc := NewServerCtx(ctx, nil)

	stream, err := svc.FormFileStream("file")
	if err != nil {
		t.Fatalf("FormFileStream() error = %v", err)
	}
	defer stream.File.Close()

	if stream.FileName != "test.txt" {
		t.Errorf("FormFileStream() FileName = %v, want %v", stream.FileName, "test.txt")
	}

	content, _ := io.ReadAll(stream.File)
	if string(content) != "file content" {
		t.Errorf("FormFileStream() content = %v, want %v", string(content), "file content")
	}
}

func TestResponseHeaders(t *testing.T) {
	ctx := newCtx(fasthttp.MethodGet, "/", nil, "")

	svc := NewServerCtx(ctx, nil)

	svc.Response().SetHeader("X-Custom", "value1")
	if got := svc.Response().GetHeader("X-Custom"); got != "value1" {
		t.Errorf("SetHeader/GetHeader = %v, want %v", got, "value1")
	}

	svc.Response().AddHeader("X-Multi", "value1")
	svc.Response().AddHeader("X-Multi", "value2")
	if got := svc.Response().GetHeader("X-Multi"); !strings.Contains(got, "value1") {
		t.Errorf("AddHeader = %v, want to contain value1", got)
	}
}
