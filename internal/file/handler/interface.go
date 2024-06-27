package handler

import "github.com/gofiber/fiber/v2"

type IHandler interface {
	GetFile() fiber.Handler
}
