package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type LoginIP struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	IpAddress   string        `bson:"ipAddress" json:"ipAddress"`
	IpAddresses []string      `bson:"ipAddresses" json:"ipAddresses"`
	AccessedAt  time.Time     `bson:"accessedAt" json:"accessedAt"`
	CreatedAt   time.Time     `bson:"createdAt" json:"-"`
	UpdatedAt   time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type LoginIpList struct {
	LoginIps         *[]LoginIP `json:"loginIps"`
	NumberOfLoginIps int        `json:"numberOfLoginIps"`
	CurrentPage      int        `json:"currentPage"`
	NumberOfPages    int        `json:"numberOfPages"`
}
