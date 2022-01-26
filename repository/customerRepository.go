package repository

import "freq/models"

type CustomerRepo interface {
	Create(customer *models.Customer) error
	FindAll(string, bool) (*models.CustomerList, error)
	FindAllByFullName(string, string, string, bool) (*models.CustomerList, error)
	FindByEmail(string) (*models.Customer, error)
	FindAllByOptInStatus(bool) (*[]models.Customer, error)
	UpdateOptInStatus(bool, string) (*models.Customer, error)
}
