package models

import (
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Purchase struct {
	Id              bson.ObjectId `bson:"_id" json:"id"`
	FirstName       string        `bson:"firstName" json:"firstName"`
	LastName        string        `bson:"lastName" json:"lastName"`
	Email           string        `bson:"email" json:"email"`
	StreetAddress   string        `bson:"streetAddress" json:"streetAddress"`
	OptionalAddress string        `bson:"optionalAddress" json:"optionalAddress"`
	City            string        `bson:"city" json:"city"`
	State           string        `bson:"state" json:"state"`
	ZipCode         string        `bson:"zipCode" json:"zipCode"`
	PurchasedItems  *[]Product    `bson:"purchasedItems" json:"purchasedItems"`
	FinalPrice      float32       `bson:"finalPrice" json:"finalPrice"`
	CouponUsed      Coupon        `bson:"couponUsed" json:"couponUsed"`
	PurchaseIntent  string        `bson:"purchaseIntent" json:"purchaseIntent"`
	TrackingId      string        `bson:"trackingId" json:"trackingId"`
	Shipped         bool          `bson:"shipped" json:"shipped"`
	Delivered       bool          `bson:"delivered" json:"delivered"`
	InfoEmailOptIn  bool          `bson:"infoEmailOptIn" json:"infoEmailOptIn"`
	Refunded        bool          `bson:"refunded" json:"refunded"`
	CreatedAt       time.Time     `bson:"createdAt" json:"-"`
	UpdatedAt       time.Time     `bson:"updatedAt" json:"-"`
}

type PurchaseAddressDTO struct {
	Id              primitive.ObjectID `json:"id"`
	StreetAddress   string             `json:"streetAddress"`
	OptionalAddress string             `json:"optionalAddress"`
	City            string             `json:"city"`
	State           string             `json:"state"`
	ZipCode         string             `json:"zipCode"`
}

type PurchaseDeliveredDTO struct {
	Id        primitive.ObjectID `json:"id"`
	Delivered bool               `json:"delivered"`
}

type PurchaseShippedDTO struct {
	Id         primitive.ObjectID `json:"id"`
	Shipped    bool               `json:"shipped"`
	TrackingId string             `json:"trackingId"`
}

type PurchaseTrackingDTO struct {
	Id         primitive.ObjectID `json:"id"`
	TrackingId string             `json:"trackingId"`
}

type PurchaseList struct {
	Purchases         *[]Purchase `json:"purchases"`
	NumberOfPurchases int         `json:"numberOfPurchases"`
	CurrentPage       int         `json:"currentPage"`
	NumberOfPages     int         `json:"numberOfPages"`
}

type Transactions struct {
	TransactionsTotal    float32 `json:"transactionsTotal"`
	NumberOfTransactions int     `json:"numberOfTransactions"`
}
