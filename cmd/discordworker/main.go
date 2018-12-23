package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"git.packetlostandfound.us/chiefnoah/fitbit-processor/pkg/common"
	nats "github.com/nats-io/go-nats"
)

func main() {
	natsHost := os.Getenv("NATS_HOST")
	discordWebhook := os.Getenv("DISCORD_WEBHOOK_URL")
	discordUsername := os.Getenv("DISCORD_WEBHOOK_USERNAME")
	log.Printf("Discord webhook URL set to %s", discordWebhook)

	log.Printf("Connecting to queue at %s", natsHost)
	nc, err := nats.Connect(natsHost)
	if err != nil {
		log.Fatalf("Unable to connect to queue: %s", err.Error())
	}
	q, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatalf("Unable to create encoded queue connection: %s", err.Error())
	}

	_, err = q.Subscribe(common.NatsQueueName, func(n *[]common.FitBitNotification) {
		log.Printf("Received message: %+v", n)
		content, _ := json.MarshalIndent(n, "", "  ")
		message := map[string]string{"content": fmt.Sprintf("Received fitbit update notification:\n```json\n%s\n```",
			content), "username": discordUsername}
		m, _ := json.Marshal(message)
		resp, err := http.Post(discordWebhook, "application/json", bytes.NewBuffer(m))
		if err != nil {
			log.Printf("Error sending discord webhook: %s", err.Error())
			return
		}
		log.Printf("Sent discord webhook. Response code: %s", resp.Status)
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
