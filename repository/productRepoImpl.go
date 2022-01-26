package repository

import (
	"errors"
	"fmt"
	"freq/database"
	"freq/models"
	"github.com/globalsign/mgo"
	bson2 "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

type ProductRepoImpl struct {
	product     models.Product
	products    []models.Product
	productList models.ProductList
}

func (p ProductRepoImpl) Create(product *models.Product) error {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.PRODUCTS).Find(bson.M{"name": product.Name}).One(&p.product)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err.Error() == "not found" {
			product.CreatedAt = time.Now()
			product.UpdatedAt = time.Now()
			product.Id = bson2.NewObjectId()

			err = conn.DB(database.DB).C(database.PRODUCTS).Insert(product)

			if err != nil {
				return fmt.Errorf("error processing data")
			}

			return nil
		}
		return fmt.Errorf("error processing data")
	}

	return errors.New("product with that name already exists")
}

func (p ProductRepoImpl) FindAll(page string, newProductQuery bool, trending bool) (*models.ProductList, error) {
	conn := database.Sess

	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}

	if trending {
		//findOptions.SetSort(bson.D{{"timesPurchased", -1}})
		//findOptions.SetLimit(3)
		//p.productList.NumberOfProducts = 3
	} else if newProductQuery {
		//findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	err = conn.DB(database.DB).C(database.PRODUCTS).Find(nil).Skip((pageNumber - 1) * perPage).Limit(perPage).All(&p.products)

	if err != nil {
		return nil, errors.New("error getting all products")
	}

	if p.productList.NumberOfProducts != 3 {
		count, err := conn.DB(database.DB).C(database.PRODUCTS).Count()

		if err != nil {
			return nil, errors.New("error getting all products")
		}

		p.productList.NumberOfProducts = count
	}

	if p.productList.NumberOfProducts < 10 {
		p.productList.NumberOfPages = 1
	} else {
		p.productList.NumberOfPages = int(p.productList.NumberOfProducts/10) + 1
	}

	p.productList.Products = &p.products
	p.productList.CurrentPage = pageNumber

	return &p.productList, nil
}

func (p ProductRepoImpl) FindAllByProductIds(ids *[]bson2.ObjectId) (*[]models.Product, error) {
	conn := database.Sess

	for _, id := range *ids {
		err := conn.DB(database.DB).C(database.PRODUCTS).Find(bson.M{"_id": id}).One(&p.product)
		if err != nil {
			return nil, errors.New("not found")
		}

		p.products = append(p.products, p.product)
	}

	return &p.products, nil
}

func (p ProductRepoImpl) FindAllByCategory(category string, page string, newProductQuery bool) (*models.ProductList, error) {
	conn := database.Sess

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}

	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	if newProductQuery {
		//findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	err = conn.DB(database.DB).C(database.PRODUCTS).Find(bson.M{"category": category}).All(p.products)

	if err != nil {
		return nil, errors.New("error finding by category")
	}

	count, err := conn.DB(database.DB).C(database.IPS).Count()

	if err != nil {
		return nil, errors.New("error getting all by category")
	}

	p.productList.NumberOfProducts = count

	if p.productList.NumberOfProducts < 10 {
		p.productList.NumberOfPages = 1
	} else {
		p.productList.NumberOfPages = int(count/10) + 1
	}

	p.productList.Products = &p.products
	p.productList.CurrentPage = pageNumber

	return &p.productList, nil
}

func (p ProductRepoImpl) FindByProductId(id bson2.ObjectId) (*models.Product, error) {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.PRODUCTS).FindId(id).One(&p.product)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err.Error() == "not found" {
			return nil, errors.New("not found")
		}
		return nil, fmt.Errorf("error processing data")
	}

	return &p.product, nil
}

func (p ProductRepoImpl) FindByProductName(name string) (*models.Product, error) {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.PRODUCTS).Find(bson.M{"name": name}).One(&p.product)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err.Error() == "not found" {
			return nil, errors.New("not found")
		}
		return nil, fmt.Errorf("error processing data")
	}

	return &p.product, nil
}

