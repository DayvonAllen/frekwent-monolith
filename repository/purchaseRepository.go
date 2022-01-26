package repository

import (
	"freq/models"
	bson2 "github.com/globalsign/mgo/bson"
)

type PurchaseRepo interface {
	Purchase(purchase *models.Purchase) error
	FindAll(string, bool) (*models.PurchaseList, error)
	FindByPurchaseById(bson2.ObjectId) (*models.Purchase, error)
	CalculateTransactionsByState(string) (*models.Transactions, error)
	UpdateShippedStatus(*models.PurchaseShippedDTO) error
	UpdateDeliveredStatus(*models.PurchaseDeliveredDTO) error
	UpdatePurchaseAddress(*models.PurchaseAddressDTO) error
	UpdateTrackingNumber(*models.PurchaseTrackingDTO) error
}
