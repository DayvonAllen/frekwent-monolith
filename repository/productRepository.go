package repository

import (
	"freq/models"
	bson2 "github.com/globalsign/mgo/bson"
)

type ProductRepo interface {
	Create(product *models.Product) error
	FindAll(string, bool, bool) (*models.ProductList, error)
	FindByProductId(bson2.ObjectId) (*models.Product, error)
	FindByProductName(string) (*models.Product, error)
	FindAllByProductIds(*[]bson2.ObjectId) (*[]models.Product, error)
	FindAllByCategory(string, string, bool) (*models.ProductList, error)
	UpdatePurchaseCount(string) error
	UpdateName(string, bson2.ObjectId) (*models.Product, error)
	UpdateQuantity(uint16, bson2.ObjectId) (*models.Product, error)
	UpdatePrice(string, bson2.ObjectId) (*models.Product, error)
	UpdateDescription(string, bson2.ObjectId) (*models.Product, error)
	UpdateIngredients(*[]string, bson2.ObjectId) (*models.Product, error)
	UpdateCategory(string, bson2.ObjectId) (*models.Product, error)
	DeleteById(bson2.ObjectId) error
}
