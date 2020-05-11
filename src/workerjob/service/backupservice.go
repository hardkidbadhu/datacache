package service

import (
	"context"
	"datamodels"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"log"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ProcessDataBackUp(database *mongo.Database, rCli *redis.Client) error {

	cmd := rCli.ZRangeByScore("userInfo", &redis.ZRangeBy{Min: "-inf", Max: "+inf"})

	data := cmd.Val()

	if len(data) <= 0 {
		log.Println("Data: ProcessDataBackUp - No data found returning")
		return nil
	}

	if cmd.Err() != nil {
		return cmd.Err()
	}

	for i := range data {
		uInfIns := datamodels.UserInfo{}
		if err := json.Unmarshal([]byte(data[i]), &uInfIns); err != nil {
			log.Printf("Error - ProcessDataBackUp - %s", err.Error())
			continue
		}

		//upserts the document
		mR, err := database.Collection("userInfo").UpdateOne(context.Background(), bson.D{{"_id", uInfIns.Id}},
			bson.D{{"$set", uInfIns}}, &options.UpdateOptions{Upsert: &[]bool{true}[0]})
		if err != nil {
			log.Printf("Error - ProcessDataBackUp - %s", err.Error())
			return err
		}

		log.Printf("Info - ProcessDataBackUp - update many [Matched- %d] - [Modified- %d] [UpsertedCount- %d] [UpsertedID- %+v]" ,
			mR.MatchedCount, mR.ModifiedCount, mR.UpsertedCount, mR.UpsertedID)
	}

	rCli.ZRemRangeByScore("userInfo", "-inf", "+inf")
	return nil
}