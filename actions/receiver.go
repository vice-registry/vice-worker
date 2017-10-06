package actions

import (
	"log"

	"github.com/vice-registry/vice-util/communication"
	"github.com/vice-registry/vice-worker/common"
)

// WaitForActions listens on RabbitMQ channel and accepts one message at a time
func WaitForActions() error {

	msgs, err := communication.NewConsumer(adaptors.WorkerType)
	if err != nil {
		log.Printf("Error while registering new consumer: %s", err)
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
