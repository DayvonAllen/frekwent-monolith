package tests

import (
	"freq/database"
	"freq/models"
	"freq/router"
	bson2 "github.com/globalsign/mgo/bson"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/goleak"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

var resource *dockertest.Resource
var pool *dockertest.Pool

func TestProductHandler_FindAll(t *testing.T) {
	defer goleak.VerifyNone(t)

	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	pool = p

	opts := dockertest.RunOptions{
		Repository:   "mongo",
		Tag:          "latest",
		ExposedPorts: []string{"27017"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"27017": {
				{HostIP: "0.0.0.0", HostPort: "27018"},
			},
		},
	}

	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	os.Setenv("DB_URL", "localhost:27018")

	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
		err := pool.Purge(r)
		if err != nil {

		}
	}(pool, resource)

	conn := database.Sess

	monId := bson2.ObjectIdHex("61caa9598eaeae5425c9780f")

	prod := &models.Product{
		Id: monId,
	}

	err = conn.DB(database.DB).C(database.PRODUCTS).Insert(prod)

	tests := []struct {
		description string

		// Test input
		route string

		// Expected output
		expectedError bool
		expectedCode  int
		expectedBody  bool
		evaluator     string
		query         string
		queryValue    string
	}{
		{
			description:   "get all products route",
			route:         "/products",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  true,
			evaluator:     "61caa9598eaeae5425c9780f",
			query:         "",
			queryValue:    "",
		},
		{
			description:   "get all products trending products error",
			route:         "/products",
			expectedError: false,
			expectedCode:  400,
			expectedBody:  true,
			evaluator:     "error",
			query:         "trending",
			queryValue:    "hello",
		},
		{
			description:   "get all products new products error",
			route:         "/products",
			expectedError: false,
			expectedCode:  400,
			expectedBody:  true,
			evaluator:     "error",
			query:         "new",
			queryValue:    "hello",
		},
		{
			description:   "no products",
			route:         "/products",
			expectedError: false,
			expectedCode:  500,
			expectedBody:  true,
			evaluator:     "no products",
			query:         "",
			queryValue:    "",
		},
	}

	// Setup the app as it is done in the main function
	app := router.Setup()

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route
		// from the test case
		req, _ := http.NewRequest(
			"GET",
			test.route,
			nil,
		)

		q := req.URL.Query()
		q.Add(test.query, test.queryValue)
		req.URL.RawQuery = q.Encode()

		if test.description == "no products" {
			_, err2 := conn.DB(database.DB).C(database.PRODUCTS).RemoveAll(bson.M{})
			if err2 != nil {
				return
			}
		}

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		// Verify, that the response body equals the expected body
		assert.Equalf(t, test.expectedBody, strings.Contains(string(body), test.evaluator), test.description)
	}
}

