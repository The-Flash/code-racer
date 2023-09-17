package v1

import "github.com/gofiber/fiber/v2"

type Router struct{}

func (r *Router) Setup(route fiber.Router) {
	route.Get("/health", r.health)
	route.Get("/runtimes", r.runtimes)
	route.Post("/execute", r.execute)
}

func (r *Router) health(ctx *fiber.Ctx) error {
	return ctx.SendString("OK")
}

func (r *Router) runtimes(ctx *fiber.Ctx) error {
	return nil
}

func (r *Router) execute(ctx *fiber.Ctx) error {
	return nil
}
