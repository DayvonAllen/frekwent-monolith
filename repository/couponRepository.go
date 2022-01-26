package repository

import (
	"freq/models"
)

type CouponRepo interface {
	Create(coupon *models.Coupon) error
	FindAll(string, bool) (*models.CouponList, error)
	FindByCode(string) (*models.Coupon, error)
	DeleteByCode(string) error
}
