package gin_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	harmoniendpoint "github.com/harmonikit/harmoni/endpoint"
	gintransport "github.com/harmonikit/kit/transport/gin"
	"github.com/gin-gonic/gin"
)

func ExampleServer() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	add := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req + 1, nil
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

	server := gintransport.NewServer(add, dec, enc)
	r.POST("/add", server.Handle)

	req := httptest.NewRequest("POST", "/add", strings.NewReader("41"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	fmt.Printf("status=%d", rec.Code)
	// Output: status=200
}
