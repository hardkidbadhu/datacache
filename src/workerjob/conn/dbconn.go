package conn

import (
	"workerjob/config"

	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	// Obj defines the mongodb session, which connects mongodb instance.
	Database    *mongo.Database
	once   sync.Once
)

func InitDb() {
	once.Do(func() {
		Database = connectDB()
	})
}

//ConnectLocalDB - Connects the local mongodb with supplied uri
func connectDB() *mongo.Database {

	clientOptions := options.Client().ApplyURI(config.Cfg.Database.URI)

	if config.Cfg.Database.UserName != "" && config.Cfg.Database.Source != "" &&
		config.Cfg.Database.Password != "" {
		clientOptions.Auth = &options.Credential{
			AuthMechanism: `SCRAM-SHA-1`,
			AuthSource:    config.Cfg.Database.Source,
			Username:      config.Cfg.Database.UserName,
			Password:      config.Cfg.Database.Password,
		}
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error: ConnectLocalDB - %s", err.Error())
		return nil
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Error: ConnectLocalDB - %s", err.Error())
		return nil
	}

	return client.Database(config.Cfg.Database.Name)
}
