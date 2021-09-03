package main

import (
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

// channel for reset signal (don't care for the value)
var reset, stop chan bool

// WG
var wg sync.WaitGroup

func main() {

	// create reset channel
	reset = make(chan bool)
	stop = make(chan bool)

	log.Println("Starting Web Server")
	go startHttpServer()

	log.Println("Subscribe to MQ")
	wg.Add(1)
	go startRabbitMqSubscriber()

	log.Println("Publish to MQ")
	go startRabbitMqPublisher()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	<-c

	// sending signal to channel stop (received by subscriber and publisher)
	log.Println("receiving SIGINT, sending stop signal to subsriber (don't care for publisher and HttpServer)")
	stop <- true

	log.Println("waiting for subscriber to stop")
	wg.Wait()

	log.Println("subscriber has stop")
	log.Println("Goodbye")
	os.Exit(0)
}

// startHttpServer starts an HTTP Server on port 9090, handle for path / and /reset
func startHttpServer() {
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
		reset <- true
		log.Println("reset signal sent")
	})
	log.Fatal(http.ListenAndServe(":9090", nil))
}

// startRabbitMqSubscriber subscribes to a local rabbitmq
// use docker for local rabbitmq installation
func startRabbitMqSubscriber() {
	mq, err := rabbitmq.NewMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("connecting to rabbitmq failed => ", err)
	}

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
			log.Println("closing connection to broker")
			mq.Close()
			wg.Done()
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
		// reset signal received
		case <-reset:
			i = 0
		// do normal publishing
		default:
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
