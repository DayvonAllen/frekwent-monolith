package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Product struct {
	Id             bson.ObjectId `bson:"_id" json:"id"`
	Name           string        `bson:"name" json:"name"`
	Images         []string      `bson:"images" json:"images"`
	Price          string        `bson:"price" json:"price"`
	Quantity       uint16        `bson:"quantity" json:"quantity"`
	Description    string        `bson:"description" json:"description"`
	Ingredients    []string      `bson:"ingredients" json:"ingredients"`
	Category       string        `bson:"category" json:"category"`
	TimesPurchased int           `bson:"timesPurchased" json:"timesPurchased"`
	CreatedAt      time.Time     `bson:"createdAt" json:"-"`
	UpdatedAt      time.Time     `bson:"updatedAt" json:"-"`
}

type ProductNameDto struct {
	Name string `bson:"name" json:"name"`
}

type ProductIdDto struct {
	Ids []bson.ObjectId `json:"ids"`
}

type ProductPriceDto struct {
	Price string `bson:"price" json:"price"`
}

type ProductQuantityDto struct {
	Quantity uint16 `bson:"quantity" json:"quantity"`
}

type ProductDescriptionDto struct {
	Description string `bson:"description" json:"description"`
}

type ProductIngredientsDto struct {
	Ingredients *[]string `bson:"ingredients" json:"ingredients"`
}

type ProductCategoryDto struct {
	Category string `bson:"category" json:"category"`
}

type ProductList struct {
	Products         *[]Product `json:"products"`
	NumberOfProducts int        `json:"numberOfProducts"`
	CurrentPage      int        `json:"currentPage"`
	NumberOfPages    int        `json:"numberOfPages"`
}
