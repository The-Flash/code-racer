package v1

import (
	"github.com/The-Flash/code-racer/internal/manifest"
	"github.com/The-Flash/code-racer/internal/names"
	"github.com/The-Flash/code-racer/pkg/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

type Router struct {
	mfest *manifest.Manifest
}

func (r *Router) Setup(route fiber.Router, ctn di.Container) {
	m := ctn.Get(names.DiManifestProvider).(*manifest.Manifest)
	r.mfest = m
	route.Get("/health", r.health)
	route.Get("/runtimes", r.runtimes)
	route.Post("/execute", r.execute)
}

func (r *Router) health(ctx *fiber.Ctx) error {
	return ctx.SendString("OK")
}

func (r *Router) runtimes(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"runtimes": r.mfest.Runtimes,
	})
}

func (r *Router) execute(ctx *fiber.Ctx) error {
	validate := validator.New()
	body := new(models.ExecutionRequest)
	if err := ctx.BodyParser(body); err != nil {
		return err
	}

	if err := validate.Struct(body); err != nil {
		return err
	}
	return ctx.JSON(body)
}