func (p ProductRepoImpl) UpdateName(name string, id bson2.ObjectId) (*models.Product, error) {
	conn := database.Sess

	prod := new(models.Product)

	err := conn.DB(database.DB).C(database.PRODUCTS).Find(bson.D{{"name", name}}).One(prod)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err.Error() == "not found" {
			update := bson.M{"name": name,
				"updatedAt": time.Now()}

			err = conn.DB(database.DB).C(database.PRODUCTS).UpdateId(id, update)

			if err != nil {
				return nil, errors.New("error updating by name")
			}

			p.product.Name = name

			return &p.product, nil
		}
		return nil, fmt.Errorf("error processing data")
	}

	return nil, errors.New("product with that name already exists")
}

func (p ProductRepoImpl) UpdateQuantity(quantity uint16, id bson2.ObjectId) (*models.Product, error) {
	conn := database.Sess

	update := bson.M{"quantity": quantity,
		"updatedAt": time.Now()}

	err := conn.DB(database.DB).C(database.PRODUCTS).UpdateId(id, update)

	if err != nil {
		return nil, errors.New("error updating quantity")
	}

	p.product.Quantity = quantity

	return &p.product, nil
}

func (p ProductRepoImpl) UpdatePurchaseCount(name string) error {
	conn := database.Sess

	update :=
		mgo.Change{
			Update: bson.M{"$inc": bson.M{"timesPurchased": 1}},
		}

	_, err := conn.DB(database.DB).C(database.PRODUCTS).Find(bson.M{"name": name}).Apply(update, &p.product)

	if err != nil {
		return errors.New("error updating purchase count")
	}

	return nil
}

func (p ProductRepoImpl) UpdatePrice(price string, id bson2.ObjectId) (*models.Product, error) {
	conn := database.Sess

	update := bson.M{"price": price,
		"updatedAt": time.Now()}

	err := conn.DB(database.DB).C(database.PRODUCTS).UpdateId(id, update)

	if err != nil {
		return nil, errors.New("error updating price")
	}

	p.product.Price = price

	return &p.product, nil
}

func (p ProductRepoImpl) UpdateDescription(desc string, id bson2.ObjectId) (*models.Product, error) {
	conn := database.Sess

	update := bson.M{"description": desc,
		"updatedAt": time.Now()}

	err := conn.DB(database.DB).C(database.PRODUCTS).UpdateId(id, update)

	if err != nil {
		return nil, errors.New("error updating description")
	}

	p.product.Description = desc

	return &p.product, nil
}

func (p ProductRepoImpl) UpdateIngredients(ingredients *[]string, id bson2.ObjectId) (*models.Product, error) {
	conn := database.Sess

	update := bson.M{"ingredients": ingredients,
		"updatedAt": time.Now()}

	err := conn.DB(database.DB).C(database.PRODUCTS).UpdateId(id, update)

	if err != nil {
		return nil, errors.New("error updating ingredients")
	}

	p.product.Ingredients = *ingredients

	return &p.product, nil
}

func (p ProductRepoImpl) UpdateCategory(category string, id bson2.ObjectId) (*models.Product, error) {
	conn := database.Sess

	update := bson.M{"category": category,
		"updatedAt": time.Now()}

	err := conn.DB(database.DB).C(database.PRODUCTS).UpdateId(id, update)

	if err != nil {
		return nil, errors.New("error updating category")
	}

	p.product.Category = category

	return &p.product, nil
}

func (p ProductRepoImpl) DeleteById(id bson2.ObjectId) error {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.PRODUCTS).RemoveId(id)

	if err != nil {
		return errors.New("error deleting by ID")
	}

	return nil
}

func NewProductRepoImpl() ProductRepoImpl {
	var productRepoImpl ProductRepoImpl

	return productRepoImpl
}
