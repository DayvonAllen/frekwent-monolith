package handlers

import (
	"fmt"
	"freq/models"
	"freq/services"
	bson2 "github.com/globalsign/mgo/bson"
	"github.com/gofiber/fiber/v2"
)

type MailMemberHandler struct {
	MailMemberService services.MailMemberService
}

func (mh *MailMemberHandler) Create(c *fiber.Ctx) error {
	c.Accepts("application/json")
	mailMember := new(models.MailMember)
	err := c.BodyParser(mailMember)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}
	
	err = mh.MailMemberService.Create(mailMember)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (mh *MailMemberHandler) FindAll(c *fiber.Ctx) error {
	mailMembers, err := mh.MailMemberService.FindAll()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": mailMembers})
}

func (mh *MailMemberHandler) DeleteById(c *fiber.Ctx) error {
	id := c.Params("id")

	monId := bson2.ObjectIdHex(id)

	err := mh.MailMemberService.DeleteById(monId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}
