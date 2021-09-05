package internal

import (
	"log"
	"sync"
	"time"

	rabbitmq "github.com/hadihammurabi/go-rabbitmq"
)

// StartRabbitMqSubscriber subscribes to a local rabbitmq
// use docker for local rabbitmq installation
func StartRabbitMqSubscriber(wg *sync.WaitGroup, stop chan bool) {
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
