package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type MailMember struct {
	Id              bson.ObjectId `bson:"_id" json:"id"`
	MemberFirstName string        `bson:"memberFirstName" json:"memberFirstName"`
	MemberLastName  string        `bson:"memberLastName" json:"memberLastName"`
	MemberEmail     string        `bson:"memberEmail" json:"memberEmail"`
	CreatedAt       time.Time     `bson:"createdAt" json:"-"`
	UpdatedAt       time.Time     `bson:"updatedAt" json:"updatedAt"`
}