//func TestProductHandler_Create(t *testing.T) {
//	p, err := dockertest.NewPool("")
//	if err != nil {
//		log.Fatalf("could not connect to docker: %s", err)
//	}
//
//	pool = p
//
//	opts := dockertest.RunOptions{
//		Repository:   "mongo",
//		Tag:          "latest",
//		ExposedPorts: []string{"27017"},
//		PortBindings: map[docker.Port][]docker.PortBinding{
//			"27017": {
//				{HostIP: "0.0.0.0", HostPort: "27017"},
//			},
//		},
//	}
//
//	resource, err = pool.RunWithOptions(&opts)
//	if err != nil {
//		_ = pool.Purge(resource)
//		log.Fatalf("could not start resource: %s", err)
//	}
//
//	os.Setenv("DB_URL", "mongodb://localhost:27017")
//	os.Setenv("SECRET", "test")
//	os.Setenv("EXPIRATION", "120000")
//
//	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
//		err := pool.Purge(r)
//		if err != nil {
//
//		}
//	}(pool, resource)
//
//	conn := database.ConnectToDB()
//
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
//
//	conn.AdminCollection.InsertOne(context.TODO(), &models.User{Email: "hdoe@gmail.com", Username: "hdoe", Password: string(hashedPassword),
//		Id: primitive.NewObjectID()})
//
//	conn.ProductCollection.InsertOne(context.TODO(), &models.Product{
//		Ingredients:    []string{""},
//		Images:         []string{""},
//		Id:             primitive.NewObjectID(),
//		Quantity:       10,
//		Description:    "TEST",
//		Price:          "12.20",
//		Category:       "faceWash",
//		TimesPurchased: 0,
//		Name:           "test",
//	})
//
//	b, _ := json.Marshal(&models.Product{
//		Ingredients:    []string{""},
//		Images:         []string{""},
//		Id:             primitive.NewObjectID(),
//		Quantity:       10,
//		Description:    "TEST",
//		Price:          "12.20",
//		Category:       "faceWash",
//		TimesPurchased: 0,
//		Name:           "new",
//	})
//
//	tests := []struct {
//		description string
//
//		// Test input
//		route string
//
//		// Expected output
//		expectedError bool
//		expectedCode  int
//		expectedBody  bool
//		evaluator     string
//		query         string
//		queryValue    string
//		body          []byte
//		product       []byte
//	}{
//		{
//			description:   "unauthorized",
//			route:         "/iriguchi/items",
//			expectedError: false,
//			expectedCode:  401,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          nil,
//			product:       []byte(`{"name": "test product5",  "price": "10.01", "description": "desc...", "category": "faceWash"}`),
//		},
//		{
//			description:   "logged in - create products",
//			route:         "/iriguchi/items",
//			expectedError: false,
//			expectedCode:  201,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			product:       b,
//		},
//		{
//			description:   "logged in - create products - validation error for price",
//			route:         "/iriguchi/items",
//			expectedError: false,
//			expectedCode:  400,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			product:       []byte(`{"name": "test product5",  "price": "10", "description": "desc...", "category": "faceWash"}`),
//		},
//		{
//			description:   "logged in - create products - validation error for category",
//			route:         "/iriguchi/items",
//			expectedError: false,
//			expectedCode:  400,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			product:       []byte(`{"name": "test product5",  "price": "10.00", "quantity": 0, "description": "desc...", "category": ""}`),
//		},
//		{
//			description:   "logged in - create products - validation error for description",
//			route:         "/iriguchi/items",
//			expectedError: false,
//			expectedCode:  400,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			product:       []byte(`{"name": "test product5",  "price": "10.00", "quantity": 0, "description": "", "category": "faceWash"}`),
//		},
//		{
//			description:   "logged in - create products - validation error for name",
//			route:         "/iriguchi/items",
//			expectedError: false,
//			expectedCode:  400,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			product:       []byte(`{"name": "",  "price": "10.00", "quantity": 0, "description": "", "category": "faceWash"}`),
//		},
//		{
//			description:   "logged in - create products - can't create two products with the same name",
//			route:         "/iriguchi/items",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			product:       []byte(`{"name": "test",  "price": "10.00", "quantity": 0, "description": "ttt", "category": "faceWash"}`),
//		},
//		{
//			description:   "logged in - create products - incorrect JSON",
//			route:         "/iriguchi/items",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "error",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			product:       []byte(`{"name": "test product5",  "price": "10.00", "quantity": "", "description": "desc...", "category": "faceWash"}`),
//		},
//	}
//
//	// Setup the app as it is done in the main function
//	app := router.Setup()
//
//	// Iterate through test single test cases
//	for _, test := range tests {
//		// Create a new http request with the route
//		// from the test case
//		req, _ := http.NewRequest(
//			"POST",
//			test.route,
//			bytes.NewBuffer(test.product),
//		)
//
//		if strings.Contains(test.description, "logged in") {
//			re, _ := http.NewRequest(
//				"POST",
//				"/iriguchi/auth/login",
//				bytes.NewBuffer(test.body),
//			)
//
//			re.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//			myR, _ := app.Test(re, -1)
//
//			c := myR.Cookies()
//
//			myc := c[0]
//			req.AddCookie(myc)
//		}
//
//		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//		// Perform the request plain with the app.
//		// The -1 disables request latency.
//		res, err := app.Test(req, -1)
//
//		// verify that no error occured, that is not expected
//		assert.Equalf(t, test.expectedError, err != nil, test.description)
//
//		// As expected errors lead to broken responses, the next
//		// test case needs to be processed
//		if test.expectedError {
//			continue
//		}
//
//		// Verify if the status code is as expected
//		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
//
//		// Read the response body
//		body, err := ioutil.ReadAll(res.Body)
//
//		// Reading the response body should work everytime, such that
//		// the err variable should be nil
//		assert.Nilf(t, err, test.description)
//
//		// Verify, that the response body equals the expected body
//		assert.Equalf(t, test.expectedBody, strings.Contains(string(body), test.evaluator), test.description)
//	}
//}
//
//func TestProductHandler_FindAllByCategory(t *testing.T) {
//	p, err := dockertest.NewPool("")
//	if err != nil {
//		log.Fatalf("could not connect to docker: %s", err)
//	}
//
//	pool = p
//
//	opts := dockertest.RunOptions{
//		Repository:   "mongo",
//		Tag:          "latest",
//		ExposedPorts: []string{"27017"},
//		PortBindings: map[docker.Port][]docker.PortBinding{
//			"27017": {
//				{HostIP: "0.0.0.0", HostPort: "27017"},
//			},
//		},
//	}
//
//	resource, err = pool.RunWithOptions(&opts)
//	if err != nil {
//		_ = pool.Purge(resource)
//		log.Fatalf("could not start resource: %s", err)
//	}
//
//	os.Setenv("DB_URL", "mongodb://localhost:27017")
//	os.Setenv("SECRET", "test")
//	os.Setenv("EXPIRATION", "120000")
//
//	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
//		err := pool.Purge(r)
//		if err != nil {
//
//		}
//	}(pool, resource)
//
//	conn := database.ConnectToDB()
//
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
//
//	conn.AdminCollection.InsertOne(context.TODO(), &models.User{Email: "hdoe@gmail.com", Username: "hdoe", Password: string(hashedPassword),
//		Id: primitive.NewObjectID()})
//
//	conn.ProductCollection.InsertOne(context.TODO(), &models.Product{
//		Ingredients:    []string{""},
//		Images:         []string{""},
//		Id:             primitive.NewObjectID(),
//		Quantity:       10,
//		Description:    "TEST",
//		Price:          "12.20",
//		Category:       "facewash",
//		TimesPurchased: 0,
//		Name:           "test",
//	})
//
//	tests := []struct {
//		description string
//
//		// Test input
//		route string
//
//		// Expected output
//		expectedError bool
//		expectedCode  int
//		expectedBody  bool
//		evaluator     string
//		query         string
//		queryValue    string
//		body          []byte
//	}{
//		{
//			description:   "get product successfully",
//			route:         "/products/category/facewash",
//			expectedError: false,
//			expectedCode:  200,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          nil,
//		},
//		{
//			description:   "no product found",
//			route:         "/products/category/test",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          nil,
//		},
//		{
//			description:   "new query param validation error",
//			route:         "/products/category/facewash?new=llslsl",
//			expectedError: false,
//			expectedCode:  400,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          nil,
//		},
//	}
//
//	// Setup the app as it is done in the main function
//	app := router.Setup()
//
//	// Iterate through test single test cases
//	for _, test := range tests {
//		// Create a new http request with the route
//		// from the test case
//		req, _ := http.NewRequest(
//			"GET",
//			test.route,
//			nil,
//		)
//
//		// Perform the request plain with the app.
//		// The -1 disables request latency.
//		res, err := app.Test(req, -1)
//
//		// verify that no error occured, that is not expected
//		assert.Equalf(t, test.expectedError, err != nil, test.description)
//
//		// As expected errors lead to broken responses, the next
//		// test case needs to be processed
//		if test.expectedError {
//			continue
//		}
//
//		// Verify if the status code is as expected
//		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
//
//		// Read the response body
//		body, err := ioutil.ReadAll(res.Body)
//
//		// Reading the response body should work everytime, such that
//		// the err variable should be nil
//		assert.Nilf(t, err, test.description)
//
//		// Verify, that the response body equals the expected body
//		assert.Equalf(t, test.expectedBody, strings.Contains(string(body), test.evaluator), test.description)
//	}
//}
//
//func TestProductHandler_FindProductById(t *testing.T) {
//	p, err := dockertest.NewPool("")
//	if err != nil {
//		log.Fatalf("could not connect to docker: %s", err)
//	}
//
//	pool = p
//
//	opts := dockertest.RunOptions{
//		Repository:   "mongo",
//		Tag:          "latest",
//		ExposedPorts: []string{"27017"},
//		PortBindings: map[docker.Port][]docker.PortBinding{
//			"27017": {
//				{HostIP: "0.0.0.0", HostPort: "27017"},
//			},
//		},
//	}
//
//	resource, err = pool.RunWithOptions(&opts)
//	if err != nil {
//		_ = pool.Purge(resource)
//		log.Fatalf("could not start resource: %s", err)
//	}
//
//	os.Setenv("DB_URL", "mongodb://localhost:27017")
//	os.Setenv("SECRET", "test")
//	os.Setenv("EXPIRATION", "120000")
//
//	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
//		err := pool.Purge(r)
//		if err != nil {
//
//		}
//	}(pool, resource)
//
//	conn := database.ConnectToDB()
//
//	monId, _ := primitive.ObjectIDFromHex("61c32a0a61d12a5f03b73fc7")
//
//	conn.ProductCollection.InsertOne(context.TODO(), &models.Product{
//		Ingredients:    []string{""},
//		Images:         []string{""},
//		Id:             monId,
//		Quantity:       10,
//		Description:    "TEST",
//		Price:          "12.20",
//		Category:       "facewash",
//		TimesPurchased: 0,
//		Name:           "test",
//	})
//
//	tests := []struct {
//		description string
//
//		// Test input
//		route string
//
//		// Expected output
//		expectedError bool
//		expectedCode  int
//		expectedBody  bool
//		evaluator     string
//		query         string
//		queryValue    string
//		body          []byte
//	}{
//		{
//			description:   "get product by ID successfully",
//			route:         "/products/61c32a0a61d12a5f03b73fc7",
//			expectedError: false,
//			expectedCode:  200,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          nil,
//		},
//		{
//			description:   "product not in DB",
//			route:         "/products/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          nil,
//		},
//		{
//			description:   "invalid ID",
//			route:         "/products/1",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          nil,
//		},
//	}
//
//	// Setup the app as it is done in the main function
//	app := router.Setup()
//
//	// Iterate through test single test cases
//	for _, test := range tests {
//		// Create a new http request with the route
//		// from the test case
//		req, _ := http.NewRequest(
//			"GET",
//			test.route,
//			nil,
//		)
//
//		// Perform the request plain with the app.
//		// The -1 disables request latency.
//		res, err := app.Test(req, -1)
//
//		// verify that no error occured, that is not expected
//		assert.Equalf(t, test.expectedError, err != nil, test.description)
//
//		// As expected errors lead to broken responses, the next
//		// test case needs to be processed
//		if test.expectedError {
//			continue
//		}
//
//		// Verify if the status code is as expected
//		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
//
//		// Read the response body
//		body, err := ioutil.ReadAll(res.Body)
//
//		// Reading the response body should work everytime, such that
//		// the err variable should be nil
//		assert.Nilf(t, err, test.description)
//
//		// Verify, that the response body equals the expected body
//		assert.Equalf(t, test.expectedBody, strings.Contains(string(body), test.evaluator), test.description)
//	}
//}
//
//func TestProductHandler_FindProductByName(t *testing.T) {
//	p, err := dockertest.NewPool("")
//	if err != nil {
//		log.Fatalf("could not connect to docker: %s", err)
//	}
//
//	pool = p
//
//	opts := dockertest.RunOptions{
//		Repository:   "mongo",
//		Tag:          "latest",
//		ExposedPorts: []string{"27017"},
//		PortBindings: map[docker.Port][]docker.PortBinding{
//			"27017": {
//				{HostIP: "0.0.0.0", HostPort: "27017"},
//			},
//		},
//	}
//
//	resource, err = pool.RunWithOptions(&opts)
//	if err != nil {
//		_ = pool.Purge(resource)
//		log.Fatalf("could not start resource: %s", err)
//	}
//
//	os.Setenv("DB_URL", "mongodb://localhost:27017")
//	os.Setenv("SECRET", "test")
//	os.Setenv("EXPIRATION", "120000")
//
//	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
//		err := pool.Purge(r)
//		if err != nil {
//
//		}
//	}(pool, resource)
//
//	conn := database.ConnectToDB()
//
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
//
//	conn.AdminCollection.InsertOne(context.TODO(), &models.User{Email: "hdoe@gmail.com", Username: "hdoe", Password: string(hashedPassword),
//		Id: primitive.NewObjectID()})
//
//	monId, _ := primitive.ObjectIDFromHex("61c32a0a61d12a5f03b73fc7")
//
//	conn.ProductCollection.InsertOne(context.TODO(), &models.Product{
//		Ingredients:    []string{""},
//		Images:         []string{""},
//		Id:             monId,
//		Quantity:       10,
//		Description:    "TEST",
//		Price:          "12.20",
//		Category:       "facewash",
//		TimesPurchased: 0,
//		Name:           "test",
//	})
//
//	tests := []struct {
//		description string
//
//		// Test input
//		route string
//
//		// Expected output
//		expectedError bool
//		expectedCode  int
//		expectedBody  bool
//		evaluator     string
//		query         string
//		queryValue    string
//		body          []byte
//	}{
//		{
//			description:   "logged in - get product by name successfully",
//			route:         "/iriguchi/items/name?name=test",
//			expectedError: false,
//			expectedCode:  200,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//		},
//		{
//			description:   "logged in - cannot find by name",
//			route:         "/iriguchi/items/name",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//		},
//	}
//
//	// Setup the app as it is done in the main function
//	app := router.Setup()
//
//	// Iterate through test single test cases
//	for _, test := range tests {
//		// Create a new http request with the route
//		// from the test case
//		req, _ := http.NewRequest(
//			"GET",
//			test.route,
//			nil,
//		)
//
//		if strings.Contains(test.description, "logged in") {
//			re, _ := http.NewRequest(
//				"POST",
//				"/iriguchi/auth/login",
//				bytes.NewBuffer(test.body),
//			)
//
//			re.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//			myR, _ := app.Test(re, -1)
//
//			c := myR.Cookies()
//
//			myc := c[0]
//
//			req.AddCookie(myc)
//		}
//
//		// Perform the request plain with the app.
//		// The -1 disables request latency.
//		res, err := app.Test(req, -1)
//
//		// verify that no error occured, that is not expected
//		assert.Equalf(t, test.expectedError, err != nil, test.description)
//
//		// As expected errors lead to broken responses, the next
//		// test case needs to be processed
//		if test.expectedError {
//			continue
//		}
//
//		// Verify if the status code is as expected
//		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
//
//		// Read the response body
//		body, err := ioutil.ReadAll(res.Body)
//
//		// Reading the response body should work everytime, such that
//		// the err variable should be nil
//		assert.Nilf(t, err, test.description)
//
//		// Verify, that the response body equals the expected body
//		assert.Equalf(t, test.expectedBody, strings.Contains(string(body), test.evaluator), test.description)
//	}
//}
//
//func TestProductHandler_UpdateName(t *testing.T) {
//	p, err := dockertest.NewPool("")
//	if err != nil {
//		log.Fatalf("could not connect to docker: %s", err)
//	}
//
//	pool = p
//
//	opts := dockertest.RunOptions{
//		Repository:   "mongo",
//		Tag:          "latest",
//		ExposedPorts: []string{"27017"},
//		PortBindings: map[docker.Port][]docker.PortBinding{
//			"27017": {
//				{HostIP: "0.0.0.0", HostPort: "27017"},
//			},
//		},
//	}
//
//	resource, err = pool.RunWithOptions(&opts)
//	if err != nil {
//		_ = pool.Purge(resource)
//		log.Fatalf("could not start resource: %s", err)
//	}
//
//	os.Setenv("DB_URL", "mongodb://localhost:27017")
//	os.Setenv("SECRET", "test")
//	os.Setenv("EXPIRATION", "120000")
//
//	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
//		err := pool.Purge(r)
//		if err != nil {
//
//		}
//	}(pool, resource)
//
//	conn := database.ConnectToDB()
//
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
//
//	conn.AdminCollection.InsertOne(context.TODO(), &models.User{Email: "hdoe@gmail.com", Username: "hdoe", Password: string(hashedPassword),
//		Id: primitive.NewObjectID()})
//
//	monId, _ := primitive.ObjectIDFromHex("61c32a0a61d12a5f03b73fc7")
//
//	conn.ProductCollection.InsertOne(context.TODO(), &models.Product{
//		Ingredients:    []string{""},
//		Images:         []string{""},
//		Id:             monId,
//		Quantity:       10,
//		Description:    "TEST",
//		Price:          "12.20",
//		Category:       "facewash",
//		TimesPurchased: 0,
//		Name:           "test",
//	})
//
//	tests := []struct {
//		description string
//
//		// Test input
//		route string
//
//		// Expected output
//		expectedError bool
//		expectedCode  int
//		expectedBody  bool
//		evaluator     string
//		query         string
//		queryValue    string
//		body          []byte
//		name          []byte
//	}{
//		{
//			description:   "logged in - update name successfully",
//			route:         "/iriguchi/items/name/61c32a0a61d12a5f03b73fc7",
//			expectedError: false,
//			expectedCode:  200,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"name": "new"}`),
//		},
//		{
//			description:   "logged in - fail to update name",
//			route:         "/iriguchi/items/name/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"name": "new"}`),
//		},
//		{
//			description:   "logged in - validation error",
//			route:         "/iriguchi/items/name/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  400,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"name": ""}`),
//		},
//		{
//			description:   "logged in - invalid id",
//			route:         "/iriguchi/items/name/dl",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"name": "test"}`),
//		},
//		{
//			description:   "logged in - invalid json",
//			route:         "/iriguchi/items/name/dl",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"name": ld}`),
//		},
//	}
//
//	// Setup the app as it is done in the main function
//	app := router.Setup()
//
//	// Iterate through test single test cases
//	for _, test := range tests {
//		// Create a new http request with the route
//		// from the test case
//		req, _ := http.NewRequest(
//			"PUT",
//			test.route,
//			bytes.NewBuffer(test.name),
//		)
//
//		if strings.Contains(test.description, "logged in") {
//			re, _ := http.NewRequest(
//				"POST",
//				"/iriguchi/auth/login",
//				bytes.NewBuffer(test.body),
//			)
//
//			re.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//			myR, _ := app.Test(re, -1)
//
//			c := myR.Cookies()
//
//			myc := c[0]
//
//			req.AddCookie(myc)
//		}
//
//		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//		// Perform the request plain with the app.
//		// The -1 disables request latency.
//		res, err := app.Test(req, -1)
//
//		// verify that no error occured, that is not expected
//		assert.Equalf(t, test.expectedError, err != nil, test.description)
//
//		// As expected errors lead to broken responses, the next
//		// test case needs to be processed
//		if test.expectedError {
//			continue
//		}
//
//		// Verify if the status code is as expected
//		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
//
//		// Read the response body
//		body, err := ioutil.ReadAll(res.Body)
//
//		// Reading the response body should work everytime, such that
//		// the err variable should be nil
//		assert.Nilf(t, err, test.description)
//
//		// Verify, that the response body equals the expected body
//		assert.Equalf(t, test.expectedBody, strings.Contains(string(body), test.evaluator), test.description)
//	}
//}
//
//func TestProductHandler_UpdatePrice(t *testing.T) {
//	p, err := dockertest.NewPool("")
//	if err != nil {
//		log.Fatalf("could not connect to docker: %s", err)
//	}
//
//	pool = p
//
//	opts := dockertest.RunOptions{
//		Repository:   "mongo",
//		Tag:          "latest",
//		ExposedPorts: []string{"27017"},
//		PortBindings: map[docker.Port][]docker.PortBinding{
//			"27017": {
//				{HostIP: "0.0.0.0", HostPort: "27017"},
//			},
//		},
//	}
//
//	resource, err = pool.RunWithOptions(&opts)
//	if err != nil {
//		_ = pool.Purge(resource)
//		log.Fatalf("could not start resource: %s", err)
//	}
//
//	os.Setenv("DB_URL", "mongodb://localhost:27017")
//	os.Setenv("SECRET", "test")
//	os.Setenv("EXPIRATION", "120000")
//
//	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
//		err := pool.Purge(r)
//		if err != nil {
//
//		}
//	}(pool, resource)
//
//	conn := database.ConnectToDB()
//
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
//
//	conn.AdminCollection.InsertOne(context.TODO(), &models.User{Email: "hdoe@gmail.com", Username: "hdoe", Password: string(hashedPassword),
//		Id: primitive.NewObjectID()})
//
//	monId, _ := primitive.ObjectIDFromHex("61c32a0a61d12a5f03b73fc7")
//
//	conn.ProductCollection.InsertOne(context.TODO(), &models.Product{
//		Ingredients:    []string{""},
//		Images:         []string{""},
//		Id:             monId,
//		Quantity:       10,
//		Description:    "TEST",
//		Price:          "12.20",
//		Category:       "facewash",
//		TimesPurchased: 0,
//		Name:           "test",
//	})
//
//	tests := []struct {
//		description string
//
//		// Test input
//		route string
//
//		// Expected output
//		expectedError bool
//		expectedCode  int
//		expectedBody  bool
//		evaluator     string
//		query         string
//		queryValue    string
//		body          []byte
//		name          []byte
//	}{
//		{
//			description:   "logged in - update price successfully",
//			route:         "/iriguchi/items/price/61c32a0a61d12a5f03b73fc7",
//			expectedError: false,
//			expectedCode:  200,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"price": "12.20"}`),
//		},
//		{
//			description:   "logged in - fail to update price",
//			route:         "/iriguchi/items/price/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"price": "120.20"}`),
//		},
//		{
//			description:   "logged in - validation error price",
//			route:         "/iriguchi/items/price/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  400,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"price": "11"}`),
//		},
//		{
//			description:   "logged in - invalid id price",
//			route:         "/iriguchi/items/price/dl",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"price": "11.22"}`),
//		},
//		{
//			description:   "logged in - invalid json price",
//			route:         "/iriguchi/items/price/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"price": ld}`),
//		},
//	}
//
//	// Setup the app as it is done in the main function
//	app := router.Setup()
//
//	// Iterate through test single test cases
//	for _, test := range tests {
//		// Create a new http request with the route
//		// from the test case
//		req, _ := http.NewRequest(
//			"PUT",
//			test.route,
//			bytes.NewBuffer(test.name),
//		)
//
//		if strings.Contains(test.description, "logged in") {
//			re, _ := http.NewRequest(
//				"POST",
//				"/iriguchi/auth/login",
//				bytes.NewBuffer(test.body),
//			)
//
//			re.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//			myR, _ := app.Test(re, -1)
//
//			c := myR.Cookies()
//
//			myc := c[0]
//
//			req.AddCookie(myc)
//		}
//
//		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//		// Perform the request plain with the app.
//		// The -1 disables request latency.
//		res, err := app.Test(req, -1)
//
//		// verify that no error occured, that is not expected
//		assert.Equalf(t, test.expectedError, err != nil, test.description)
//
//		// As expected errors lead to broken responses, the next
//		// test case needs to be processed
//		if test.expectedError {
//			continue
//		}
//
//		// Verify if the status code is as expected
//		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
//
//		// Read the response body
//		body, err := ioutil.ReadAll(res.Body)
//
//		// Reading the response body should work everytime, such that
//		// the err variable should be nil
//		assert.Nilf(t, err, test.description)
//
//		// Verify, that the response body equals the expected body
//		assert.Equalf(t, test.expectedBody, strings.Contains(string(body), test.evaluator), test.description)
//	}
//}
//
//func TestProductHandler_UpdateDescription(t *testing.T) {
//	p, err := dockertest.NewPool("")
//	if err != nil {
//		log.Fatalf("could not connect to docker: %s", err)
//	}
//
//	pool = p
//
//	opts := dockertest.RunOptions{
//		Repository:   "mongo",
//		Tag:          "latest",
//		ExposedPorts: []string{"27017"},
//		PortBindings: map[docker.Port][]docker.PortBinding{
//			"27017": {
//				{HostIP: "0.0.0.0", HostPort: "27017"},
//			},
//		},
//	}
//
//	resource, err = pool.RunWithOptions(&opts)
//	if err != nil {
//		_ = pool.Purge(resource)
//		log.Fatalf("could not start resource: %s", err)
//	}
//
//	os.Setenv("DB_URL", "mongodb://localhost:27017")
//	os.Setenv("SECRET", "test")
//	os.Setenv("EXPIRATION", "120000")
//
//	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
//		err := pool.Purge(r)
//		if err != nil {
//
//		}
//	}(pool, resource)
//
//	conn := database.ConnectToDB()
//
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
//
//	conn.AdminCollection.InsertOne(context.TODO(), &models.User{Email: "hdoe@gmail.com", Username: "hdoe", Password: string(hashedPassword),
//		Id: primitive.NewObjectID()})
//
//	monId, _ := primitive.ObjectIDFromHex("61c32a0a61d12a5f03b73fc7")
//
//	conn.ProductCollection.InsertOne(context.TODO(), &models.Product{
//		Ingredients:    []string{""},
//		Images:         []string{""},
//		Id:             monId,
//		Quantity:       10,
//		Description:    "TEST",
//		Price:          "12.20",
//		Category:       "facewash",
//		TimesPurchased: 0,
//		Name:           "test",
//	})
//
//	tests := []struct {
//		description string
//
//		// Test input
//		route string
//
//		// Expected output
//		expectedError bool
//		expectedCode  int
//		expectedBody  bool
//		evaluator     string
//		query         string
//		queryValue    string
//		body          []byte
//		name          []byte
//	}{
//		{
//			description:   "logged in - update description successfully",
//			route:         "/iriguchi/items/description/61c32a0a61d12a5f03b73fc7",
//			expectedError: false,
//			expectedCode:  200,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"description": "new"}`),
//		},
//		{
//			description:   "logged in - fail to update description",
//			route:         "/iriguchi/items/description/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"description": "ld"}`),
//		},
//		{
//			description:   "logged in - validation error description",
//			route:         "/iriguchi/items/description/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  400,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"description": ""}`),
//		},
//		{
//			description:   "logged in - invalid id description",
//			route:         "/iriguchi/items/description/dl",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"description": "11.22"}`),
//		},
//		{
//			description:   "logged in - invalid json description",
//			route:         "/iriguchi/items/description/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"description": ld}`),
//		},
//	}
//
//	// Setup the app as it is done in the main function
//	app := router.Setup()
//
//	// Iterate through test single test cases
//	for _, test := range tests {
//		// Create a new http request with the route
//		// from the test case
//		req, _ := http.NewRequest(
//			"PUT",
//			test.route,
//			bytes.NewBuffer(test.name),
//		)
//
//		if strings.Contains(test.description, "logged in") {
//			re, _ := http.NewRequest(
//				"POST",
//				"/iriguchi/auth/login",
//				bytes.NewBuffer(test.body),
//			)
//
//			re.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//			myR, _ := app.Test(re, -1)
//
//			c := myR.Cookies()
//
//			myc := c[0]
//
//			req.AddCookie(myc)
//		}
//
//		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//		// Perform the request plain with the app.
//		// The -1 disables request latency.
//		res, err := app.Test(req, -1)
//
//		// verify that no error occured, that is not expected
//		assert.Equalf(t, test.expectedError, err != nil, test.description)
//
//		// As expected errors lead to broken responses, the next
//		// test case needs to be processed
//		if test.expectedError {
//			continue
//		}
//
//		// Verify if the status code is as expected
//		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
//
//		// Read the response body
//		body, err := ioutil.ReadAll(res.Body)
//
//		// Reading the response body should work everytime, such that
//		// the err variable should be nil
//		assert.Nilf(t, err, test.description)
//
//		// Verify, that the response body equals the expected body
//		assert.Equalf(t, test.expectedBody, strings.Contains(string(body), test.evaluator), test.description)
//	}
//}
//
//func TestProductHandler_UpdateQuantity(t *testing.T) {
//	p, err := dockertest.NewPool("")
//	if err != nil {
//		log.Fatalf("could not connect to docker: %s", err)
//	}
//
//	pool = p
//
//	opts := dockertest.RunOptions{
//		Repository:   "mongo",
//		Tag:          "latest",
//		ExposedPorts: []string{"27017"},
//		PortBindings: map[docker.Port][]docker.PortBinding{
//			"27017": {
//				{HostIP: "0.0.0.0", HostPort: "27017"},
//			},
//		},
//	}
//
//	resource, err = pool.RunWithOptions(&opts)
//	if err != nil {
//		_ = pool.Purge(resource)
//		log.Fatalf("could not start resource: %s", err)
//	}
//
//	os.Setenv("DB_URL", "mongodb://localhost:27017")
//	os.Setenv("SECRET", "test")
//	os.Setenv("EXPIRATION", "120000")
//
//	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
//		err := pool.Purge(r)
//		if err != nil {
//
//		}
//	}(pool, resource)
//
//	conn := database.ConnectToDB()
//
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
//
//	conn.AdminCollection.InsertOne(context.TODO(), &models.User{Email: "hdoe@gmail.com", Username: "hdoe", Password: string(hashedPassword),
//		Id: primitive.NewObjectID()})
//
//	monId, _ := primitive.ObjectIDFromHex("61c32a0a61d12a5f03b73fc7")
//
//	conn.ProductCollection.InsertOne(context.TODO(), &models.Product{
//		Ingredients:    []string{""},
//		Images:         []string{""},
//		Id:             monId,
//		Quantity:       10,
//		Description:    "TEST",
//		Price:          "12.20",
//		Category:       "facewash",
//		TimesPurchased: 0,
//		Name:           "test",
//	})
//
//	tests := []struct {
//		description string
//
//		// Test input
//		route string
//
//		// Expected output
//		expectedError bool
//		expectedCode  int
//		expectedBody  bool
//		evaluator     string
//		query         string
//		queryValue    string
//		body          []byte
//		name          []byte
//	}{
//		{
//			description:   "logged in - update quantity successfully",
//			route:         "/iriguchi/items/quantity/61c32a0a61d12a5f03b73fc7",
//			expectedError: false,
//			expectedCode:  200,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"quantity": 1}`),
//		},
//		{
//			description:   "logged in - fail to quantity",
//			route:         "/iriguchi/items/quantity/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"quantity": 1}`),
//		},
//		{
//			description:   "logged in - invalid id quantity",
//			route:         "/iriguchi/items/quantity/dl",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"quantity": 1}`),
//		},
//		{
//			description:   "logged in - invalid json quantity",
//			route:         "/iriguchi/items/quantity/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"quantity" 0}`),
//		},
//	}
//
//	// Setup the app as it is done in the main function
//	app := router.Setup()
//
//	// Iterate through test single test cases
//	for _, test := range tests {
//		// Create a new http request with the route
//		// from the test case
//		req, _ := http.NewRequest(
//			"PUT",
//			test.route,
//			bytes.NewBuffer(test.name),
//		)
//
//		if strings.Contains(test.description, "logged in") {
//			re, _ := http.NewRequest(
//				"POST",
//				"/iriguchi/auth/login",
//				bytes.NewBuffer(test.body),
//			)
//
//			re.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//			myR, _ := app.Test(re, -1)
//
//			c := myR.Cookies()
//
//			myc := c[0]
//
//			req.AddCookie(myc)
//		}
//
//		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//		// Perform the request plain with the app.
//		// The -1 disables request latency.
//		res, err := app.Test(req, -1)
//
//		// verify that no error occured, that is not expected
//		assert.Equalf(t, test.expectedError, err != nil, test.description)
//
//		// As expected errors lead to broken responses, the next
//		// test case needs to be processed
//		if test.expectedError {
//			continue
//		}
//
//		// Verify if the status code is as expected
//		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
//
//		// Read the response body
//		body, err := ioutil.ReadAll(res.Body)
//
//		// Reading the response body should work everytime, such that
//		// the err variable should be nil
//		assert.Nilf(t, err, test.description)
//
//		// Verify, that the response body equals the expected body
//		assert.Equalf(t, test.expectedBody, strings.Contains(string(body), test.evaluator), test.description)
//	}
//}
//
//func TestProductHandler_UpdateIngredients(t *testing.T) {
//	p, err := dockertest.NewPool("")
//	if err != nil {
//		log.Fatalf("could not connect to docker: %s", err)
//	}
//
//	pool = p
//
//	opts := dockertest.RunOptions{
//		Repository:   "mongo",
//		Tag:          "latest",
//		ExposedPorts: []string{"27017"},
//		PortBindings: map[docker.Port][]docker.PortBinding{
//			"27017": {
//				{HostIP: "0.0.0.0", HostPort: "27017"},
//			},
//		},
//	}
//
//	resource, err = pool.RunWithOptions(&opts)
//	if err != nil {
//		_ = pool.Purge(resource)
//		log.Fatalf("could not start resource: %s", err)
//	}
//
//	os.Setenv("DB_URL", "mongodb://localhost:27017")
//	os.Setenv("SECRET", "test")
//	os.Setenv("EXPIRATION", "120000")
//
//	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
//		err := pool.Purge(r)
//		if err != nil {
//
//		}
//	}(pool, resource)
//
//	conn := database.ConnectToDB()
//
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
//
//	conn.AdminCollection.InsertOne(context.TODO(), &models.User{Email: "hdoe@gmail.com", Username: "hdoe", Password: string(hashedPassword),
//		Id: primitive.NewObjectID()})
//
//	monId, _ := primitive.ObjectIDFromHex("61c32a0a61d12a5f03b73fc7")
//
//	conn.ProductCollection.InsertOne(context.TODO(), &models.Product{
//		Ingredients:    []string{""},
//		Images:         []string{""},
//		Id:             monId,
//		Quantity:       10,
//		Description:    "TEST",
//		Price:          "12.20",
//		Category:       "facewash",
//		TimesPurchased: 0,
//		Name:           "test",
//	})
//
//	tests := []struct {
//		description string
//
//		// Test input
//		route string
//
//		// Expected output
//		expectedError bool
//		expectedCode  int
//		expectedBody  bool
//		evaluator     string
//		query         string
//		queryValue    string
//		body          []byte
//		name          []byte
//	}{
//		{
//			description:   "logged in - update ingredients successfully",
//			route:         "/iriguchi/items/ingredients/61c32a0a61d12a5f03b73fc7",
//			expectedError: false,
//			expectedCode:  200,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"ingredients": ["hello"]}`),
//		},
//		{
//			description:   "logged in - fail to ingredients",
//			route:         "/iriguchi/items/ingredients/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"ingredients": []}`),
//		},
//		{
//			description:   "logged in - invalid id ingredients",
//			route:         "/iriguchi/items/ingredients/dl",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"ingredients": []}`),
//		},
//		{
//			description:   "logged in - invalid json ingredients",
//			route:         "/iriguchi/items/ingredients/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"ingredients" []}`),
//		},
//	}
//
//	// Setup the app as it is done in the main function
//	app := router.Setup()
//
//	// Iterate through test single test cases
//	for _, test := range tests {
//		// Create a new http request with the route
//		// from the test case
//		req, _ := http.NewRequest(
//			"PUT",
//			test.route,
//			bytes.NewBuffer(test.name),
//		)
//
//		if strings.Contains(test.description, "logged in") {
//			re, _ := http.NewRequest(
//				"POST",
//				"/iriguchi/auth/login",
//				bytes.NewBuffer(test.body),
//			)
//
//			re.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//			myR, _ := app.Test(re, -1)
//
//			c := myR.Cookies()
//
//			myc := c[0]
//
//			req.AddCookie(myc)
//		}
//
//		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//		// Perform the request plain with the app.
//		// The -1 disables request latency.
//		res, err := app.Test(req, -1)
//
//		// verify that no error occured, that is not expected
//		assert.Equalf(t, test.expectedError, err != nil, test.description)
//
//		// As expected errors lead to broken responses, the next
//		// test case needs to be processed
//		if test.expectedError {
//			continue
//		}
//
//		// Verify if the status code is as expected
//		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
//
//		// Read the response body
//		body, err := ioutil.ReadAll(res.Body)
//
//		// Reading the response body should work everytime, such that
//		// the err variable should be nil
//		assert.Nilf(t, err, test.description)
//
//		// Verify, that the response body equals the expected body
//		assert.Equalf(t, test.expectedBody, strings.Contains(string(body), test.evaluator), test.description)
//	}
//}
//
//func TestProductHandler_UpdateCategory(t *testing.T) {
//	p, err := dockertest.NewPool("")
//	if err != nil {
//		log.Fatalf("could not connect to docker: %s", err)
//	}
//
//	pool = p
//
//	opts := dockertest.RunOptions{
//		Repository:   "mongo",
//		Tag:          "latest",
//		ExposedPorts: []string{"27017"},
//		PortBindings: map[docker.Port][]docker.PortBinding{
//			"27017": {
//				{HostIP: "0.0.0.0", HostPort: "27017"},
//			},
//		},
//	}
//
//	resource, err = pool.RunWithOptions(&opts)
//	if err != nil {
//		_ = pool.Purge(resource)
//		log.Fatalf("could not start resource: %s", err)
//	}
//
//	os.Setenv("DB_URL", "mongodb://localhost:27017")
//	os.Setenv("SECRET", "test")
//	os.Setenv("EXPIRATION", "120000")
//
//	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
//		err := pool.Purge(r)
//		if err != nil {
//
//		}
//	}(pool, resource)
//
//	conn := database.ConnectToDB()
//
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
//
//	conn.AdminCollection.InsertOne(context.TODO(), &models.User{Email: "hdoe@gmail.com", Username: "hdoe", Password: string(hashedPassword),
//		Id: primitive.NewObjectID()})
//
//	monId, _ := primitive.ObjectIDFromHex("61c32a0a61d12a5f03b73fc7")
//
//	conn.ProductCollection.InsertOne(context.TODO(), &models.Product{
//		Ingredients:    []string{""},
//		Images:         []string{""},
//		Id:             monId,
//		Quantity:       10,
//		Description:    "TEST",
//		Price:          "12.20",
//		Category:       "facewash",
//		TimesPurchased: 0,
//		Name:           "test",
//	})
//
//	tests := []struct {
//		description string
//
//		// Test input
//		route string
//
//		// Expected output
//		expectedError bool
//		expectedCode  int
//		expectedBody  bool
//		evaluator     string
//		query         string
//		queryValue    string
//		body          []byte
//		name          []byte
//	}{
//		{
//			description:   "logged in - update category successfully",
//			route:         "/iriguchi/items/category/61c32a0a61d12a5f03b73fc7",
//			expectedError: false,
//			expectedCode:  200,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"category": "new"}`),
//		},
//		{
//			description:   "logged in - fail to category",
//			route:         "/iriguchi/items/category/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"category": "ld"}`),
//		},
//		{
//			description:   "logged in - validation error category",
//			route:         "/iriguchi/items/category/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  400,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"category": ""}`),
//		},
//		{
//			description:   "logged in - invalid id category",
//			route:         "/iriguchi/items/category/dl",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"category": "new"}`),
//		},
//		{
//			description:   "logged in - invalid json category",
//			route:         "/iriguchi/items/category/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//			name:          []byte(`{"category": ld}`),
//		},
//	}
//
//	// Setup the app as it is done in the main function
//	app := router.Setup()
//
//	// Iterate through test single test cases
//	for _, test := range tests {
//		// Create a new http request with the route
//		// from the test case
//		req, _ := http.NewRequest(
//			"PUT",
//			test.route,
//			bytes.NewBuffer(test.name),
//		)
//
//		if strings.Contains(test.description, "logged in") {
//			re, _ := http.NewRequest(
//				"POST",
//				"/iriguchi/auth/login",
//				bytes.NewBuffer(test.body),
//			)
//
//			re.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//			myR, _ := app.Test(re, -1)
//
//			c := myR.Cookies()
//
//			myc := c[0]
//
//			req.AddCookie(myc)
//		}
//
//		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//		// Perform the request plain with the app.
//		// The -1 disables request latency.
//		res, err := app.Test(req, -1)
//
//		// verify that no error occured, that is not expected
//		assert.Equalf(t, test.expectedError, err != nil, test.description)
//
//		// As expected errors lead to broken responses, the next
//		// test case needs to be processed
//		if test.expectedError {
//			continue
//		}
//
//		// Verify if the status code is as expected
//		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
//
//		// Read the response body
//		body, err := ioutil.ReadAll(res.Body)
//
//		// Reading the response body should work everytime, such that
//		// the err variable should be nil
//		assert.Nilf(t, err, test.description)
//
//		// Verify, that the response body equals the expected body
//		assert.Equalf(t, test.expectedBody, strings.Contains(string(body), test.evaluator), test.description)
//	}
//}
//
//func TestProductHandler_DeleteById(t *testing.T) {
//	p, err := dockertest.NewPool("")
//	if err != nil {
//		log.Fatalf("could not connect to docker: %s", err)
//	}
//
//	pool = p
//
//	opts := dockertest.RunOptions{
//		Repository:   "mongo",
//		Tag:          "latest",
//		ExposedPorts: []string{"27017"},
//		PortBindings: map[docker.Port][]docker.PortBinding{
//			"27017": {
//				{HostIP: "0.0.0.0", HostPort: "27017"},
//			},
//		},
//	}
//
//	resource, err = pool.RunWithOptions(&opts)
//	if err != nil {
//		_ = pool.Purge(resource)
//		log.Fatalf("could not start resource: %s", err)
//	}
//
//	os.Setenv("DB_URL", "mongodb://localhost:27017")
//	os.Setenv("SECRET", "test")
//	os.Setenv("EXPIRATION", "120000")
//
//	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
//		err := pool.Purge(r)
//		if err != nil {
//
//		}
//	}(pool, resource)
//
//	conn := database.ConnectToDB()
//
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
//
//	conn.AdminCollection.InsertOne(context.TODO(), &models.User{Email: "hdoe@gmail.com", Username: "hdoe", Password: string(hashedPassword),
//		Id: primitive.NewObjectID()})
//
//	monId, _ := primitive.ObjectIDFromHex("61c32a0a61d12a5f03b73fc7")
//
//	conn.ProductCollection.InsertOne(context.TODO(), &models.Product{
//		Ingredients:    []string{""},
//		Images:         []string{""},
//		Id:             monId,
//		Quantity:       10,
//		Description:    "TEST",
//		Price:          "12.20",
//		Category:       "facewash",
//		TimesPurchased: 0,
//		Name:           "test",
//	})
//
//	tests := []struct {
//		description string
//
//		// Test input
//		route string
//
//		// Expected output
//		expectedError bool
//		expectedCode  int
//		expectedBody  bool
//		evaluator     string
//		query         string
//		queryValue    string
//		body          []byte
//	}{
//		{
//			description:   "logged in - delete product successfully",
//			route:         "/iriguchi/items/delete/61c32a0a61d12a5f03b73fc7",
//			expectedError: false,
//			expectedCode:  204,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//		},
//		{
//			description:   "logged in - fail to delete",
//			route:         "/iriguchi/items/delete/61c32a0a61d12a5f03b73fc8",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//		},
//
//		{
//			description:   "logged in - invalid id delete",
//			route:         "/iriguchi/items/delete/dl",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//		},
//	}
//
//	// Setup the app as it is done in the main function
//	app := router.Setup()
//
//	// Iterate through test single test cases
//	for _, test := range tests {
//		// Create a new http request with the route
//		// from the test case
//		req, _ := http.NewRequest(
//			"DELETE",
//			test.route,
//			nil,
//		)
//
//		if strings.Contains(test.description, "logged in") {
//			re, _ := http.NewRequest(
//				"POST",
//				"/iriguchi/auth/login",
//				bytes.NewBuffer(test.body),
//			)
//
//			re.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//			myR, _ := app.Test(re, -1)
//
//			c := myR.Cookies()
//
//			myc := c[0]
//
//			req.AddCookie(myc)
//		}
//
//		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//		// Perform the request plain with the app.
//		// The -1 disables request latency.
//		res, err := app.Test(req, -1)
//
//		// verify that no error occured, that is not expected
//		assert.Equalf(t, test.expectedError, err != nil, test.description)
//
//		// As expected errors lead to broken responses, the next
//		// test case needs to be processed
//		if test.expectedError {
//			continue
//		}
//
//		// Verify if the status code is as expected
//		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
//
//		// Read the response body
//		body, err := ioutil.ReadAll(res.Body)
//
//		// Reading the response body should work everytime, such that
//		// the err variable should be nil
//		assert.Nilf(t, err, test.description)
//
//		// Verify, that the response body equals the expected body
//		assert.Equalf(t, test.expectedBody, strings.Contains(string(body), test.evaluator), test.description)
//	}
//}
