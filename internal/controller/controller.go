package controller

import "github.com/gofiber/fiber/v2"

type Controller interface {
	Start(c *fiber.Ctx) error
	Stop(c *fiber.Ctx) error
	GetMessages(c *fiber.Ctx) error
}

type controller struct{}

func NewController() Controller {
	return &controller{}
}

func (r *controller) Start(c *fiber.Ctx) error {
	return c.SendString("Server started")
}

func (r *controller) Stop(c *fiber.Ctx) error {
	return c.SendString("Server stopped")
}

func (r *controller) GetMessages(c *fiber.Ctx) error {
	return c.SendString("List of sent messages")
}
