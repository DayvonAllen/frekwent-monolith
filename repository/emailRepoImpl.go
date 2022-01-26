package repository

import (
	"errors"
	"fmt"
	"freq/config"
	"freq/database"
	"freq/models"
	bson2 "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"time"
)

type EmailRepoImpl struct {
	email     models.Email
	emails    []models.Email
	emailList models.EmailList
}

func (e EmailRepoImpl) Create(email *models.Email) error {
	conn := database.Sess

	email.CreatedAt = time.Now()
	email.UpdatedAt = time.Now()

	err := conn.DB(database.DB).C(database.EMAILS).Insert(&email)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (e EmailRepoImpl) SendMassEmail(emails *[]string, coupon string) error {
	conn := database.Sess

	emailsArr := make([]interface{}, 0, len(*emails))

	for _, em := range *emails {
		emailsArr = append(emailsArr, models.Email{
			Id:            bson2.NewObjectId(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Content:       coupon,
			Subject:       "coupon status",
			Status:        models.Pending,
			Type:          models.CouponPromotion,
			From:          config.Config("BUSINESS_EMAIL"),
			CustomerEmail: em,
		})
	}

	err := conn.DB(database.DB).C(database.EMAILS).Insert(emailsArr)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (e EmailRepoImpl) FindAll(page string, newEmailQuery bool) (*models.EmailList, error) {
	conn := database.Sess

	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}

	if newEmailQuery {
		//findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	err = conn.DB(database.DB).C(database.EMAILS).Find(nil).Skip((pageNumber - 1) * perPage).Limit(perPage).All(&e.emails)

	if err != nil {
		return nil, err
	}

	count, err := conn.DB(database.DB).C(database.EMAILS).Count()

	if err != nil {
		panic(err)
	}

	e.emailList.NumberOfEmails = count

	if e.emailList.NumberOfEmails < 10 {
		e.emailList.NumberOfPages = 1
	} else {
		e.emailList.NumberOfPages = int(count/10) + 1
	}

	e.emailList.Emails = &e.emails
	e.emailList.CurrentPage = pageNumber

	return &e.emailList, nil
}

func (e EmailRepoImpl) FindAllByEmail(page string, newEmailQuery bool, email string) (*models.EmailList, error) {
	conn := database.Sess

	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}

	if newEmailQuery {
		//findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	err = conn.DB(database.DB).C(database.EMAILS).Find(bson.M{"email": email}).Skip((pageNumber - 1) * perPage).Limit(perPage).All(&e.emails)

	if err != nil {
		return nil, errors.New("error finding email")
	}

	count, err := conn.DB(database.DB).C(database.EMAILS).Find(bson.M{"email": email}).Count()

	if err != nil {
		panic(err)
	}

	e.emailList.NumberOfEmails = count

	if e.emailList.NumberOfEmails < 10 {
		e.emailList.NumberOfPages = 1
	} else {
		e.emailList.NumberOfPages = int(count/10) + 1
	}

	e.emailList.Emails = &e.emails
	e.emailList.CurrentPage = pageNumber

	return &e.emailList, nil
}

func (e EmailRepoImpl) FindAllByStatus(page string, newEmailQuery bool, status *models.Status) (*models.EmailList, error) {
	conn := database.Sess

	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}

	if newEmailQuery {
		//findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	err = conn.DB(database.DB).C(database.EMAILS).Find(bson.M{"status": status}).Skip((pageNumber - 1) * perPage).Limit(perPage).All(&e.emails)

	if err != nil {
		return nil, errors.New("error finding email")
	}

	count, err := conn.DB(database.DB).C(database.EMAILS).Find(bson.M{"status": status}).Count()

	if err != nil {
		panic(err)
	}

	e.emailList.NumberOfEmails = count

	if e.emailList.NumberOfEmails < 10 {
		e.emailList.NumberOfPages = 1
	} else {
		e.emailList.NumberOfPages = int(count/10) + 1
	}

	e.emailList.Emails = &e.emails
	e.emailList.CurrentPage = pageNumber

	return &e.emailList, nil
}

func (e EmailRepoImpl) UpdateEmailStatus(id bson2.ObjectId, status models.Status) error {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.EMAILS).UpdateId(id, bson.M{"updatedAt": time.Now(), "status": status})

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err.Error() == "not found" {
			return err
		}
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (e EmailRepoImpl) FindAllByStatusRaw(status *models.Status) (*[]models.Email, error) {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.EMAILS).Find(bson.M{"status": status}).All(&e.emails)

	if err != nil {
		return nil, err
	}

	return &e.emails, nil
}

func NewEmailRepoImpl() EmailRepoImpl {
	var emailRepoImpl EmailRepoImpl

	return emailRepoImpl
}
