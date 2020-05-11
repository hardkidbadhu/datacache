package main

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/go-chi/chi"
	"net/http"
	"os/signal"
	"strconv"
	"time"
	"workerservice/config"
	"workerservice/conn"
	"workerservice/service"

	"context"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"runtime"
	"sync"
)

var (
	ctl   = &sync.WaitGroup{}
)

func main() {

	//passing the global config files as a flag during the runtime
	var configFile = flag.String("conf", "", "configuration file")
	flag.Parse()
	if flag.NFlag() != 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.Println("INFO: Initializing worker server...")

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Configuring Goroutine to use all the available CPUs
	if nCPU := runtime.NumCPU(); nCPU > 1 {
		prevnCPU := runtime.GOMAXPROCS(nCPU)
		log.Println("INFO: prevnCPU, nCPU - ", prevnCPU, nCPU)
	}


	//Parsing configuration file
	log.Println("Parsing configuration...")
	config.Parse(*configFile)

	rCli := conn.ConnectRedis(1)
	defer rCli.Close()

	log.Println("INFO: Initializing mongodb master session.")
	conn.InitDb()
	defer conn.Database.Client().Disconnect(context.Background())
	log.Println("INFO: Mongodb connected.")


	defer func() {
		if err := recover(); err != nil {
			log.Println(`panic:main`, err)
		}
	}()

	ctl.Add(2)
	spawn(WebServer,  conn.Database, rCli, "worker: Webserver") //worker - 1
	spawn(RestoreToCache, conn.Database, rCli, "worker: Mongo->Redis") //worker - 2

	ctl.Wait()

	log.Println("Shutting down worker service...")
}

//Helper function that takes in func and deploys child worker go-routines
func spawn(f func(database *mongo.Database, rCli *redis.Client), database *mongo.Database,
	rCli *redis.Client, label string) {
	go func() {
		log.Printf("Info: Spawning - %s", label)
		defer func() {
			if err := recover(); err != nil {
				log.Println(`panic: `+label, err)
				f(database, rCli)
			}
		}()
		f(database, rCli)
	}()
}

//DbBackup worker for persistence db back up every 15minutes from redis -> mongodb
func WebServer(database *mongo.Database, rCli *redis.Client) {
	defer ctl.Done()

	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
		}{
			"Hey buddy! Am alright...",
		})
	})

	r.Get("/db_stats", func(w http.ResponseWriter, r *http.Request) {

		c := rCli.Get("counter")
		zr := rCli.ZRangeByScore("userInfo", &redis.ZRangeBy{Min: "-inf", Max: "+inf"})

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			RecordsInCache string `json:"records_in_cache"`
			RecordsInDB    string `json:"records_in_db"`
		}{
			RecordsInCache: strconv.Itoa(len(zr.Val())),
			RecordsInDB: c.Val(),
		})
	})

	server := &http.Server{
		Addr:              ":9096",
		Handler:           r,
		ReadHeaderTimeout: 15 * time.Second,
		WriteTimeout:      150 * time.Second,
		IdleTimeout:       180 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	log.Println("Listening server on: ", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Listen: %s\n", err)
	}

	log.Println("Server stopped...")
}

//RestoreToCache restores the data in cache on notification from message brokers
func RestoreToCache(database *mongo.Database, rCli *redis.Client) {
	defer ctl.Done()

	var (
		err error
		con sarama.Consumer
		part sarama.PartitionConsumer
		counter int
		consumerMessage *sarama.ConsumerMessage
		sy sync.WaitGroup
	)

	conf := sarama.NewConfig()
	conf.Consumer.Return.Errors = true

	//Consumer group to consume the topic that is produced by the server
	con, err = sarama.NewConsumer(config.Cfg.Kafka.BrokerAddr, conf)
	if err != nil {
		log.Printf("Error: RestoreToCache - %s", err.Error())
		return
	}

	defer func() {
		if err := con.Close(); err != nil {
			log.Printf("Error: RestoreToCache - %s", err.Error())
			return
		}
	}()

	part, err = con.ConsumePartition(config.Cfg.Kafka.Topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Printf("Error: RestoreToCache - %s", err.Error())
		return
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	//counter to hold the no of messages processed
	sy.Add(1)
	go func() {
		defer sy.Done()
		for {
			select {
			case err := <-part.Errors():
				fmt.Println(err)
			case consumerMessage = <-part.Messages():
				log.Println(part.Messages())
				fmt.Println("Received messages", string(consumerMessage.Key), string(consumerMessage.Value))
				if err := service.ProcessDataSync(consumerMessage, database, rCli); err != nil {
					log.Printf("Error: Data restoring error - %s", err.Error())
					continue
				}
				counter++
			case <-signals:
				log.Println("Interupt deducted...")
				log.Println("Service going down!..")
				os.Exit(1)
				return
			}
		}
	}()

	sy.Wait()
	fmt.Printf("Info: Processed - %d, messages", counter)
}

