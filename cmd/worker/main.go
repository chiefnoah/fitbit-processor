package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"git.packetlostandfound.us/chiefnoah/fitbit-processor/pkg/common"
	nats "github.com/nats-io/go-nats"
)

func main() {
	natsHost := os.Getenv("NATS_HOST")

	log.Printf("Connecting to queue at %s", natsHost)
	nc, err := nats.Connect(natsHost)
	if err != nil {
		log.Fatalf("Unable to connect to queue: %s", err.Error())
	}
	q, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatalf("Unable to create encoded queue connection: %s", err.Error())
	}

	_, err = q.QueueSubscribe(common.NatsQueueName, common.NatsQueueGroup, func(n *common.FitBitNotification) {
		log.Printf("Received message: %+v", n)
	})
	if err != nil {
		log.Fatalf("Unable to subscrbie to %s as queue group %s: %s", common.NatsQueueName, common.NatsQueueGroup, err.Error())
	}

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	// block here until we get told to shutdown
	sig := <-gracefulStop
	fmt.Printf("caught sig: %+v\n", sig)
	fmt.Println("closing queue connection...")
	q.Close()
	os.Exit(0)
}
