package listener

import (
	"encoding/json"
	"log"
	"net/http"

	"git.packetlostandfound.us/chiefnoah/fitbit-processor/pkg/common"

	nats "github.com/nats-io/go-nats"
)

// Service defines an instance of a listner service for receving data events from Fitbit
type Service struct {
	SubscriberVerificationCode string
	q                          *nats.EncodedConn
}

//HandleFitbitNotification handles a notification event from Fitbit
func (s *Service) HandleFitbitNotification(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling request from fitbit: %s", r.URL.Path)
	verify, ok := r.URL.Query()["verify"]
	if ok {
		if verify[0] == s.SubscriberVerificationCode {
			log.Println("Verifying!")
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("Valid"))
			return
		} else {
			log.Println("Verifying bad!")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Invalid"))
			return
		}
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var notifications []common.FitBitNotification
	err := decoder.Decode(&notifications)
	if err != nil {
		log.Printf("Called: %s - Invalid JSON payload: %s", r.URL.Path, err)
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}
	for _, noticiation := range notifications {
		err = s.PublishNotification(&noticiation)
		if err != nil {
			log.Printf("Error queuing notification: %s", err.Error())
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
