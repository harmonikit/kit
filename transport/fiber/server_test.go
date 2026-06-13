package fiber_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	harmoniendpoint "github.com/harmonikit/harmoni/endpoint"
	fibertransport "github.com/harmonikit/kit/transport/fiber"
	"github.com/gofiber/fiber/v3"
)

func testApp() *fiber.App {
	return fiber.New()
}

func TestServer_Success(t *testing.T) {
	app := testApp()

	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req * 2, nil
	})

	dec := func(ctx context.Context, c fiber.Ctx) (int, error) {
		var v int
		json.Unmarshal(c.Body(), &v)
		return v, nil
	}

	enc := func(ctx context.Context, c fiber.Ctx, resp int) error {
		return c.JSON(resp)
	}

	server := fibertransport.NewServer(ep, dec, enc)
	app.Post("/", server.Handle)

	req := httptest.NewRequest("POST", "/", strings.NewReader("21"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("got status %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestServer_DecodeError(t *testing.T) {
	app := testApp()

	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})

	dec := func(ctx context.Context, c fiber.Ctx) (int, error) {
		return 0, errors.New("bad request")
	}

	enc := func(ctx context.Context, c fiber.Ctx, resp int) error { return nil }

	server := fibertransport.NewServer(ep, dec, enc)
	app.Post("/", server.Handle)

	req := httptest.NewRequest("POST", "/", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("got status %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}
}

func TestServer_EndpointError(t *testing.T) {
	app := testApp()

	wantErr := errors.New("business logic error")
	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, wantErr
	})

	dec := func(ctx context.Context, c fiber.Ctx) (int, error) { return 1, nil }
	enc := func(ctx context.Context, c fiber.Ctx, resp int) error { return nil }

	server := fibertransport.NewServer(ep, dec, enc)
	app.Post("/", server.Handle)

	req := httptest.NewRequest("POST", "/", nil)
	resp, _ := app.Test(req)
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("got status %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}
	if !strings.Contains(string(body), wantErr.Error()) {
		t.Fatalf("body %q should contain %q", string(body), wantErr.Error())
	}
}

func TestServer_CustomErrorEncoder(t *testing.T) {
	app := testApp()

	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, errors.New("not found")
	})

	dec := func(ctx context.Context, c fiber.Ctx) (int, error) { return 1, nil }
	enc := func(ctx context.Context, c fiber.Ctx, resp int) error { return nil }

	customErr := func(ctx context.Context, c fiber.Ctx, err error) error {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	server := fibertransport.NewServer(ep, dec, enc,
		fibertransport.WithErrorEncoder[int, int](customErr))
	app.Post("/", server.Handle)

	req := httptest.NewRequest("POST", "/", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("got status %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}
