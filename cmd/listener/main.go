package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"git.packetlostandfound.us/chiefnoah/fitbit-processor/pkg/listener"
)

func main() {
	port := os.Getenv("HTTP_PORT")
	natsHost := os.Getenv("NATS_HOST")
	if port == "" {
		port = "8008"
	}
	listenerService := listener.Service{
		SubscriberVerificationCode: os.Getenv("SUBSCRIBER_VERIFICATION_CODE"),
	}

	err := listenerService.InitQueue(natsHost)
	if err != nil {
		log.Fatalf("Unable to connect to queue: %s", err.Error())
	}
	defer listenerService.Cleanup()
	http.HandleFunc("/webhook", listenerService.HandleFitbitNotification)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
