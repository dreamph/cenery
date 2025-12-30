package fiber

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dreamph/cenery"
	"github.com/gofiber/fiber/v2"
)

func TestFiberParams(t *testing.T) {
	app := fiber.New()

	app.Get("/:id", func(c *fiber.Ctx) error {
		ctx := NewServerCtx(c)

		// Test with value
		id := ctx.Params("id")
		if id != "123" {
			return c.SendString("Expected 123, got " + id)
		}

		// Test with default value
		missing := ctx.Params("missing", "default")
		if missing != "default" {
			return c.SendString("Expected default, got " + missing)
		}

		return c.SendString("OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/123", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "OK" {
		t.Errorf("Test failed: %s", string(body))
	}
}

func TestFiberRouteParam(t *testing.T) {
	server := fiber.New()
	a := New(server).(*app)
	a.Get("/users/:id", func(c cenery.Ctx) error {
		if got := c.Params("id"); got != "123" {
			return c.SendString(http.StatusBadRequest, got)
		}
		return c.SendString(http.StatusOK, "OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "OK" {
		t.Errorf("Test failed: %s", string(body))
	}
}

func TestFiberQueryParam(t *testing.T) {
	app := fiber.New()

	app.Get("/search", func(c *fiber.Ctx) error {
		ctx := NewServerCtx(c)

		// Test with value
		q := ctx.QueryParam("q")
		if q != "test" {
			return c.SendString("Expected test, got " + q)
		}

		// Test with default value
		missing := ctx.QueryParam("missing", "default")
		if missing != "default" {
			return c.SendString("Expected default, got " + missing)
		}

		return c.SendString("OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/search?q=test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "OK" {
		t.Errorf("Test failed: %s", string(body))
	}
}

func TestFiberBodyParser(t *testing.T) {
	app := fiber.New()

	app.Post("/api/test", func(c *fiber.Ctx) error {
		ctx := NewServerCtx(c)

		var result struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}

		if err := ctx.BodyParser(&result); err != nil {
			return c.Status(400).SendString("Parse error: " + err.Error())
		}

		if result.Name != "test" || result.Value != 123 {
			return c.SendString("Parsed incorrectly")
		}

		return c.SendString("OK")
	})

	jsonData := `{"name":"test","value":123}`
	req := httptest.NewRequest(http.MethodPost, "/api/test", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "OK" {
		t.Errorf("Test failed: %s", string(body))
	}
}

func TestFiberBodyParserStream(t *testing.T) {
	app := fiber.New()

	app.Post("/api/test", func(c *fiber.Ctx) error {
		ctx := NewServerCtx(c)

		var result struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}

		if err := ctx.BodyParserStream(&result); err != nil {
			return c.Status(400).SendString("Parse error: " + err.Error())
		}

		if result.Name != "test" || result.Value != 123 {
			return c.SendString("Parsed incorrectly")
		}

		return c.SendString("OK")
	})

	jsonData := `{"name":"test","value":123}`
	req := httptest.NewRequest(http.MethodPost, "/api/test", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "OK" {
		t.Errorf("Test failed: %s", string(body))
	}
}

func TestFiberBodyStream(t *testing.T) {
	app := fiber.New()

	app.Post("/api/test", func(c *fiber.Ctx) error {
		ctx := NewServerCtx(c)
		body := ctx.BodyStream()
		if body == nil {
			return c.Status(400).SendString("BodyStream nil")
		}
		defer body.Close()

		data, err := io.ReadAll(body)
		if err != nil {
			return c.Status(500).SendString("Read error: " + err.Error())
		}
		if string(data) != `{"name":"test"}` {
			return c.SendString("Body mismatch")
		}
		return c.SendString("OK")
	})

	jsonData := `{"name":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/test", strings.NewReader(jsonData))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "OK" {
		t.Errorf("Test failed: %s", string(body))
	}
}

func TestFiberSendJSON(t *testing.T) {
	app := fiber.New()

	app.Get("/api/test", func(c *fiber.Ctx) error {
		ctx := NewServerCtx(c)

		data := map[string]any{
			"status":  "ok",
			"message": "success",
		}

		return ctx.SendJSON(200, data)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), `"status":"ok"`) {
		t.Errorf("Expected JSON with status:ok, got: %s", string(body))
	}
}

func TestFiberSendStream(t *testing.T) {
	app := fiber.New()

	app.Get("/stream", func(c *fiber.Ctx) error {
		ctx := NewServerCtx(c)
		data := strings.NewReader("Hello, World!")
		return ctx.SendStream(200, "text/plain", data)
	})

	req := httptest.NewRequest(http.MethodGet, "/stream", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got: %s", string(body))
	}
}

func TestFiberFormFileStream(t *testing.T) {
	app := fiber.New()

	app.Post("/upload", func(c *fiber.Ctx) error {
		ctx := NewServerCtx(c)

		stream, err := ctx.FormFileStream("file")
		if err != nil {
			return c.Status(400).SendString("Error: " + err.Error())
		}
		defer stream.File.Close()

		if stream.FileName != "test.txt" {
			return c.SendString("Wrong filename: " + stream.FileName)
		}

		content, _ := io.ReadAll(stream.File)
		if string(content) != "file content" {
			return c.SendString("Wrong content: " + string(content))
		}

		return c.SendString("OK")
	})

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	io.WriteString(part, "file content")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	if string(responseBody) != "OK" {
		t.Errorf("Test failed: %s", string(responseBody))
	}
}

func TestFiberResponseHeaders(t *testing.T) {
	app := fiber.New()

	app.Get("/headers", func(c *fiber.Ctx) error {
		ctx := NewServerCtx(c)

		ctx.Response().SetHeader("X-Custom", "value1")
		ctx.Response().AddHeader("X-Multi", "value1")
		ctx.Response().AddHeader("X-Multi", "value2")

		return c.SendString("OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/headers", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if custom := resp.Header.Get("X-Custom"); custom != "value1" {
		t.Errorf("Expected X-Custom=value1, got: %s", custom)
	}

	multiValues := resp.Header.Values("X-Multi")
	if len(multiValues) < 1 || !strings.Contains(strings.Join(multiValues, ","), "value1") {
		t.Errorf("Expected X-Multi to contain value1, got: %v", multiValues)
	}
}

// NOTE: Fiber benchmarks use app.Test() which includes routing overhead
// This is different from Echo benchmarks which test pure operations
// Fiber's routing cannot be easily separated from context operations
// For fair comparison of REAL performance, use load tests (loadtest/loadtest.go)
//
// In production: Fiber is typically 30-50% FASTER than Echo due to fasthttp

func BenchmarkFiberParams(b *testing.B) {
	app := fiber.New()

	app.Get("/:id", func(c *fiber.Ctx) error {
		ctx := NewServerCtx(c)
		_ = ctx.Params("id")
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/123", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = app.Test(req, -1) // -1 = no timeout
	}
}

func BenchmarkFiberSendJSON(b *testing.B) {
	app := fiber.New()

	data := map[string]any{
		"status":  "ok",
		"message": "success",
	}

	app.Get("/bench", func(c *fiber.Ctx) error {
		ctx := NewServerCtx(c)
		return ctx.SendJSON(200, data)
	})

	req := httptest.NewRequest(http.MethodGet, "/bench", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = app.Test(req, -1)
	}
}
