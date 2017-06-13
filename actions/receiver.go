package actions

import (
	"log"

	"github.com/vice-registry/vice-worker/common"
)

// WaitForActions listens on RabbitMQ channel and accepts one message at a time
func WaitForActions() error {

	queue, err := rabbitmqCredentials.Channel.QueueDeclare(
		adaptors.WorkerType, // name
		true,                // durable
		false,               // delete when usused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		log.Printf("Unable to connect to RabbitMQ: %s", err)
		return err
	}

	err = rabbitmqCredentials.Channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Printf("Unable to connect to RabbitMQ: %s", err)
		return err
	}

	msgs, err := rabbitmqCredentials.Channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Printf("Unable to connect to RabbitMQ: %s", err)
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			reference := string(d.Body)
			action := Action{
				reference: reference,
			}
			err := handleAction(action)
			if err != nil {
				log.Printf("Failed to handle action for reference %s: %s", reference, err)
			} else {
				d.Ack(false)
			}
		}
	}()

	// wait forever until interrupted
	<-forever

	return nil
}
