package api

import (
	"fmt"
	"os"
	"time"

	metricsRouter "github.com/The-Flash/code-racer/internal/api/metrics"
	v1 "github.com/The-Flash/code-racer/internal/api/v1"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sarulabs/di/v2"
)

type API struct {
	bindPort    string
	metricsPort string
	app         *fiber.App
	metricsApp  *fiber.App
}

func New(ctn di.Container) (api *API, err error) {
	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "8000"
	}
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "8080"
	}

	api = &API{
		bindPort:    fmt.Sprintf(":%s", appPort),
		metricsPort: fmt.Sprintf(":%s", metricsPort),
	}
	api.SetupAPI(ctn)
	api.SetupMetricsAPI()
	return
}

func (api *API) SetupAPI(ctn di.Container) {

	api.app = fiber.New(fiber.Config{
		ServerHeader: "code-racer",
	})

	api.app.Use(recover.New())

	api.app.Use(cors.New())

	// 5 requests per second
	api.app.Use(limiter.New(limiter.Config{
		Max:               5,
		Expiration:        1 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	new(v1.Router).Setup(api.app.Group("/api/v1"), ctn)
}

func (api *API) SetupMetricsAPI() {
	api.metricsApp = fiber.New(fiber.Config{
		ServerHeader: "code-racer-metrics",
	})

	new(metricsRouter.MetricsRouter).Setup(api.metricsApp.Group("/metrics"))
}

func (api *API) ListenAndServeBlocking() error {
	return api.app.Listen(api.bindPort)
}

func (api *API) ListenAndServeBlockingMetrics() error {
	return api.metricsApp.Listen(api.metricsPort)
}
