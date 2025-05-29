package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ecoderat/dispatch-go/internal/service/message"
	"github.com/ecoderat/dispatch-go/internal/service/scheduler"
)

type MessageController interface {
	Start(c *fiber.Ctx) error
	Stop(c *fiber.Ctx) error
	GetMessages(c *fiber.Ctx) error
}

type messageController struct {
	services services
}

type services struct {
	scheduler scheduler.Scheduler
	message   message.Service
}

func NewMessageController(msgService message.Service, schedService scheduler.Scheduler) MessageController {
	return &messageController{
		services: services{
			scheduler: schedService,
			message:   msgService,
		},
	}
}

func (ctrl *messageController) Start(c *fiber.Ctx) error {
	err := ctrl.services.scheduler.Start(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to start server")
	}

	return c.SendString("Server started")
}

func (ctrl *messageController) Stop(c *fiber.Ctx) error {
	err := ctrl.services.scheduler.Stop(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to stop server")
	}

	return c.SendString("Server stopped")
}

func (ctrl *messageController) GetMessages(c *fiber.Ctx) error {
	messages, err := ctrl.services.message.GetSentMessages(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve messages")
	}

	return c.JSON(messages)
}
