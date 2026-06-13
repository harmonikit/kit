package gin_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	harmoniendpoint "github.com/harmonikit/harmoni/endpoint"
	gintransport "github.com/harmonikit/kit/transport/gin"
	"github.com/gin-gonic/gin"
)

func TestServer_EncodeError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	ep := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})

	dec := func(ctx context.Context, c *gin.Context) (int, error) { return 1, nil }
	enc := func(ctx context.Context, c *gin.Context, resp int) error {
		return errors.New("encode fail")
	}

	server := gintransport.NewServer(ep, dec, enc)
	r.POST("/", server.Handle)

	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
