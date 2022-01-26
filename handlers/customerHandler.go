package handlers

import (
	"fmt"
	"freq/helper"
	"freq/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

type CustomerHandler struct {
	CustomerService services.CustomerService
}

func (ch *CustomerHandler) FindAll(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	newCustomerQuery := c.Query("new", "false")

	isNew, err := strconv.ParseBool(newCustomerQuery)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("must provide a valid value")})
	}

	ips, err := ch.CustomerService.FindAll(page, isNew)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": ips})
}

func (ch *CustomerHandler) FindAllByFullName(c *fiber.Ctx) error {
	c.Accepts("application/json")
	page := c.Query("page", "1")
	newCustomerQuery := c.Query("new", "false")
	firstName := c.Query("firstName", "")
	lastName := c.Query("lastName", "")

	isNew, err := strconv.ParseBool(newCustomerQuery)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("must provide a valid value")})
	}

	customers, err := ch.CustomerService.FindAllByFullName(strings.ToLower(firstName), strings.ToLower(lastName), page, isNew)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": customers})
}

func (ch *CustomerHandler) FindAllByOptInStatus(c *fiber.Ctx) error {
	customers, err := ch.CustomerService.FindAllByOptInStatus(true)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": customers})
}

func (ch *CustomerHandler) UpdateOptInStatus(c *fiber.Ctx) error {
	email := c.Params("email")

	if !helper.IsEmail(email) {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": "invalid email"})
	}

	customer, err := ch.CustomerService.UpdateOptInStatus(false, email)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": customer.Email})
}
