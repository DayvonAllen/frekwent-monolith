package repository

import (
	"freq/models"
	bson2 "github.com/globalsign/mgo/bson"
)

type MailMemberRepo interface {
	Create(mailMember *models.MailMember) error
	FindAll() (*[]models.MailMember, error)
	DeleteById(id bson2.ObjectId) error
}
