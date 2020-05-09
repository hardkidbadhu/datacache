package data_seeds

import (
	"context"
	"datamodels"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strconv"
)

func InitSeeds(db *mongo.Database) {
	log.Println("Intializing data seeders...")

	for i := 1; i < 21; i ++ {
		u := datamodels.UserInfo{}
		u.Name = "Sample_name_" + strconv.Itoa(i)
		u.Id = primitive.NewObjectID()
		u.SNo = float64(i)
		u.Age = int32(i + 1)
		u.Email = "sampleemail@" + strconv.Itoa(i) + ".com"
		u.PhoneNumber = "9998899" + strconv.Itoa(i)
		db.Collection("userinfo").InsertOne(context.Background(), u)
	}

}
