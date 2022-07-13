package grpc

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/getjump/pubsub/pubsub-server/internal"
	"github.com/getjump/pubsub/pubsub-server/pkg/proto"
)

type (
	clientIDType  string
	messageIDType string
)

type clientMessageChanMap struct {
	messages map[messageIDType]chan struct{}
	Lock     sync.RWMutex
}

type retryMessage struct {
	*proto.PubSubMessage
	retries uint
}

func NewPubSubServer(ctx context.Context, retrySendTimeout time.Duration, maxRetries uint) *PubSubServer {
	server := PubSubServer{
		// TODO: that map can grow pretty fast rework or cleanup?
		messagesToValidate: make(map[clientIDType]*clientMessageChanMap),
		retrySendTimeout:   retrySendTimeout,
		maxRetries:         maxRetries,
		Processor: internal.NewProcessor(
			internal.NewSimpleStringConsumerMatcher(ctx), &internal.SimpleConsumerFactory{},
		),
	}

	return &server
}

type PubSubServer struct {
	proto.UnimplementedPubSubServer

	*internal.Processor

	// TODO: that map can grow pretty fast rework or cleanup?
	messagesToValidate     map[clientIDType]*clientMessageChanMap
	messagesToValidateLock sync.RWMutex

	retrySendTimeout time.Duration
	maxRetries       uint
}

// Publish Producer can Publish Message to topic and it should be delivered to all consumers
// TODO: probably sync and async modes should be accepted
// SYNC should wait for message delivery to every consumer (also with validation)
// ASYNC should return result after processing and spawn goroutines to send in parallel.
func (server *PubSubServer) Publish(ctx context.Context, req *proto.PublishRequest) (*proto.PublishResponse, error) {
	response := &proto.PublishResponse{}
	response.PublishStatus = proto.PublishResponse_FAILED

	var err error

	switch req.PublishMethod {
	case proto.PublishRequest_ASYNC:
		// TODO: Once handler leaves context is done, do better?
		err = server.Processor.PublishAsync(context.Background(), req.Message.Topic, req.Message)
	case proto.PublishRequest_SYNC:
		// TODO: sync send should probably have some timeout
		// ctx, _ := context.WithTimeout(ctx, time.Millisecond*500)
		err = server.Processor.PublishSync(ctx, req.Message.Topic, req.Message)
	}

	if err != nil {
		return nil, err
	}

	response.PublishStatus = proto.PublishResponse_OK

	return response, nil
}

func (server *PubSubServer) retrySendIfNotRetrieved(
	clientId clientIDType,
	message *proto.PubSubMessage,
	doneChan chan struct{},
	consumer internal.Consumer,
	retryMessageChan chan *retryMessage,
) {
	valChan := make(chan struct{})

	server.messagesToValidateLock.RLock()
	clientPool := server.messagesToValidate[clientIDType(clientId)]
	server.messagesToValidateLock.RUnlock()

	clientPool.Lock.Lock()
	clientPool.messages[messageIDType(message.Id)] = valChan
	clientPool.Lock.Unlock()

	retryMessage := &retryMessage{message, 0}

	for {
		select {
		case <-valChan:
			if doneChan != nil {
				clientPool.Lock.Lock()
				delete(clientPool.messages, messageIDType(message.Id))
				clientPool.Lock.Unlock()
				// doneChan <- struct{}{}
				close(doneChan)
			}

			return
		case <-time.After(server.retrySendTimeout):
			// TODO: probably we should also do something with that consumer? like delete it completely
			// because it doesn't respond anymore (gRPC have ping/pong internally)
			if retryMessage.retries >= server.maxRetries {
				server.ConsumerMatcher.RemoveConsumerFromTopic(message.Topic, consumer)
				close(consumer.ShutdownChan())

				return
			}

			retryMessageChan <- retryMessage
			retryMessage.retries += 1
		}
	}
}

func (server *PubSubServer) Subscribe(req *proto.SubscribeRequest, stream proto.PubSub_SubscribeServer) error {
	consumer, err := server.Processor.Subscribe(stream.Context(), req.Topic)
	if err != nil {
		return err
	}

	server.messagesToValidateLock.Lock()
	server.messagesToValidate[clientIDType(req.ClientId)] = &clientMessageChanMap{
		messages: make(map[messageIDType]chan struct{}),
	}
	server.messagesToValidateLock.Unlock()

	// Should be buffered?
	retryChan := make(chan *retryMessage)

	// t := uint64(0)

	for {
		select {
		case <-stream.Context().Done():
			close(consumer.ShutdownChan())
			server.messagesToValidateLock.Lock()
			delete(server.messagesToValidate, clientIDType(req.ClientId))
			server.messagesToValidateLock.Unlock()

			return stream.Context().Err()
		case m := <-consumer.GetMessageChan():
			psm, ok := m.Message.(*proto.PubSubMessage)

			if !ok {
				return fmt.Errorf("wrong type %T", m.Message)
			}

			err := stream.Send(psm)
			if err != nil {
				close(consumer.ShutdownChan())
				server.messagesToValidateLock.Lock()
				delete(server.messagesToValidate, clientIDType(req.ClientId))
				server.messagesToValidateLock.Unlock()

				server.RemoveConsumerFromTopic(req.Topic, consumer)

				return err
			}

			if m.Done != nil {
				// TODO|FIXME: that goroutines could probably leak in case of client stop responding, etc
				go func(doneChan chan struct{}, validationRequired *bool) {
					// In other case server.retrySendIfNotRetrieved will send once done
					if validationRequired == nil || (validationRequired != nil && !*validationRequired) {
						// doneChan <- struct{}{}
						// atomic.AddUint64(&t, uint64(1))
						close(doneChan)
					}
				}(m.Done, psm.ValidationRequired)
			}

			// fmt.Println(atomic.LoadUint64(&t))

			if psm.ValidationRequired != nil && *psm.ValidationRequired {
				// maybe implement it in that loop/select-case?
				go server.retrySendIfNotRetrieved(clientIDType(req.ClientId), psm, m.Done, consumer, retryChan)
			}
		case m := <-retryChan:
			stream.Send(m.PubSubMessage)
		}
	}
}

// ValidateRetrieval validates retrieval of Message by client if it is required.
func (server *PubSubServer) ValidateRetrieval(stream proto.PubSub_ValidateRetrievalServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		default:
			break
		}

		message, err := stream.Recv()
		if err != nil {
			return err
		}

		server.messagesToValidateLock.RLock()
		client, ok := server.messagesToValidate[clientIDType(message.ClientId)]

		if !ok {
			server.messagesToValidateLock.RUnlock()

			return errors.New("unknown client")
		}
		server.messagesToValidateLock.RUnlock()

		client.Lock.RLock()
		if valChan, ok := client.messages[messageIDType(message.Message.Id)]; ok {
			client.Lock.RUnlock()
			valChan <- struct{}{}
		} else {
			client.Lock.RUnlock()

			// TODO: Custom errors
			return errors.New("unknown message")
		}
	}
}
