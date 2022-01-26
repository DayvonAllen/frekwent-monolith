package database

import (
	"freq/config"
	"github.com/globalsign/mgo"
	"log"
	"time"
)

const ADMIN = "admin"
const PRODUCTS = "products"
const EMAILS = "emails"
const COUPONS = "coupons"
const IPS = "loginIPs"
const CUSTOMERS = "customers"
const PURCHASES = "purchases"
const MAIL_MEMBERS = "mailMembers"
const DB = "frekwent"

var Sess = ConnectToDB()

func ConnectToDB() *mgo.Session {
	u := config.Config("DB_URL")

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:     []string{u},
		Timeout:   60 * time.Second,
		PoolLimit: 20,
		Database:  "Frekwent",
	}

	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	mongoSession.SetMode(mgo.Monotonic, true)

	return mongoSession
}
