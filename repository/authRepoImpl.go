package repository

import (
	"fmt"
	"freq/database"
	"freq/helper"
	"freq/models"
	"freq/util"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepoImpl struct {
	user models.User
}

func (a AuthRepoImpl) Login(username string, password string, ip string, ips *[]string) (*models.User, string, error) {
	var login util.Authentication

	conn := database.Sess

	if helper.IsEmail(username) {
		err := conn.DB(database.DB).C(database.ADMIN).Find(bson.M{"email": username}).One(&a.user)

		if err != nil {
			return nil, "", fmt.Errorf("error finding by email")
		}

	} else {
		err := conn.DB(database.DB).C(database.ADMIN).Find(bson.M{"username": username}).One(&a.user)
		if err != nil {
			return nil, "", fmt.Errorf("error finding by username")
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(a.user.Password), []byte(password))

	if err != nil {
		return nil, "", fmt.Errorf("error comparing password")
	}

	token, err := login.GenerateJWT(a.user)

	if err != nil {
		return nil, "", fmt.Errorf("error generating token")
	}

	ipAddress := new(models.LoginIP)

	ipAddress.IpAddress = ip
	ipAddress.IpAddresses = *ips

	go func() {
		err = LoginIpRepoImpl{}.Create(ipAddress)
		if err != nil {
			return
		}
		return
	}()

	return &a.user, token, nil
}

func NewAuthRepoImpl() AuthRepoImpl {
	var authRepoImpl AuthRepoImpl

	return authRepoImpl
}
