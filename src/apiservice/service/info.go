package service

import (
	"apiservice/config"
	"datamodels"
	"github.com/go-redis/redis"
)

var _ InfoService = &infoSrv{}

type InfoService interface {
	//FetchDataFromCache returns the requested data if available in cache
	FetchDataFromCache(*datamodels.Pagination) ([]*datamodels.UserInfo, error)

	//SendNotificationToKafka If data is not present in cache it sends a notification to kafka where the data is
	//restored in cache
	SendNotificationToKafka()
}

type infoSrv struct {
	rCli *redis.Client
	cfg *config.Configuration
}

func NewInfoService(r *redis.Client, c *config.Configuration) *infoSrv {
	return &infoSrv{
		r,
		c,
	}
}

func (i infoSrv) FetchDataFromCache(p *datamodels.Pagination) ([]*datamodels.UserInfo, error){

	return nil, nil
}

func (i infoSrv) SendNotificationToKafka() {
	panic("implement me")
}
