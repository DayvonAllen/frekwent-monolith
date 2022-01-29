package services

import (
	"freq/models"
	"freq/repository"
	bson2 "github.com/globalsign/mgo/bson"
)

type ProductService interface {
	Create(product *models.Product) error
	FindAll(string, bool, bool) (*models.ProductList, error)
	FindAllByCategory(string, string, bool) (*models.ProductList, error)
	FindAllByProductIds(*[]bson2.ObjectId) (*[]models.Product, error)
	FindByProductId(bson2.ObjectId) (*models.Product, error)
	FindByProductName(string) (*models.Product, error)
	UpdateName(string, bson2.ObjectId) (*models.Product, error)
	UpdateQuantity(uint16, bson2.ObjectId) (*models.Product, error)
	UpdatePrice(string, bson2.ObjectId) (*models.Product, error)
	UpdateDescription(string, bson2.ObjectId) (*models.Product, error)
	UpdateIngredients(*[]string, bson2.ObjectId) (*models.Product, error)
	UpdateCategory(string, bson2.ObjectId) (*models.Product, error)
	DeleteById(bson2.ObjectId) error
}

type DefaultProductService struct {
	repo repository.ProductRepo
}

func (p DefaultProductService) Create(product *models.Product) error {
	err := p.repo.Create(product)

	if err != nil {
		return err
	}

	return nil
}

func (p DefaultProductService) FindAll(page string, newQuery bool, trending bool) (*models.ProductList, error) {
	products, err := p.repo.FindAll(page, newQuery, trending)

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p DefaultProductService) FindAllByProductIds(id *[]bson2.ObjectId) (*[]models.Product, error) {
	products, err := p.repo.FindAllByProductIds(id)

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p DefaultProductService) FindAllByCategory(category string, page string, newQuery bool) (*models.ProductList, error) {
	products, err := p.repo.FindAllByCategory(category, page, newQuery)

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p DefaultProductService) FindByProductId(id bson2.ObjectId) (*models.Product, error) {
	product, err := p.repo.FindByProductId(id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p DefaultProductService) FindByProductName(name string) (*models.Product, error) {
	product, err := p.repo.FindByProductName(name)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p DefaultProductService) UpdateName(name string, id bson2.ObjectId) (*models.Product, error) {
	product, err := p.repo.UpdateName(name, id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p DefaultProductService) UpdateQuantity(quan uint16, id bson2.ObjectId) (*models.Product, error) {
	product, err := p.repo.UpdateQuantity(quan, id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p DefaultProductService) UpdatePrice(price string, id bson2.ObjectId) (*models.Product, error) {
	product, err := p.repo.UpdatePrice(price, id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p DefaultProductService) UpdateDescription(desc string, id bson2.ObjectId) (*models.Product, error) {
	product, err := p.repo.UpdateDescription(desc, id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p DefaultProductService) UpdateIngredients(ingredients *[]string, id bson2.ObjectId) (*models.Product, error) {
	product, err := p.repo.UpdateIngredients(ingredients, id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p DefaultProductService) UpdateCategory(category string, id bson2.ObjectId) (*models.Product, error) {
	product, err := p.repo.UpdateCategory(category, id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p DefaultProductService) DeleteById(id bson2.ObjectId) error {
	err := p.repo.DeleteById(id)

	if err != nil {
		return err
	}

	return nil
}

func NewProductService(repository repository.ProductRepo) DefaultProductService {
	return DefaultProductService{repository}
}
