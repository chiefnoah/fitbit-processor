package listener

import (
	"log"
	"net/http"
)

// FitbitListenerService defines an instance of a listner service for receving data events from Fitbit
type Service struct {
	SubscriberVerificationCode string
}

//HandleFitbitNotification handles a notification event from Fitbit
func (s *Service) HandleFitbitNotification(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling request from fitbit: %+f", r)
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
	w.WriteHeader(http.StatusAccepted)
}
