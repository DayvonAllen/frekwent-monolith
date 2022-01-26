package services

import (
	"fmt"
	"freq/models"
	"freq/repository"
)

type AuthService interface {
	Login(username string, password string, ip string, ips *[]string) (*models.User, string, error)
}

type DefaultAuthService struct {
	repo repository.AuthRepo
}

func (a DefaultAuthService) Login(username string, password string, ip string, ips *[]string) (*models.User, string, error) {
	u, token, err := a.repo.Login(username, password, ip, ips)
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}
	return u, token, nil
}

func NewAuthService(repository repository.AuthRepo) DefaultAuthService {
	return DefaultAuthService{repository}
}
