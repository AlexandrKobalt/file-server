package http

import (
	"github.com/AlexandrKobalt/trip-track_file-server/internal/file/handler"
	"github.com/gofiber/fiber/v2"
)

func Map(group fiber.Router, h handler.IHandler) {
	group.Get("/:key", h.GetFile())
}
