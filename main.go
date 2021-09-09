package main

import (
	"log"
	"os"
	"os/signal"
	"reload-config/internal"
	"sync"
	"syscall"
)

func main() {

	// channel for stop, reset and shutdown signal
	stop := make(chan bool, 3)
	reset := make(chan bool)
	shutdown := make(chan bool)

	// used to synchronize subscriber, publisher and httpserver to finish
	var wg sync.WaitGroup
	wg.Add(3)

	log.Println("Starting Htpp Server")
	go internal.StartHttpServer(9090, &wg, stop, reset, shutdown)

	log.Println("Starting Subscriber")
	go internal.StartRabbitMqSubscriber(&wg, stop)

	log.Println("Starting Publisher")
	go internal.StartRabbitMqPublisher(&wg, stop, reset)

	// waiting for terminate request
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	select {
	case <-c:
		log.Println("terminate request from console")
	case <-shutdown:
		log.Println("terminate request from http request")
	}

	// sending stop to all (Subscriber, Publisher and HttpServer)
	log.Println("stopping Subscriber, Publisher and HttpServer")
	for i := 0; i < 3; i++ {
		stop <- true
	}

	// rendezvous
	wg.Wait()
	log.Println("Subscriber, Publisher and HttpServer all stop, Goodbye")
	os.Exit(0)
}
