package conn

import (
	"apiservice/config"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
	"time"
)

func TestInitDb(t *testing.T) {
	config.Parse("/config_files/apiservice_conf.json")
	InitDb()

	assert.IsType(t, &mongo.Database{}, Database)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := Database.Client().Ping(ctx, readpref.Primary()); err != nil {
		assert.Error(t, err)
	}

	assert.NotNil(t, Database)

}
