package metrics

import "github.com/gofiber/fiber/v2"

type MetricsRouter struct{}

func (r *MetricsRouter) Setup(route fiber.Router) {
	route.Get("/", r.getMetrics)
}

func (r *MetricsRouter) getMetrics(ctx *fiber.Ctx) error {
	return ctx.SendString("OK")
}
