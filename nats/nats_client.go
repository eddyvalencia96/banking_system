package nats_client

import (
	"github.com/nats-io/nats.go"
)

// NewClient crea y devuelve una nueva instancia del cliente NATS
func NewClient(url string) (*nats.Conn, error) {
	nc, err := nats.Connect(url, nats.Name("NATS Client"))
	if err != nil {
		return nil, err
	}
	return nc, nil
}

// CreateStream crea un nuevo flujo (stream) en JetStream
func CreateStream(nc *nats.Conn, streamName string, subjects []string) error {
	js, err := nc.JetStream()
	if err != nil {
		return err
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: subjects,
	})
	return err
}

// PublishToStream publica un mensaje en un flujo específico de JetStream
func PublishToStream(nc *nats.Conn, subject string, message []byte) error {
	js, err := nc.JetStream()
	if err != nil {
		return err
	}

	_, err = js.Publish(subject, message)
	return err
}
