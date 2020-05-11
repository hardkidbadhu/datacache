package datamodels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserInfo struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	SNo         float64            `json:"sNo" bson:"sNo"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Age         int32              `json:"age,omitempty" bson:"age,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	PhoneNumber string             `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}

type Pagination struct {
	PageNum   int    `json:"pageNum,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Count     int    `json:"count,omitempty"`
}