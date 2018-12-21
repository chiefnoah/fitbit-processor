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
	if port == "" {
		port = "8008"
	}
	listerService := listener.Service{
		SubscriberVerificationCode: os.Getenv("SUBSCRIBER_VERIFICATION_CODE"),
	}
	http.HandleFunc("/webhook", listerService.HandleFitbitNotification)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
