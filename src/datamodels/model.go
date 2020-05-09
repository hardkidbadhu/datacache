package datamodels

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserInfo struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	SNo         float64            `json:"sNo" bson:"sNo"`
	Name        string             `json:"name" bson:"name"`
	Age         int32              `json:"age" bson:"age"`
	Email       string             `json:"email" bson:"email"`
	PhoneNumber string             `json:"phoneNumber" bson:"phoneNumber"`
}

type Pagination struct {
	PageNum   int    `json:"pageNum,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Count     int    `json:"count,omitempty"`
	NextPage  int    `json:"nextPage,omitempty"`
	PrevPage  int    `json:"prevPage,omitempty"`
	TotalPage int    `json:"pageCount,omitepmty"`
}