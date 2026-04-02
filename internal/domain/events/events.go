package events

type EventProducer interface {
	Publish(jsonEvent, topic string, key []byte) error
}
