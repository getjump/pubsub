package internal

import (
	"context"
	"sync"
)

type Processor struct {
	ConsumerMatcher

	consumerFactory ConsumerFactory
}

// NewProcessor returns new instance of Processor.
func NewProcessor(matcher ConsumerMatcher, consumerFactory ConsumerFactory) *Processor {
	processor := Processor{
		ConsumerMatcher: matcher,
		consumerFactory: consumerFactory,
	}

	return &processor
}

// Publish Producer can Publish Message to topic and it should be delivered to all consumers
// SYNC should wait for message delivery to every consumer (also with validation)
// ASYNC should return result after processing and spawn goroutines to send in parallel
// we only guarantee that consumers that subscribed before call to that RPC will receive that message.
func (p *Processor) PublishAsync(ctx context.Context, topic string, message any) error {
	consumers := p.ConsumerMatcher.MatchTopic(topic)

	for _, consumer := range consumers {
		go func(ctx context.Context, out chan ConsumerMessage, shutdownChan chan struct{}) {
			select {
			case out <- ConsumerMessage{Message: message}:
				return
			case <-ctx.Done():
				return
			case <-shutdownChan:
				return
			}
		}(ctx, consumer.GetMessageChan(), consumer.ShutdownChan())
	}

	return nil
}

// Publish Producer can Publish Message to topic and it should be delivered to all consumers
// SYNC should wait for message delivery to every consumer (also with validation)
// we only guarantee that consumers that subscribed before call to that RPC will receive that message.
func (p *Processor) PublishSync(ctx context.Context, topic string, message any) error {
	consumers := p.ConsumerMatcher.MatchTopic(topic)

	wg := sync.WaitGroup{}
	wg.Add(len(consumers))

	for _, consumer := range consumers {
		// TODO|FIXME: that goroutines could probably leak in case of client stop responding, etc
		go func(doneChan <-chan struct{}, in chan ConsumerMessage, shutdownChan chan struct{}) {
			done := make(chan struct{})

			once := in

			defer wg.Done()

			for {
				select {
				// send only once no need to join here after
				case once <- ConsumerMessage{Message: message, Done: done}:
					once = nil
				// done is message delivered (with or without validation)
				case <-done:
					return
				// context done
				case <-doneChan:
					return
				case <-shutdownChan:
					return
				}
			}
		}(ctx.Done(), consumer.GetMessageChan(), consumer.ShutdownChan())
	}

	wg.Wait()

	return ctx.Err()
}

func (p *Processor) Subscribe(ctx context.Context, topic string) (Consumer, error) {
	consumer := p.consumerFactory.NewConsumer()
	p.ConsumerMatcher.AddConsumers(topic, consumer)

	return consumer, nil
}
