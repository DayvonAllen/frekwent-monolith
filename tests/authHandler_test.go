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
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"go.uber.org/goleak"
//	"golang.org/x/crypto/bcrypt"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"os"
//	"strings"
//	"testing"
//)
//
//func TestAuthHandler_Login(t *testing.T) {
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
//				{HostIP: "0.0.0.0", HostPort: "27018"},
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
//	os.Setenv("DB_URL", "mongodb://localhost:27018")
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
//			description:   "log in",
//			route:         "/iriguchi/auth/login",
//			expectedError: false,
//			expectedCode:  200,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com", "password": "password"}`),
//		},
//		{
//			description:   "body parsing error",
//			route:         "/iriguchi/auth/login",
//			expectedError: false,
//			expectedCode:  500,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "hdoe@gmail.com" "password": "password"}`),
//		},
//		{
//			description:   "email validation error",
//			route:         "/iriguchi/auth/login",
//			expectedError: false,
//			expectedCode:  400,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "", "password": "password"}`),
//		},
//		{
//			description:   "password validation error",
//			route:         "/iriguchi/auth/login",
//			expectedError: false,
//			expectedCode:  400,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "test@test.com", "password": ""}`),
//		},
//		{
//			description:   "wrong credentials error",
//			route:         "/iriguchi/auth/login",
//			expectedError: false,
//			expectedCode:  401,
//			expectedBody:  true,
//			evaluator:     "",
//			query:         "",
//			queryValue:    "",
//			body:          []byte(`{"email": "test@test.com", "password": "password"}`),
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
//			"/iriguchi/auth/login",
//			bytes.NewBuffer(test.body),
//		)
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
//func TestAuthHandler_Logout(t *testing.T) {
//	defer goleak.VerifyNone(t)
//
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
//				{HostIP: "0.0.0.0", HostPort: "27018"},
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
//	os.Setenv("DB_URL", "mongodb://localhost:27018")
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
//			description:   "logout",
//			route:         "/iriguchi/auth/logout",
//			expectedError: false,
//			expectedCode:  200,
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
//			"POST",
//			"/iriguchi/auth/login",
//			bytes.NewBuffer(test.body),
//		)
//
//		re, _ := http.NewRequest(
//			"GET",
//			"/iriguchi/auth/logout",
//			bytes.NewBuffer(test.body),
//		)
//
//		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//		myR, _ := app.Test(req, -1)
//
//		c := myR.Cookies()
//
//		myc := c[0]
//		re.AddCookie(myc)
//
//		// Perform the request plain with the app.
//		// The -1 disables request latency.
//		res, err := app.Test(re, -1)
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
