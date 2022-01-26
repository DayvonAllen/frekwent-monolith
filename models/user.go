package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type User struct {
	Id        bson.ObjectId `bson:"_id" json:"id"`
	Username  string        `bson:"username" json:"username"`
	Email     string        `bson:"email" json:"email"`
	Password  string        `bson:"password" json:"-"`
	CreatedAt time.Time     `bson:"createdAt" json:"-"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"-"`
}
