package fiberapp

import (
	"context"
	"log"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Config struct {
	Host             string `validate:"required"`
	IPHeader         string
	AllowOrigins     string `validate:"required"`
	AllowMethods     string `validate:"required"`
	AllowCredentials bool
	AllowHeaders     string `validate:"required"`
	ExposeHeaders    string `validate:"required"`
}

type FiberApp struct {
	cfg    Config
	logger *slog.Logger
	App    *fiber.App
}

func New(cfg Config, logger *slog.Logger) *FiberApp {
	return &FiberApp{
		cfg:    cfg,
		logger: logger,
		App: fiber.New(fiber.Config{
			DisableStartupMessage: true,
		}),
	}
}

func (a *FiberApp) Start(ctx context.Context) error {
	a.App.Use(cors.New(cors.Config{
		AllowOrigins:     a.cfg.AllowOrigins,
		AllowMethods:     a.cfg.AllowMethods,
		AllowCredentials: a.cfg.AllowCredentials,
		AllowHeaders:     a.cfg.AllowHeaders,
		ExposeHeaders:    a.cfg.ExposeHeaders,
	}))

	a.App.Get("/health_check", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	go func() {
		if err := a.App.Listen(a.cfg.Host); err != nil {
			log.Fatal(err.Error())
		}
	}()

	return nil
}

func (a *FiberApp) Stop(ctx context.Context) error {
	return a.App.ShutdownWithContext(ctx)
}
