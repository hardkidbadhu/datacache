package conn

import (
	"apiservice/config"
	"log"
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestConnectRedis(t *testing.T) {
	config.Parse("/config_files/apiservice_conf.json")
	cli := ConnectRedis(1)
	defer cli.Close()

	if cli != nil {

		// Making sure the connection alive
		if err := cli.Ping().Err(); err != nil {
			if cli != nil {
				cli.Close()
			}
			log.Fatalln(`ERROR: Couldn't connect to redis: `, err)
		}
	}

	assert.IsType(t, &redis.Client{}, cli)
	assert.NotNil(t, cli)
	assert.Equal(t, "ping: PONG", cli.Ping().String())
}