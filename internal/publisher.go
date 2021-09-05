package internal

import (
	"fmt"
	"log"
	"sync"
	"time"

	rabbitmq "github.com/hadihammurabi/go-rabbitmq"
	"github.com/streadway/amqp"
)

// StartRabbitMqPublisher publish to a local rabbitmq
// use docker for local rabbitmq installation
func StartRabbitMqPublisher(wg *sync.WaitGroup, stop, reset chan bool) {
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
