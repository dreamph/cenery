package gin

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/users/123", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123"}}

	ctx := NewServerCtx(c)

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
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/search?q=test&page=1", nil)

	ctx := NewServerCtx(c)

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

func TestBodyParser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	jsonData := `{"name":"test","value":123}`
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	ctx := NewServerCtx(c)

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

func TestSendJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	ctx := NewServerCtx(c)

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
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	ctx := NewServerCtx(c)

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
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	io.WriteString(part, "file content")
	writer.Close()

	c.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	ctx := NewServerCtx(c)

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
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	ctx := NewServerCtx(c)

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
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/users/123", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123"}}

	ctx := NewServerCtx(c)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx.Params("id")
	}
}

func BenchmarkSendJSON(b *testing.B) {
	gin.SetMode(gin.TestMode)
	data := map[string]any{
		"status":  "ok",
		"message": "success",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := NewServerCtx(c)
		_ = ctx.SendJSON(200, data)
	}
}
