package repository

import "freq/models"

type LoginIpRepo interface {
	Create(ip *models.LoginIP) error
	FindAll(string, bool) (*models.LoginIpList, error)
	FindByIp(string) (*models.LoginIP, error)
	UpdateLoginIp(ip *models.LoginIP) error
}
