package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type handler struct {
	fileDirectory string
}

func New(fileDirectory string) IHandler {
	return &handler{fileDirectory: fileDirectory}
}

func (h *handler) GetFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fileKey := c.Params("key")
		filePath := fmt.Sprintf("%s/%s", h.fileDirectory, fileKey)
		return c.SendFile(filePath)
	}
}
