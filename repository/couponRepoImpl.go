package repository

import (
	"errors"
	"fmt"
	"freq/database"
	"freq/models"
	bson2 "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"time"
)

type CouponRepoImpl struct {
	Coupon     models.Coupon
	Coupons    []models.Coupon
	CouponList models.CouponList
}

func (c CouponRepoImpl) Create(coupon *models.Coupon) error {
	conn := database.Sess

	_, err := c.FindByCode(coupon.Code)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err.Error() == "not found" {
			coupon.CreatedAt = time.Now()
			coupon.UpdatedAt = time.Now()
			coupon.Id = bson2.NewObjectId()

			err = conn.DB(database.DB).C(database.COUPONS).Insert(coupon)

			if err != nil {
				return fmt.Errorf("error processing data")
			}

			return nil
		}
		return fmt.Errorf("error processing data")
	}

	return errors.New("coupon already exists")
}

func (c CouponRepoImpl) FindAll(page string, newCouponQuery bool) (*models.CouponList, error) {
	conn := database.Sess

	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}

	if newCouponQuery {
		//findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	err = conn.DB(database.DB).C(database.COUPONS).Find(nil).Skip((pageNumber - 1) * perPage).Limit(perPage).All(&c.Coupons)

	if err != nil {
		return nil, err
	}

	count, err := conn.DB(database.DB).C(database.COUPONS).Count()

	if err != nil {
		panic(err)
	}

	c.CouponList.NumberOfCoupons = count

	if c.CouponList.NumberOfCoupons < 10 {
		c.CouponList.NumberOfPages = 1
	} else {
		c.CouponList.NumberOfPages = int(count/10) + 1
	}

	c.CouponList.Coupons = &c.Coupons
	c.CouponList.CurrentPage = pageNumber

	return &c.CouponList, nil
}

func (c CouponRepoImpl) FindByCode(code string) (*models.Coupon, error) {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.COUPONS).Find(bson.M{"code": code}).One(&c.Coupon)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err.Error() == "not found" {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	return &c.Coupon, nil
}

func (c CouponRepoImpl) DeleteByCode(code string) error {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.COUPONS).Remove(bson.M{"code": code})

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func NewCouponRepoImpl() CouponRepoImpl {
	var couponRepoImpl CouponRepoImpl

	return couponRepoImpl
}
