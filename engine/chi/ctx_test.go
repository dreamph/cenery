package chi

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dreamph/cenery"
	"github.com/go-chi/chi/v5"
)

func TestParams(t *testing.T) {
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "123")
	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
	rec := httptest.NewRecorder()

	ctx := NewServerCtx(rec, req, nil)

	if got := ctx.Params("id"); got != "123" {
		t.Errorf("Params() = %v, want %v", got, "123")
	}

	if got := ctx.Params("id", "default"); got != "123" {
		t.Errorf("Params() with default = %v, want %v", got, "123")
	}

	if got := ctx.Params("missing", "default"); got != "default" {
		t.Errorf("Params() with missing param = %v, want %v", got, "default")
	}
}

func TestQueryParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/search?q=test&page=1", nil)
	rec := httptest.NewRecorder()

	ctx := NewServerCtx(rec, req, nil)

	if got := ctx.QueryParam("q"); got != "test" {
		t.Errorf("QueryParam() = %v, want %v", got, "test")
	}

	if got := ctx.QueryParam("q", "default"); got != "test" {
		t.Errorf("QueryParam() with default = %v, want %v", got, "test")
	}

	if got := ctx.QueryParam("missing", "default"); got != "default" {
		t.Errorf("QueryParam() with missing param = %v, want %v", got, "default")
	}
}

func TestRouteParam(t *testing.T) {
	server := chi.NewRouter()
	a := New(server).(*app)
	a.Get("/users/:id", func(c cenery.Ctx) error {
		if got := c.Params("id"); got != "123" {
			return c.SendString(http.StatusBadRequest, got)
		}
		return c.SendString(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %v, want %v", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("body = %v, want %v", rec.Body.String(), "ok")
	}
}

func TestBodyParser(t *testing.T) {
	jsonData := `{"name":"test","value":123}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ctx := NewServerCtx(rec, req, nil)

	var result struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	if err := ctx.BodyParser(&result); err != nil {
		t.Errorf("BodyParser() error = %v", err)
	}

	if result.Name != "test" || result.Value != 123 {
		t.Errorf("BodyParser() parsed incorrectly: got %+v", result)
	}
}

func TestBodyParserStream(t *testing.T) {
	jsonData := `{"name":"test","value":123}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ctx := NewServerCtx(rec, req, nil)

	var result struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	if err := ctx.BodyParserStream(&result); err != nil {
		t.Errorf("BodyParserStream() error = %v", err)
	}

	if result.Name != "test" || result.Value != 123 {
		t.Errorf("BodyParserStream() parsed incorrectly: got %+v", result)
	}
}

func TestBodyStream(t *testing.T) {
	jsonData := `{"name":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonData))
	rec := httptest.NewRecorder()

	ctx := NewServerCtx(rec, req, nil)

	body := ctx.BodyStream()
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
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	ctx := NewServerCtx(rec, req, nil)

	data := map[string]any{
		"status":  "ok",
		"message": "success",
	}

	if err := ctx.SendJSON(200, data); err != nil {
		t.Errorf("SendJSON() error = %v", err)
	}

	if rec.Code != 200 {
		t.Errorf("SendJSON() status = %v, want %v", rec.Code, 200)
	}

	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Errorf("SendJSON() body = %v, want to contain status:ok", rec.Body.String())
	}
}

func TestSendStream(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	ctx := NewServerCtx(rec, req, nil)

	data := strings.NewReader("Hello, World!")

	if err := ctx.SendStream(200, "text/plain", data); err != nil {
		t.Errorf("SendStream() error = %v", err)
	}

	if rec.Code != 200 {
		t.Errorf("SendStream() status = %v, want %v", rec.Code, 200)
	}

	if rec.Header().Get("Content-Type") != "text/plain" {
		t.Errorf("SendStream() Content-Type = %v, want %v", rec.Header().Get("Content-Type"), "text/plain")
	}

	if rec.Body.String() != "Hello, World!" {
		t.Errorf("SendStream() body = %v, want %v", rec.Body.String(), "Hello, World!")
	}
}

func TestFormFileStream(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	io.WriteString(part, "file content")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	ctx := NewServerCtx(rec, req, nil)

	stream, err := ctx.FormFileStream("file")
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
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	ctx := NewServerCtx(rec, req, nil)

	ctx.Response().SetHeader("X-Custom", "value1")
	if got := ctx.Response().GetHeader("X-Custom"); got != "value1" {
		t.Errorf("SetHeader/GetHeader = %v, want %v", got, "value1")
	}

	ctx.Response().AddHeader("X-Multi", "value1")
	ctx.Response().AddHeader("X-Multi", "value2")
	if got := ctx.Response().GetHeader("X-Multi"); !strings.Contains(got, "value1") {
		t.Errorf("AddHeader = %v, want to contain value1", got)
	}
}

func BenchmarkParams(b *testing.B) {
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "123")
	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
	rec := httptest.NewRecorder()

	ctx := NewServerCtx(rec, req, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx.Params("id")
	}
}

func BenchmarkSendJSON(b *testing.B) {
	data := map[string]any{
		"status":  "ok",
		"message": "success",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := NewServerCtx(rec, req, nil)
		_ = ctx.SendJSON(200, data)
	}
}
