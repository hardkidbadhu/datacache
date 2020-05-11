package service

import (
	"workerservice/dal"

	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"

)

func ProcessDataSync(message *sarama.ConsumerMessage, database *mongo.Database, rCli *redis.Client) error {

	//key := string(message.Key)
	//if key != "userData" {
	//	return fmt.Errorf("error: ProcessDataSync - No key found in the message - %+v", message)
	//}

	res := struct {
		PageNo int `json:"pageNo"`
		Limit  int `json:"limit"`
	}{}

	if err := json.Unmarshal(message.Value, &res); err != nil {
		return fmt.Errorf("error: ProcessDataSync - %s", err.Error())
	}

	//fetch requested data form database
	val, err := dal.NewUserInfDal(database.Collection("userInfo")).GetData(res.PageNo, res.Limit)
	if err != nil {
		return err
	}

	//After values are received from db restore into redis with the particular score(sno)
	for i := range val {
		byt, _ := json.Marshal(*val[i])

		//validates the record if already present in cache or not
		if ccmd := rCli.ZRangeByScore("userInfo", &redis.ZRangeBy{}); len(ccmd.Val()) == 0 {
			cmd := rCli.ZAdd("userInfo", &redis.Z{Score: val[i].SNo, Member: string(byt)})
			if cmd.Err() != nil {
				log.Printf("Error: ProcessDataSync - %s", cmd.Err().Error())
			}
		}

		log.Printf("Info: ProcessDataSync - Data already exists - %+v", val[i])
	}

	return nil
}
