package handlers

import (
	"fmt"
	"freq/services"
	"freq/util"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type AuthHandler struct {
	AuthService services.AuthService
}

func (ah *AuthHandler) Login(c *fiber.Ctx) error {
	c.Accepts("application/json")
	details := new(util.LoginDetails)
	err := c.BodyParser(details)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	if len(details.Email) <= 1 || len(details.Email) > 60 {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("invalid username")})
	}

	if len(details.Password) <= 5 || len(details.Password) > 60 {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("invalid password")})
	}

	var auth util.Authentication

	ips := c.IPs()

	user, token, err := ah.AuthService.Login(strings.ToLower(details.Email), details.Password, c.IP(), &ips)

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	signedToken := make([]byte, 0, 100)
	signedToken = append(signedToken, []byte("Bearer "+token+"|")...)
	t, err := auth.SignToken([]byte(token))

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	signedToken = append(signedToken, t...)

	cookie := new(fiber.Cookie)
	cookie.Name = "auth"
	cookie.Value = string(signedToken)
	cookie.HTTPOnly = true
	cookie.Expires = time.Now().Add(24 * time.Hour)

	// Set cookie
	c.Cookie(cookie)

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": user})
}

func (ah *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name: "auth",
		// Set expiry date to the past
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: true,
	})

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}
