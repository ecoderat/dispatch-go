package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ecoderat/dispatch-go/internal/service"
)

type MessageController interface {
	Start(c *fiber.Ctx) error
	Stop(c *fiber.Ctx) error
	GetMessages(c *fiber.Ctx) error
}

type messageController struct {
	service service.MessageService
}

func NewMessageController(messageService service.MessageService) MessageController {
	return &messageController{
		service: messageService,
	}
}

func (ctrl *messageController) Start(c *fiber.Ctx) error {
	err := ctrl.service.Start(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to start server")
	}

	return c.SendString("Server started")
}

func (ctrl *messageController) Stop(c *fiber.Ctx) error {
	err := ctrl.service.Stop(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to stop server")
	}

	return c.SendString("Server stopped")
}

func (ctrl *messageController) GetMessages(c *fiber.Ctx) error {
	messages, err := ctrl.service.GetMessages(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve messages")
	}

	return c.JSON(messages)
}
