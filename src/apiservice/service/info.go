package service

import (
	"apiservice/Errs"
	"apiservice/config"
	"apiservice/constants"
	"apiservice/utils"
	"time"

	"context"
	"datamodels"
	"encoding/json"
	"log"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ InfoService = &infoSrv{}

type InfoService interface {
	//AddDataToCache adds the newly generated data into cache
	AddDataToCache(*datamodels.UserInfo) error

	//FetchDataFromCache returns the requested data if available in cache
	FetchDataFromCache(*datamodels.Pagination) (interface{}, error)

	//FetchDataFromDB fetches the requested data from database
	FetchDataFromDB(int64, int64) ([]datamodels.UserInfo, error)

	//SendNotificationToKafka If data is not present in cache it sends a notification to kafka where the data is
	//restored in cache
	SendNotificationToKafka(int, int)

	//SetCounterData sets the counter in cache to paginate data
	SetCounterData() error
}

type infoSrv struct {
	rCli *redis.Client
	col  *mongo.Collection
}

func NewInfoService(r *redis.Client, col *mongo.Collection) *infoSrv {
	return &infoSrv{
		r,
		col,
	}
}

func (i infoSrv) AddDataToCache(info *datamodels.UserInfo) error {

	info.Id = primitive.NewObjectID()
	ccmd := i.rCli.Get(constants.RecCounter)
	info.SNo, _ = strconv.ParseFloat(ccmd.Val(), 64)
	info.CreatedAt = time.Now().UTC()

	byt, err := json.Marshal(info)
	if err != nil {
		return &Errs.AppErr{
			Message: "Something went wrong, Please try after sometime!.",
			Err:     err.Error(),
		}
	}

	cmd := i.rCli.ZAdd(constants.RedisPrefixKey, &redis.Z{Score: info.SNo+1, Member: string(byt)})
	if cmd.Err() != nil {
		return &Errs.AppErr{
			Message: "Error in adding members!.",
			Err:     cmd.Err().Error(),
		}
	}

	i.rCli.Set(constants.RecCounter, info.SNo+1, 0)

	return nil
}

func (i infoSrv) FetchDataFromCache(p *datamodels.Pagination) (interface{}, error) {

	counterCmd := i.rCli.Get(constants.RecCounter)
	totalRec, _ := strconv.Atoi(counterCmd.Val())


	pageRec := utils.TotalRecordsInPage(p.PageNum, p.Limit, totalRec)
	if pageRec <= 0 {
		return nil, nil
	}

	min := (p.PageNum - 1) * p.Limit
	max := min + pageRec

	p.Count = totalRec

	cmd := i.rCli.ZRangeByScore(constants.RedisPrefixKey,
		&redis.ZRangeBy{Min: strconv.Itoa(min), Max: strconv.Itoa(max)})
	if len(cmd.Val()) != pageRec {
		//Async call to apache notification worker
		go i.SendNotificationToKafka(p.PageNum, p.Limit)
		return i.FetchDataFromDB(int64((p.PageNum - 1) * p.Limit), int64(p.Limit))
	}

	uinf := []datamodels.UserInfo{}
	for p := range cmd.Val() {
		uIns := datamodels.UserInfo{}

		if err := json.Unmarshal([]byte(cmd.Val()[p]), &uIns); err != nil {
			log.Printf("Error: Unmarshalling redis data - %s", err.Error())
			continue
		}

		uinf = append(uinf, uIns)
	}

	return uinf, nil

}

func (i infoSrv) SendNotificationToKafka(pno, limit int) {

	log.Println("SendNotificationToKafka", pno, limit)
	kfCfg := sarama.NewConfig()

	kfCfg.Producer.RequiredAcks = sarama.WaitForAll
	kfCfg.Producer.Retry.Max = 10
	kfCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(config.Cfg.Kafka.BrokerAddr, kfCfg)
	if err != nil {
		log.Println("Error: SendNotificationToKafka", err.Error())
		return
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Println("Error: SendNotificationToKafka", err.Error())
			return
		}
	}()

	//Encoded json message
	byt, _ := json.Marshal(struct {
		PageNo int `json:"pageNo"`
		Limit  int `json:"limit"`
	}{
		pno,
		limit,
	})

	msg := &sarama.ProducerMessage{
		Topic: config.Cfg.Kafka.Topic,
		Key:   sarama.StringEncoder("userData"),
		Value: sarama.ByteEncoder(byt),
	}

	pt, of, err := producer.SendMessage(msg)
	if err != nil {
		log.Println("Error: SendNotificationToKafka", err.Error())
		return
	}

	log.Printf("Message published to topic - %s ,partition - %d, offset - %d",
		config.Cfg.Kafka.Topic, pt, of)

}

func (i infoSrv) SetCounterData() error {
	count, err := i.col.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return err
	}

	i.rCli.Set(constants.RecCounter, count, 0)
	return nil
}

func (i infoSrv) FetchDataFromDB(skip, limit int64) (data []datamodels.UserInfo, err error) {

	opts := &options.FindOptions{Skip: &[]int64{skip}[0], Limit: &[]int64{limit}[0]}

	cur := &mongo.Cursor{}
	if cur, err = i.col.Find(context.Background(), bson.D{}, opts); err != nil {
		log.Printf("Error: FetchDataFromDB - fetching data - %s", err.Error())
		return
	}

	if err = cur.All(context.Background(), &data); err != nil {
		log.Printf("Error: FetchDataFromDB - cursor conversion - %s", err.Error())
		return
	}

	log.Println("data", data, err)
	return
}
