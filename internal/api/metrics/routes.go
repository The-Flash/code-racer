package metrics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type MetricsRouter struct{}

func (r *MetricsRouter) Setup(route fiber.Router) {
	route.Get("/", r.getMetrics)
}

func (r *MetricsRouter) getMetrics(ctx *fiber.Ctx) error {
	rHandler := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
	rHandler(ctx.Context())
	return nil
}
