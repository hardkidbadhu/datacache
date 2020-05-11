package main

import (
	"workerjob/config"
	"workerjob/conn"
	"workerjob/service"

	"context"
	"flag"
	"log"
	"os"
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

	//Parsing configuration file
	log.Println("Parsing configuration...")
	config.Parse(*configFile)

	rCli := conn.ConnectRedis(1)
	defer rCli.Close()

	log.Println("INFO: Initializing mongodb master session.")
	conn.InitDb()
	defer conn.Database.Client().Disconnect(context.Background())

	log.Println("INFO: Mongodb connected...")

	service.ProcessDataBackUp(conn.Database, rCli)

	log.Println("INFO: DB backup completed...")
}
