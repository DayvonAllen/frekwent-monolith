package services

import (
	"freq/models"
	"freq/repository"
)

type EmailService interface {
	Create(email *models.Email) error
	SendMassEmail(emails *[]string, coupon string) error
	FindAll(string, bool) (*models.EmailList, error)
	FindAllByEmail(string, bool, string) (*models.EmailList, error)
	FindAllByStatus(string, bool, *models.Status) (*models.EmailList, error)
}

func (e DefaultEmailService) Create(email *models.Email) error {
	err := e.repo.Create(email)

	if err != nil {
		return err
	}

	return nil
}

func (e DefaultEmailService) SendMassEmail(emails *[]string, coupon string) error {
	err := e.repo.SendMassEmail(emails, coupon)

	if err != nil {
		return err
	}

	return nil
}

func (e DefaultEmailService) FindAll(page string, newQuery bool) (*models.EmailList, error) {
	emails, err := e.repo.FindAll(page, newQuery)

	if err != nil {
		return nil, err
	}

	return emails, nil
}

func (e DefaultEmailService) FindAllByEmail(page string, newQuery bool, email string) (*models.EmailList, error) {
	emails, err := e.repo.FindAllByEmail(page, newQuery, email)

	if err != nil {
		return nil, err
	}

	return emails, nil
}

func (e DefaultEmailService) FindAllByStatus(page string, newQuery bool, status *models.Status) (*models.EmailList, error) {
	emails, err := e.repo.FindAllByStatus(page, newQuery, status)

	if err != nil {
		return nil, err
	}

	return emails, nil
}

type DefaultEmailService struct {
	repo repository.EmailRepo
}

func NewEmailService(repository repository.EmailRepo) DefaultEmailService {
	return DefaultEmailService{repository}
}
