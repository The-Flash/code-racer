package api

import (
	"fmt"
	"os"

	v1 "github.com/The-Flash/code-racer/internal/api/v1"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sarulabs/di/v2"
)

type API struct {
	bindPort string
	app      *fiber.App
}

func NewAPI(ctn di.Container) (r *API, err error) {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	r = &API{
		bindPort: fmt.Sprintf(":%s", port),
	}

	r.app = fiber.New(fiber.Config{
		ServerHeader: "code-racer",
	})

	r.app.Use(recover.New())

	// TODO: add rate limtter

	new(v1.Router).Setup(r.app.Group("/api/v1"), ctn)
	return
}

func (r *API) ListenAndServeBlocking() error {
	return r.app.Listen(r.bindPort)
}
