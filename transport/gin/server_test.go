package gin_test

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
	gintransport "github.com/harmonikit/kit/transport/gin"
	"github.com/gin-gonic/gin"
)

func TestServer_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req * 2, nil
	})

	dec := func(ctx context.Context, c *gin.Context) (int, error) {
		body, _ := io.ReadAll(c.Request.Body)
		var v int
		json.Unmarshal(body, &v)
		return v, nil
	}

	enc := func(ctx context.Context, c *gin.Context, resp int) error {
		c.JSON(http.StatusOK, resp)
		return nil
	}

	server := gintransport.NewServer(ep, dec, enc)
	r.POST("/", server.Handle)

	req := httptest.NewRequest("POST", "/", strings.NewReader("21"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestServer_DecodeError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})

	dec := func(ctx context.Context, c *gin.Context) (int, error) {
		return 0, errors.New("bad request")
	}

	enc := func(ctx context.Context, c *gin.Context, resp int) error { return nil }

	server := gintransport.NewServer(ep, dec, enc)
	r.POST("/", server.Handle)

	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestServer_EndpointError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	wantErr := errors.New("business logic error")
	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, wantErr
	})

	dec := func(ctx context.Context, c *gin.Context) (int, error) { return 1, nil }
	enc := func(ctx context.Context, c *gin.Context, resp int) error { return nil }

	server := gintransport.NewServer(ep, dec, enc)
	r.POST("/", server.Handle)

	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if !strings.Contains(rec.Body.String(), wantErr.Error()) {
		t.Fatalf("body %q should contain %q", rec.Body.String(), wantErr.Error())
	}
}

func TestServer_CustomErrorEncoder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, errors.New("not found")
	})

	dec := func(ctx context.Context, c *gin.Context) (int, error) { return 1, nil }
	enc := func(ctx context.Context, c *gin.Context, resp int) error { return nil }

	customErr := func(ctx context.Context, c *gin.Context, err error) error {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return nil
	}

	server := gintransport.NewServer(ep, dec, enc,
		gintransport.WithErrorEncoder[int, int](customErr))
	r.POST("/", server.Handle)

	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusNotFound)
	}
}
