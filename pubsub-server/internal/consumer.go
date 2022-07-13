package internal

type ConsumerMessage struct {
	Message any
	Done    chan struct{}
}

type Consumer interface {
	GetMessageChan() chan ConsumerMessage
	ShutdownChan() chan struct{}
	GetID() uint64
}

type ConsumerFactory interface {
	NewConsumer() Consumer
}

type ConsumerMatcher interface {
	MatchTopic(topic string) []Consumer
	AddConsumers(topic string, consumers ...Consumer)
	RemoveConsumerFromTopic(topic string, consumer Consumer)
}
