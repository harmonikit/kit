package fiber_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"

	harmoniendpoint "github.com/harmonikit/harmoni/endpoint"
	fibertransport "github.com/harmonikit/kit/transport/fiber"
	"github.com/gofiber/fiber/v3"
)

func ExampleServer() {
	app := fiber.New()

	add := harmoniendpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req + 1, nil
	})

	dec := func(ctx context.Context, c fiber.Ctx) (int, error) {
		var v int
		json.Unmarshal(c.Body(), &v)
		return v, nil
	}

	enc := func(ctx context.Context, c fiber.Ctx, resp int) error {
		return c.JSON(resp)
	}

	server := fibertransport.NewServer(add, dec, enc)
	app.Post("/add", server.Handle)

	req := httptest.NewRequest("POST", "/add", strings.NewReader("41"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	fmt.Printf("status=%d", resp.StatusCode)
	// Output: status=200
}
