package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	rabbitmq "github.com/hadihammurabi/go-rabbitmq"
	"github.com/streadway/amqp"
)

// channel for reset, stop and shutdown signal
var reset, stop, shutdown chan bool

// used to wait for subscriber, publisher and httpserver to finish
var wg sync.WaitGroup

func main() {

	// create channels
	reset = make(chan bool)
	stop = make(chan bool, 3)
	shutdown = make(chan bool)

	log.Println("Starting Htpp Server")
	go startHttpServer()

	log.Println("Starting Subscriber")
	wg.Add(3)
	go startRabbitMqSubscriber()

	log.Println("Starting Publisher")
	go startRabbitMqPublisher()

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
	wg.Wait()
	log.Println("Subscriber, Publisher and HttpServer all stop, Goodbye")
	os.Exit(0)
}

// startHttpServer starts an HTTP Server on port 9090, handle for path / and /reset
func startHttpServer() {
	defer wg.Done()

	// handle path /
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write([]byte("<html>click <a href='/reset'>reset</a> to reset Publisher</html>"))
		if err != nil {
			log.Fatal(err)
		}
	})
	// handle path /reset
	http.HandleFunc("/reset", func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write([]byte("<html>Publisher has been reset<br/><a href='/'>back to main page</a></html>"))
		if err != nil {
			log.Fatal(err)
		}
		// send reset signal
		log.Println("got http request to reset publisher")
		reset <- true
	})
	// handle path /stop
	http.HandleFunc("/shutdown", func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write([]byte("<html>app is about to shutdown</html>"))
		if err != nil {
			log.Fatal(err)
		}
		// send shutdown signal
		log.Println("got http request to shutdown the app")
		shutdown <- true
	})

	httpServer := http.Server{
		Addr: ":9090",
	}

	// monitor stop signal
	go func() {
		<-stop
		log.Println("HttpServer is shutting down")
		httpServer.Shutdown(context.Background())
	}()

	log.Println("HttpServer is starting")
	err := httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}

	log.Println("HttpServer is shutdown")
}

// startRabbitMqSubscriber subscribes to a local rabbitmq
// use docker for local rabbitmq installation
func startRabbitMqSubscriber() {
	defer wg.Done()

	log.Println("subscriber establishing connection to broker")
	mq, err := rabbitmq.NewMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("subscriber connecting to rabbitmq failed", err)
	}
	defer func() {
		log.Println("subscriber closing connection to broker")
		mq.Close()
	}()

	_, err = mq.QueueDeclare(rabbitmq.NewQueueOptions().SetName("test"))
	if err != nil {
		log.Fatal(err)
	}

	// consume message on queue (this runs forever and ever and ever.. until.. the end of time #halah)
	msgs, err := mq.Consume(nil)
	if err != nil {
		log.Fatal(err)
	}

	// show the message
	for msg := range msgs {
		select {
		case <-stop:
			return
		default:
			log.Println("subscriber processing", string(msg.Body), "=== start ===")
			for i := 1; i <= 5; i++ {
				log.Println("subscriber processing", string(msg.Body), "phase", i)
				time.Sleep(time.Second)
			}
			log.Println("subscriber processing", string(msg.Body), ">>>>> finish <<<<<")
			msg.Ack(false)
		}
	}
}

func startRabbitMqPublisher() {
	defer wg.Done()

	mq, err := rabbitmq.NewMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer mq.Close()

	_, err = mq.QueueDeclare(rabbitmq.NewQueueOptions().SetName("test"))
	if err != nil {
		log.Fatal(err)
	}

	err = mq.ExchangeDeclare(rabbitmq.NewExchangeOptions().SetName("test").SetType(rabbitmq.ExchangeTypeFanout))
	if err != nil {
		log.Fatal(err)
	}

	err = mq.QueueBind(rabbitmq.NewQueueBindOptions().SetName("test").SetExchange("test"))
	if err != nil {
		log.Fatal(err)
	}

	// the counters
	var i int = 0

	for {
		select {
		case <-stop: // stop signal received
			return
		case <-reset: // reset signal received
			i = 0
		default: // do normal publishing
			i++
			log.Println("publisher sending", i)
			err = mq.Publish(
				rabbitmq.NewPublishOptions().
					SetExchange("test").
					SetMessage(amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(fmt.Sprintf("%d", i)),
					}),
			)
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(3 * time.Second)
		}
	}
}
