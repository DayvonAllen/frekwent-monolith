package handlers

import (
	"errors"
	"fmt"
	"freq/helper"
	"freq/models"
	"freq/services"
	bson2 "github.com/globalsign/mgo/bson"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

type ProductHandler struct {
	ProductService services.ProductService
}

func (ph *ProductHandler) Create(c *fiber.Ctx) error {
	c.Accepts("application/json")
	product := new(models.Product)
	err := c.BodyParser(product)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = helper.ValidateData(product.Category, "Category", 1, 60)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = helper.ValidateData(product.Name, "Name", 1, 60)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = helper.ValidateData(product.Description, "Description", 1, 120)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	if !helper.IsValidPrice(product.Price) {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", errors.New("invalid price"))})
	}

	product.Category = strings.ToLower(product.Category)
	product.Name = strings.ToLower(product.Name)

	err = ph.ProductService.Create(product)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ph *ProductHandler) FindAll(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	newProductQuery := c.Query("new", "false")
	trending := c.Query("trending", "false")

	isNew, err := strconv.ParseBool(newProductQuery)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("must provide a valid value")})
	}

	isTrending, err := strconv.ParseBool(trending)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("must provide a valid value")})
	}

	products, err := ph.ProductService.FindAll(page, isNew, isTrending)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": products})
}

func (ph *ProductHandler) FindAllByCategory(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	newProductQuery := c.Query("new", "false")
	category := c.Params("category")

	isNew, err := strconv.ParseBool(newProductQuery)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("must provide a valid value")})
	}

	products, err := ph.ProductService.FindAllByCategory(strings.ToLower(category), page, isNew)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": products})
}

func (ph *ProductHandler) FindByProductId(c *fiber.Ctx) error {
	id := c.Params("id")

	monId := bson2.ObjectIdHex(id)

	product, err := ph.ProductService.FindByProductId(monId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": product})
}

func (ph *ProductHandler) FindAllByProductIds(c *fiber.Ctx) error {
	product := new(models.ProductIdDto)
	err := c.BodyParser(product)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	products, err := ph.ProductService.FindAllByProductIds(&product.Ids)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": products})
}

func (ph *ProductHandler) FindByProductName(c *fiber.Ctx) error {
	name := c.Query("name")

	product, err := ph.ProductService.FindByProductName(name)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": product})
}

func (ph *ProductHandler) UpdateName(c *fiber.Ctx) error {
	c.Accepts("application/json")
	product := new(models.ProductNameDto)
	err := c.BodyParser(product)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = helper.ValidateData(product.Name, "Name", 1, 60)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	id := c.Params("id")

	monId := bson2.ObjectIdHex(id)

	updatedProduct, err := ph.ProductService.UpdateName(product.Name, monId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": updatedProduct})
}

func (ph *ProductHandler) UpdatePrice(c *fiber.Ctx) error {
	c.Accepts("application/json")
	product := new(models.ProductPriceDto)
	err := c.BodyParser(product)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	id := c.Params("id")

	monId := bson2.ObjectIdHex(id)

	if !helper.IsValidPrice(product.Price) {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", errors.New("invalid price"))})
	}

	updatedProduct, err := ph.ProductService.UpdatePrice(product.Price, monId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": updatedProduct})
}

func (ph *ProductHandler) UpdateDescription(c *fiber.Ctx) error {
	c.Accepts("application/json")
	product := new(models.ProductDescriptionDto)
	err := c.BodyParser(product)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = helper.ValidateData(product.Description, "Description", 1, 120)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	id := c.Params("id")

	monId := bson2.ObjectIdHex(id)

	updatedProduct, err := ph.ProductService.UpdateDescription(product.Description, monId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": updatedProduct})
}

func (ph *ProductHandler) UpdateQuantity(c *fiber.Ctx) error {
	c.Accepts("application/json")
	product := new(models.ProductQuantityDto)
	err := c.BodyParser(product)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	id := c.Params("id")

	monId := bson2.ObjectIdHex(id)

	updatedProduct, err := ph.ProductService.UpdateQuantity(product.Quantity, monId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": updatedProduct})
}

func (ph *ProductHandler) UpdateIngredients(c *fiber.Ctx) error {
	c.Accepts("application/json")
	product := new(models.ProductIngredientsDto)
	err := c.BodyParser(product)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	id := c.Params("id")

	monId := bson2.ObjectIdHex(id)

	updatedProduct, err := ph.ProductService.UpdateIngredients(product.Ingredients, monId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": updatedProduct})
}

func (ph *ProductHandler) UpdateCategory(c *fiber.Ctx) error {
	c.Accepts("application/json")
	product := new(models.ProductCategoryDto)
	err := c.BodyParser(product)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = helper.ValidateData(product.Category, "Category", 1, 60)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	id := c.Params("id")

	monId := bson2.ObjectIdHex(id)

	updatedProduct, err := ph.ProductService.UpdateCategory(strings.ToLower(product.Category), monId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": updatedProduct})
}

func (ph *ProductHandler) DeleteById(c *fiber.Ctx) error {
	id := c.Params("id")

	monId := bson2.ObjectIdHex(id)

	err := ph.ProductService.DeleteById(monId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(204).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}
