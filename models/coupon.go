package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Coupon struct {
	Id             bson.ObjectId `bson:"_id" json:"id"`
	Percentage     uint8         `bson:"percentage" json:"percentage"`
	Code           string        `bson:"code" json:"code"`
	ExpirationDate time.Time     `bson:"expirationDate" json:"expirationDate"`
	CreatedAt      time.Time     `bson:"createdAt" json:"-"`
	UpdatedAt      time.Time     `bson:"updatedAt" json:"-"`
}

type CouponList struct {
	Coupons         *[]Coupon `json:"coupons"`
	NumberOfCoupons int       `json:"numberOfCoupons"`
	CurrentPage     int       `json:"currentPage"`
	NumberOfPages   int       `json:"numberOfPages"`
}
