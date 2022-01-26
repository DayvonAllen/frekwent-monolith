package tests

//
//import (
//	"bytes"
//	"context"
//	"freq/database"
//	"freq/models"
//	"freq/router"
//	"github.com/ory/dockertest/v3"
//	"github.com/ory/dockertest/v3/docker"
//	"github.com/stretchr/testify/assert"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"golang.org/x/crypto/bcrypt"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"os"
//	"strings"
//	"testing"
//)
//
//func TestLoginIpHandler_FindAll(t *testing.T) {
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
//	conn.LoginIPCollection.InsertOne(context.TODO(), &models.LoginIP{
//		Id:        primitive.NewObjectID(),
//		IpAddress: "127.0.0.1",
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
//			description:   "logged in - get all IP",
//			route:         "/iriguchi/ip",
//			expectedError: false,
//			expectedCode:  200,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//		},
//		{
//			description:   "logged in - bad query param",
//			route:         "/iriguchi/ip",
//			expectedError: false,
//			expectedCode:  400,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "new",
//			queryValue:    "ddjdjdj",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//		},
//		{
//			description:   "logged in - no products",
//			route:         "/iriguchi/ip",
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
//			req.AddCookie(myc)
//		}
//
//		if strings.Contains(test.description, "no products") {
//			conn.LoginIPCollection.DeleteMany(context.TODO(), bson.M{})
//		}
//
//		q := req.URL.Query()
//		q.Add(test.query, test.queryValue)
//		req.URL.RawQuery = q.Encode()
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
//func TestLoginIpHandler_FindById(t *testing.T) {
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
//	conn.LoginIPCollection.InsertOne(context.TODO(), &models.LoginIP{
//		Id:        primitive.NewObjectID(),
//		IpAddress: "127.0.0.1",
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
//			description:   "logged in - get by IP",
//			route:         "/iriguchi/ip/get/127.0.0.1",
//			expectedError: false,
//			expectedCode:  200,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//		},
//		{
//			description:   "logged in - no IP found",
//			route:         "/iriguchi/ip/get/ldldld",
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
