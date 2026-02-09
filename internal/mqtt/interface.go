package mqtt

import "time"

type Publisher interface {
	SetAvailability(topic string, onlinePayload []byte)
	Connect(timeout time.Duration) error
	Publish(topic string, qos byte, retain bool, payload []byte) error
	Close()
}
