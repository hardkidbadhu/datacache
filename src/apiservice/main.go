package main

import (
	"apiservice/config"
	"apiservice/conn"
	"apiservice/data_seeds"
	"apiservice/handler"
	"apiservice/router"

	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	//passing the global config files as a flag during the runtime
	var configFile = flag.String("conf", "", "configuration file")
	flag.Parse()
	if flag.NFlag() != 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	//Configuring logger for the app
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.New(os.Stdout, "datacache_api_service: ", log.LstdFlags|log.Lshortfile)

	log.Println("Initializing api service...")
	//Parsing configuration file
	log.Println("Parsing configuration...")
	config.Parse(*configFile)

	cli := conn.ConnectRedis(1)
	defer cli.Close()

	log.Println("Connecting database...")
	conn.InitDb()
	defer conn.Database.Client().Disconnect(context.Background())

	data_seeds.InitSeeds(conn.Database)

	p := handler.NewProvider(config.Cfg, cli, conn.Database)
	r := router.NewRouter(p)

	server := &http.Server{
		Addr:           config.Cfg.HTTPAddress + ":" + config.Cfg.Port,
		Handler:        r,
		ReadTimeout:    time.Duration(config.Cfg.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.Cfg.WriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(config.Cfg.MaxIdleTimeOut) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Graceful shut down of server
	graceful := make(chan os.Signal)
	signal.Notify(graceful, syscall.SIGINT)
	signal.Notify(graceful, syscall.SIGTERM)
	go func() {
		<-graceful
		log.Println("Shutting down server...")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not do graceful shutdown: %v\n", err)
		}
	}()

	log.Println("Listening server on ", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Listen: %s\n", err)
	}

	log.Println("Server gracefully stopped...")
}
