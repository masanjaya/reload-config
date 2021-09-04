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

// channel for reset and stop signal (don't care for the value)
var reset, stop, shutdown chan bool

// used to wait for subscriber to finish
var wg sync.WaitGroup

func main() {

	// create reset channel
	reset = make(chan bool)
	stop = make(chan bool, 3)
	shutdown = make(chan bool)

	log.Println("Starting Web Server")
	go startHttpServer()

	log.Println("Subscribe to MQ")
	wg.Add(3)
	go startRabbitMqSubscriber()

	log.Println("Publish to MQ")
	go startRabbitMqPublisher()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	select {
	case <-c:
		log.Println("terminate request from console")
	case <-shutdown:
		log.Println("terminate request from http request")
	}

	// sending stop to all (subscriber, publisher and http server)
	log.Println("stopping subsriber, publisher and HttpServer")
	for i := 0; i < 3; i++ {
		stop <- true
	}
	wg.Wait()
	log.Println("subsriber, publisher and HttpServer all stop")
	log.Println("Goodbye")
	os.Exit(0)
}

// startHttpServer starts an HTTP Server on port 9090, handle for path / and /reset
func startHttpServer() {
	defer wg.Done()

	// handle path /
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write([]byte("<html>click <a href='/reset'>reset</a> to reset the counter</html>"))
		if err != nil {
			log.Fatal(err)
		}
	})
	// handle path /reset
	http.HandleFunc("/reset", func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write([]byte("<html>Counter has been Resetet<br/><a href='/'>back to main page</a></html>"))
		if err != nil {
			log.Fatal(err)
		}
		// send reset signal
		log.Println("got request to reset publisher")
		reset <- true
	})
	// handle path /stop
	http.HandleFunc("/shutdown", func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write([]byte("<html>app is about to shutdown</html>"))
		if err != nil {
			log.Fatal(err)
		}
		// send shutdown signal
		log.Println("got request to shutdown the app")
		shutdown <- true
	})

	httpServer := http.Server{
		Addr: ":9090",
	}

	// monitor stop signal
	go func() {
		<-stop
		log.Println("Http Server is shutting down")
		httpServer.Shutdown(context.Background())
	}()

	log.Println("Http Server is starting")
	err := httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}

	log.Println("Http Server is shutdown")
}

// startRabbitMqSubscriber subscribes to a local rabbitmq
// use docker for local rabbitmq installation
func startRabbitMqSubscriber() {
	defer wg.Done()

	log.Println("subscriber establishing connection to broker")
	mq, err := rabbitmq.NewMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("subscriber connecting to rabbitmq failed => ", err)
	}
	defer func() {
		log.Println("subscriber closing connection to broker")
		mq.Close()
	}()

	_, err = mq.QueueDeclare(rabbitmq.NewQueueOptions().SetName("hello"))
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

	_, err = mq.QueueDeclare(rabbitmq.NewQueueOptions().SetName("hello"))
	if err != nil {
		log.Fatal(err)
	}

	err = mq.ExchangeDeclare(rabbitmq.NewExchangeOptions().SetName("hello").SetType(rabbitmq.ExchangeTypeFanout))
	if err != nil {
		log.Fatal(err)
	}

	err = mq.QueueBind(rabbitmq.NewQueueBindOptions().SetName("hello").SetExchange("hello"))
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
					SetExchange("hello").
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
