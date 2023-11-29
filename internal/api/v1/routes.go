package v1

import (
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/The-Flash/code-racer/internal/doc"
	"github.com/The-Flash/code-racer/internal/execution"
	"github.com/The-Flash/code-racer/internal/file_system"
	"github.com/The-Flash/code-racer/internal/manifest"
	"github.com/The-Flash/code-racer/internal/names"
	"github.com/The-Flash/code-racer/pkg/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/sarulabs/di/v2"
)

type Router struct {
	mfest    *manifest.Manifest
	fp       *file_system.FileProvider
	executor *execution.Executor
}

func (r *Router) Setup(route fiber.Router, ctn di.Container) {
	m := ctn.Get(names.DiManifestProvider).(*manifest.Manifest)
	fp := ctn.Get(names.DiFileProvider).(*file_system.FileProvider)
	executor := ctn.Get(names.DiExecutorProvider).(*execution.Executor)
	r.mfest = m
	r.fp = fp
	r.executor = executor
	// api endpoints
	route.Get("/health", r.health)
	route.Get("/runtimes", r.runtimes)
	route.Post("/execute", r.execute)
	route.Get("/swagger/*", swagger.HandlerDefault)
}

// @summary Health check
// @description Health check endpoint
// @produce plain
// @success 200 {string} OK
// @router /health [get]
func (r *Router) health(ctx *fiber.Ctx) error {
	return ctx.SendString("OK")
}

// @summary Runtimes
// @description Get available runtimes
// @accept json
// @produce json
// @router /runtimes [get]
// @success 200 {array} manifest.ManifestRuntime
func (r *Router) runtimes(ctx *fiber.Ctx) error {
	return ctx.JSON(r.mfest.Runtimes)

}

// @summary Execute
// @description Execute code
// @accept json
// @produce json
// @param body body models.ExecutionRequest true "Execution request"
// @router /execute [post]
// @success 200 {object} models.ExecutionResponse
// @failure 400 {string} Bad Request
// @failure 500 {string} Internal Server Error
func (r *Router) execute(ctx *fiber.Ctx) error {
	validate := validator.New()
	body := new(models.ExecutionRequest)
	if err := ctx.BodyParser(body); err != nil {
		return err
	}
	if err := validate.Struct(body); err != nil {
		return err
	}
	runtime, ok := r.mfest.GetRuntimeForLanguage(body.Language)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "runtime not found")
	}
	if ok := r.executor.IsExecutorAvailable(runtime); !ok {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("no executors available for %v", runtime.Language))
	}

	executionStartTime := time.Now()
	resp, err := r.executor.Execute(body.Files, &execution.ExecutionConfig{
		// ExecutionId is attached in Execute
		EntryPoint: body.EntryPoint,
		Runtime:    runtime,
	})
	if err != nil {
		log.Println(err)
		return errors.New("execution failed")
	}
	elapsedTime := time.Since(executionStartTime)
	resp.ExecutionTime = elapsedTime.String()
	return ctx.JSON(resp)
}
