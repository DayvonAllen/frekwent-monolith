package repository

import (
	"errors"
	"fmt"
	"freq/config"
	"freq/database"
	"freq/helper"
	"freq/models"
	bson2 "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PurchaseRepoImpl struct {
	purchase     models.Purchase
	purchases    []models.Purchase
	purchaseList models.PurchaseList
	transaction  models.Transactions
}

func (p PurchaseRepoImpl) Purchase(purchase *models.Purchase) error {
	conn := database.Sess

	purchase = helper.EncryptPI(purchase)

	customer := new(models.Customer)

	customer.Id = bson2.NewObjectId()
	customer.UpdatedAt = time.Now()
	customer.CreatedAt = time.Now()
	customer.FirstName = purchase.FirstName
	customer.LastName = purchase.LastName
	customer.Email = purchase.Email
	customer.PurchasedItems = purchase.PurchasedItems
	customer.StreetAddress = purchase.StreetAddress
	customer.OptionalAddress = purchase.OptionalAddress
	customer.City = purchase.City
	customer.State = purchase.State
	customer.ZipCode = purchase.ZipCode
	customer.InfoEmailOptIn = purchase.InfoEmailOptIn

	err := CustomerRepoImpl{}.Create(customer)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	purchase.CouponUsed.Id = bson2.NewObjectId()

	err = conn.DB(database.DB).C(database.PURCHASES).Insert(purchase)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	go func() {
		emailAdd := config.Config("BUSINESS_EMAIL")

		email := new(models.Email)
		email.Id = bson2.NewObjectId()
		email.CustomerEmail = customer.Email
		email.From = emailAdd
		email.Subject = "test subject"
		email.Content = "test content"
		email.Type = models.Purchased
		email.Status = models.Pending

		err := EmailRepoImpl{}.Create(email)

		if err != nil {
			log.Println(errors.New(fmt.Sprintf("error sending an email to %s", email.CustomerEmail)))
		}

		// Todo send email and update status to success or failure
	}()

	go func() {
		repo := ProductRepoImpl{}
		for _, product := range *purchase.PurchasedItems {
			err = repo.UpdatePurchaseCount(product.Name)
			if err != nil {
				return
			}
		}
	}()

	return nil
}

func (p PurchaseRepoImpl) FindAll(page string, newPurchaseQuery bool) (*models.PurchaseList, error) {
	conn := database.Sess

	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}

	if newPurchaseQuery {
		//findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	err = conn.DB(database.DB).C(database.PURCHASES).Find(nil).Skip((pageNumber - 1) * perPage).Limit(perPage).All(&p.purchases)

	if err != nil {
		return nil, errors.New("error finding purchases")
	}

	count, err := conn.DB(database.DB).C(database.PRODUCTS).Count()

	if err != nil {
		return nil, errors.New("error finding purchases")
	}

	p.purchaseList.NumberOfPurchases = count

	if p.purchaseList.NumberOfPurchases < 10 {
		p.purchaseList.NumberOfPages = 1
	} else {
		p.purchaseList.NumberOfPages = int(count/10) + 1
	}

	p.purchaseList.Purchases = &p.purchases
	p.purchaseList.CurrentPage = pageNumber

	return &p.purchaseList, nil
}

func (p PurchaseRepoImpl) CalculateTransactionsByState(state string) (*models.Transactions, error) {
	conn := database.Sess

	var filter bson.M

	if strings.ToLower(state) == "all" {
		filter = bson.M{"refunded": false}
	} else {
		filter = bson.M{"state": state, "refunded": false}
	}

	err := conn.DB(database.DB).C(database.PURCHASES).Find(filter).All(&p.purchases)

	if err != nil {
		return nil, errors.New("error calculating by state")
	}

	count, err := conn.DB(database.DB).C(database.PURCHASES).Find(filter).Count()

	if err != nil {
		return nil, errors.New("error calculating by state")
	}

	for _, pur := range p.purchases {
		p.transaction.TransactionsTotal = pur.FinalPrice + p.transaction.TransactionsTotal
	}

	p.transaction.NumberOfTransactions = count

	return &p.transaction, nil
}

func (p PurchaseRepoImpl) FindByPurchaseById(id bson2.ObjectId) (*models.Purchase, error) {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.PURCHASES).FindId(id).One(&p.purchase)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err.Error() == "not found" {
			return nil, errors.New(fmt.Sprintf("no purchase with the ID %v", id))
		}
		return nil, fmt.Errorf("error processing data")
	}

	p.purchase = *helper.DecryptPI(&p.purchase)

	return &p.purchase, nil
}

func (p PurchaseRepoImpl) UpdateShippedStatus(dto *models.PurchaseShippedDTO) error {
	conn := database.Sess

	update := bson.M{"shipped": dto.Shipped,
		"trackingId": dto.TrackingId}

	err := conn.DB(database.DB).C(database.PURCHASES).UpdateId(dto.Id, update)

	if err != nil {
		return errors.New("error updating purchase's shipping status")
	}

	return nil
}

func (p PurchaseRepoImpl) UpdateDeliveredStatus(dto *models.PurchaseDeliveredDTO) error {
	conn := database.Sess

	update := bson.M{"delivered": dto.Delivered}

	err := conn.DB(database.DB).C(database.PURCHASES).UpdateId(dto.Id, update)

	if err != nil {
		return errors.New("error updating purchase's delivered status")
	}

	return nil
}

func (p PurchaseRepoImpl) UpdatePurchaseAddress(dto *models.PurchaseAddressDTO) error {
	conn := database.Sess

	key := config.Config("KEY")

	encrypt := helper.Encryption{Key: []byte(key)}

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(dto.StreetAddress)

		if err != nil {
			log.Println(err)
		}

		dto.StreetAddress = pi
	}()

	go func() {
		defer wg.Done()

		if len(dto.OptionalAddress) > 0 {
			pi, err := encrypt.Encrypt(dto.OptionalAddress)

			if err != nil {
				log.Println(err)
			}

			dto.OptionalAddress = pi
		}
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(dto.City)

		if err != nil {
			log.Println(err)
		}

		dto.City = pi
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(dto.State)

		if err != nil {
			log.Println(err)
		}

		dto.State = pi
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(dto.ZipCode)

		if err != nil {
			log.Println(err)
		}

		dto.ZipCode = pi
	}()

	wg.Wait()

	update := bson.M{"streetAddress": dto.StreetAddress,
		"optionalAddress": dto.OptionalAddress,
		"city":            dto.City,
		"state":           dto.State,
		"zipCode":         dto.ZipCode}

	err := conn.DB(database.DB).C(database.PURCHASES).UpdateId(dto.Id, update)

	if err != nil {
		return errors.New("error updating purchase's address")
	}

	return nil
}

func (p PurchaseRepoImpl) UpdateTrackingNumber(dto *models.PurchaseTrackingDTO) error {
	conn := database.Sess

	update := bson.M{"trackingId": dto.TrackingId}

	err := conn.DB(database.DB).C(database.PURCHASES).UpdateId(dto.Id, update)

	if err != nil {
		return errors.New("error updating purchase's tracking ID")
	}

	return nil
}

func NewPurchaseRepoImpl() PurchaseRepoImpl {
	var purchaseRepoImpl PurchaseRepoImpl

	return purchaseRepoImpl
}
