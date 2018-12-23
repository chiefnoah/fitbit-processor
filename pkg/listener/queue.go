package listener

import (
	"log"

	"git.packetlostandfound.us/chiefnoah/fitbit-processor/pkg/common"
	nats "github.com/nats-io/go-nats"
)

//InitQueue establishes a connection to a queue server
func (s *Service) InitQueue(host string) error {
	log.Printf("Connecting to queue at %s", host)
	nc, err := nats.Connect(host)
	if err != nil {
		return err
	}
	q, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return err
	}
	s.q = q

	return nil
}

//Cleanup cleans up connections set up for a Service to run
func (s *Service) Cleanup() {
	log.Print("Closing connection to queue")
	s.q.Close()
}

//PublishNotification sends a fitbit notification to the queue for processing
func (s *Service) PublishNotification(n *common.FitBitNotification) error {
	return s.q.Publish(common.NatsQueueName, n)
}
