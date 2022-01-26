package repository

import (
	"freq/models"
	bson2 "github.com/globalsign/mgo/bson"
)

type EmailRepo interface {
	Create(email *models.Email) error
	SendMassEmail(emails *[]string, coupon string) error
	FindAll(string, bool) (*models.EmailList, error)
	FindAllByEmail(string, bool, string) (*models.EmailList, error)
	FindAllByStatus(string, bool, *models.Status) (*models.EmailList, error)
	UpdateEmailStatus(bson2.ObjectId, models.Status) error
	FindAllByStatusRaw(status *models.Status) (*[]models.Email, error)
}
