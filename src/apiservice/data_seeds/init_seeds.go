package data_seeds

import (
	"apiservice/constants"
	"context"
	"datamodels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strconv"
	"time"
)

func InitSeeds(db *mongo.Database) {
	log.Println("Intializing data seeders...")

	//Load data only if collection is empty
	c, _ := db.Collection(constants.COLUserInfo).CountDocuments(context.Background(), bson.D{})
	if c > 0 {
		return
	}

	for i := 1; i < 21; i ++ {
		u := datamodels.UserInfo{}
		u.Name = "Sample_name_" + strconv.Itoa(i)
		u.Id = primitive.NewObjectID()
		u.SNo = float64(i)
		u.Age = int32(i + 1)
		u.Email = "sampleemail@" + strconv.Itoa(i) + ".com"
		u.PhoneNumber = "9998899" + strconv.Itoa(i)
		u.CreatedAt = time.Now().UTC()

		db.Collection(constants.COLUserInfo).InsertOne(context.Background(), u)
	}

}
