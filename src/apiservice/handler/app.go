package handler

import (
	"apiservice/config"
	"apiservice/utils"
	"go.mongodb.org/mongo-driver/mongo"

	"net/http"

	"github.com/go-redis/redis"
)

type Provider struct {
	Cfg  *config.Configuration
	RCli *redis.Client
	Db   *mongo.Database
}

func NewProvider(cfg *config.Configuration, rCli *redis.Client,db *mongo.Database) *Provider {
	return &Provider{
		Cfg:  cfg,
		RCli: rCli,
		Db: db,
	}
}

func (p *Provider) Ping (rw http.ResponseWriter, r *http.Request) {
		utils.RenderJson(rw, http.StatusOK, "Pong!.")
}