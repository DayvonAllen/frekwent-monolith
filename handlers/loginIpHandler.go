package handlers

import (
	"fmt"
	"freq/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type LoginIpHandler struct {
	LoginIpService services.LoginIpService
}

func (lh *LoginIpHandler) FindAll(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	newLoginQuery := c.Query("new", "false")

	isNew, err := strconv.ParseBool(newLoginQuery)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("must provide a valid value")})
	}

	ips, err := lh.LoginIpService.FindAll(page, isNew)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": ips})
}

func (lh *LoginIpHandler) FindByIp(c *fiber.Ctx) error {
	ip := c.Params("ip")

	foundIP, err := lh.LoginIpService.FindByIp(ip)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": foundIP})
}
