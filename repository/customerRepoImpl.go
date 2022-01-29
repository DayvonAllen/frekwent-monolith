package repository

import (
	"errors"
	"fmt"
	"freq/config"
	"freq/database"
	"freq/helper"
	"freq/models"
	"github.com/globalsign/mgo"
	bson2 "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"sync"
	"time"
)

type CustomerRepoImpl struct {
	customer     models.Customer
	customers    []models.Customer
	customerList models.CustomerList
}

func (c CustomerRepoImpl) Create(customer *models.Customer) error {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.CUSTOMERS).Find(bson.M{"email": customer.Email}).One(&c.customer)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err.Error() == "not found" {
			// Index
			index := mgo.Index{
				Key:        []string{"firstName", "lastName", "email", "infoEmailOptIn"},
				Unique:     true,
				DropDups:   true,
				Background: true,
				Sparse:     true,
			}

			coll := conn.DB(database.DB).C(database.CUSTOMERS)
			err = coll.EnsureIndex(index)

			if err != nil {
				panic(err)
			}

			err = coll.Insert(customer)

			if err != nil {
				return fmt.Errorf("error processing data")
			}

			go func() {
				err = MailMemberRepoImpl{}.Create(&models.MailMember{
					Id:              bson2.NewObjectId(),
					MemberFirstName: customer.FirstName,
					MemberLastName:  customer.LastName,
					MemberEmail:     customer.Email,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				})
			}()

			return nil
		}
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (c CustomerRepoImpl) FindAll(page string, newCustomerQuery bool) (*models.CustomerList, error) {
	conn := database.Sess

	//findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}

	if newCustomerQuery {
		//findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	err = conn.DB(database.DB).C(database.CUSTOMERS).Find(nil).Skip((pageNumber - 1) * perPage).Limit(perPage).All(&c.customers)

	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	key := config.Config("KEY")

	encrypt := helper.Encryption{Key: []byte(key)}

	decryptedCustomers := make([]models.Customer, 0, len(c.customers))

	for _, customer := range c.customers {
		wg.Add(5)
		go func() {
			defer wg.Done()
			pi, err := encrypt.Decrypt(customer.StreetAddress)

			if err != nil {
				panic(err)
			}

			customer.StreetAddress = pi
		}()

		go func() {
			defer wg.Done()

			if len(customer.OptionalAddress) > 0 {
				pi, err := encrypt.Decrypt(customer.OptionalAddress)

				if err != nil {
					panic(err)
				}

				customer.OptionalAddress = pi
			}
		}()

		go func() {
			defer wg.Done()
			pi, err := encrypt.Decrypt(customer.City)

			if err != nil {
				panic(err)
			}

			customer.City = pi
		}()

		go func() {
			defer wg.Done()
			pi, err := encrypt.Decrypt(customer.State)

			if err != nil {
				panic(err)
			}

			customer.State = pi
		}()

		go func() {
			defer wg.Done()
			pi, err := encrypt.Decrypt(customer.ZipCode)

			if err != nil {
				panic(err)
			}

			customer.ZipCode = pi
		}()
		wg.Wait()
		decryptedCustomers = append(decryptedCustomers, customer)
	}

	count, err := conn.DB(database.DB).C(database.CUSTOMERS).Count()

	if err != nil {
		panic(err)
	}

	c.customerList.NumberOfCustomers = count

	if c.customerList.NumberOfCustomers < 10 {
		c.customerList.NumberOfPages = 1
	} else {
		c.customerList.NumberOfPages = int(count/10) + 1
	}

	c.customerList.Customers = &decryptedCustomers
	c.customerList.CurrentPage = pageNumber

	return &c.customerList, nil
}

func (c CustomerRepoImpl) FindAllByFullName(firstName string, lastName string, page string, newLoginQuery bool) (*models.CustomerList, error) {
	conn := database.Sess

	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}

	if newLoginQuery {
		//findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	err = conn.DB(database.DB).C(database.CUSTOMERS).Find(bson.M{"firstName": firstName, "lastName": lastName}).Skip((pageNumber - 1) * perPage).Limit(perPage).All(&c.customers)

	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	key := config.Config("KEY")

	encrypt := helper.Encryption{Key: []byte(key)}

	decryptedCustomers := make([]models.Customer, 0, len(c.customers))

	for _, customer := range c.customers {
		wg.Add(5)
		go func() {
			defer wg.Done()
			pi, err := encrypt.Decrypt(customer.StreetAddress)

			if err != nil {
				panic(err)
			}

			customer.StreetAddress = pi
		}()

		go func() {
			defer wg.Done()

			if len(customer.OptionalAddress) > 0 {
				pi, err := encrypt.Decrypt(customer.OptionalAddress)

				if err != nil {
					panic(err)
				}

				customer.OptionalAddress = pi
			}
		}()

		go func() {
			defer wg.Done()
			pi, err := encrypt.Decrypt(customer.City)

			if err != nil {
				panic(err)
			}

			customer.City = pi
		}()

		go func() {
			defer wg.Done()
			pi, err := encrypt.Decrypt(customer.State)

			if err != nil {
				panic(err)
			}

			customer.State = pi
		}()

		go func() {
			defer wg.Done()
			pi, err := encrypt.Decrypt(customer.ZipCode)

			if err != nil {
				panic(err)
			}

			customer.ZipCode = pi
		}()
		wg.Wait()
		decryptedCustomers = append(decryptedCustomers, customer)
	}

	count, err := conn.DB(database.DB).C(database.CUSTOMERS).Count()

	if err != nil {
		panic(err)
	}

	c.customerList.NumberOfCustomers = count

	if c.customerList.NumberOfCustomers < 10 {
		c.customerList.NumberOfPages = 1
	} else {
		c.customerList.NumberOfPages = int(count/10) + 1
	}

	c.customerList.Customers = &decryptedCustomers
	c.customerList.CurrentPage = pageNumber

	return &c.customerList, nil
}

func (c CustomerRepoImpl) FindByEmail(email string) (*models.Customer, error) {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.CUSTOMERS).Find(bson.M{"email": email}).One(&c.customer)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err.Error() == "not found" {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	return &c.customer, nil
}

func (c CustomerRepoImpl) FindAllByOptInStatus(optIn bool) (*[]models.Customer, error) {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.CUSTOMERS).Find(bson.M{"infoEmailOptIn": optIn}).All(&c.customers)

	if err != nil {
		return nil, err
	}

	return &c.customers, nil
}

func (c CustomerRepoImpl) UpdateOptInStatus(status bool, email string) (*models.Customer, error) {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.CUSTOMERS).Update(bson.M{"email": email}, bson.M{"infoEmailOptIn": status,
		"updatedAt": time.Now()})

	if err != nil {
		return nil, errors.New("cannot update opt in status")
	}

	c.customer.InfoEmailOptIn = status

	return &c.customer, nil
}

func NewCustomerRepoImpl() CustomerRepoImpl {
	var customerRepoImpl CustomerRepoImpl

	return customerRepoImpl
}
