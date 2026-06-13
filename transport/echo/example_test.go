package echo_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	harmoniendpoint "github.com/harmonikit/harmoni/endpoint"
	echotransport "github.com/harmonikit/kit/transport/echo"
	"github.com/labstack/echo/v4"
)

func ExampleServer() {
	e := echo.New()

	add := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req + 1, nil
	})

	dec := func(ctx context.Context, c echo.Context) (int, error) {
		body, _ := io.ReadAll(c.Request().Body)
		var v int
		json.Unmarshal(body, &v)
		return v, nil
	}

	enc := func(ctx context.Context, c echo.Context, resp int) error {
		return c.JSON(http.StatusOK, resp)
	}

	server := echotransport.NewServer(add, dec, enc)
	e.POST("/add", server.Handle)

	req := httptest.NewRequest("POST", "/add", strings.NewReader("41"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	fmt.Printf("status=%d", rec.Code)
	// Output: status=200
}
