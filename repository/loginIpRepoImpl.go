package repository

import (
	"fmt"
	"freq/database"
	"freq/models"
	bson2 "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"time"
)

type LoginIpRepoImpl struct {
	loginIp     models.LoginIP
	loginIps    []models.LoginIP
	loginIpList models.LoginIpList
}

func (l LoginIpRepoImpl) Create(ip *models.LoginIP) error {
	_, err := l.FindByIp(ip.IpAddress)

	if err != nil {
		conn := database.Sess

		ip.AccessedAt = time.Now()
		ip.CreatedAt = time.Now()
		ip.UpdatedAt = time.Now()
		ip.Id = bson2.NewObjectId()

		err := conn.DB(database.DB).C(database.IPS).Insert(&ip)

		if err != nil {
			return fmt.Errorf("error processing data")
		}
	}

	err = l.UpdateLoginIp(ip)

	if err != nil {
		return err
	}

	return nil
}

func (l LoginIpRepoImpl) FindAll(page string, newLoginQuery bool) (*models.LoginIpList, error) {
	conn := database.Sess

	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}

	if newLoginQuery {
		//findOptions.SetSort(bson.D{{"updatedAt", -1}})
	}

	err = conn.DB(database.DB).C(database.IPS).Find(nil).Skip((pageNumber - 1) * perPage).Limit(perPage).All(&l.loginIps)

	if err != nil {
		return nil, err
	}

	count, err := conn.DB(database.DB).C(database.IPS).Count()

	if err != nil {
		panic(err)
	}

	l.loginIpList.NumberOfLoginIps = count

	if l.loginIpList.NumberOfLoginIps < 10 {
		l.loginIpList.NumberOfPages = 1
	} else {
		l.loginIpList.NumberOfPages = int(count/10) + 1
	}

	l.loginIpList.LoginIps = &l.loginIps
	l.loginIpList.CurrentPage = pageNumber

	return &l.loginIpList, nil
}

func (l LoginIpRepoImpl) FindByIp(ip string) (*models.LoginIP, error) {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.IPS).Find(bson.M{"ipAddress": ip}).One(&l.loginIp)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err.Error() == "not found" {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	return &l.loginIp, nil
}

func (l LoginIpRepoImpl) UpdateLoginIp(ip *models.LoginIP) error {
	conn := database.Sess

	ip.AccessedAt = time.Now()
	ip.UpdatedAt = time.Now()

	err := conn.DB(database.DB).C(database.IPS).UpdateId(ip.Id, ip)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err.Error() == "not found" {
			return err
		}
		return fmt.Errorf("error processing data")
	}

	return nil
}

func NewLoginIpRepoImpl() LoginIpRepoImpl {
	var loginIpRepoImpl LoginIpRepoImpl

	return loginIpRepoImpl
}
