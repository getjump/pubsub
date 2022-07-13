package internal

import "sync"

type SimpleConsumerFactory struct {
	ID uint64
	sync.Mutex
}

type SimpleConsumer struct {
	messageChan  chan ConsumerMessage
	shutdownChan chan struct{}
	ID           uint64
}

func (factory *SimpleConsumerFactory) NewConsumer() Consumer {
	factory.Mutex.Lock()
	c := SimpleConsumer{
		messageChan:  make(chan ConsumerMessage),
		shutdownChan: make(chan struct{}),
		ID:           factory.ID,
	}
	factory.ID += 1
	factory.Mutex.Unlock()

	return c
}

func (consumer SimpleConsumer) GetMessageChan() chan ConsumerMessage {
	return consumer.messageChan
}

func (consumer SimpleConsumer) ShutdownChan() chan struct{} {
	return consumer.shutdownChan
}

func (consumer SimpleConsumer) GetID() uint64 {
	return consumer.ID
}
