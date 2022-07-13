package pkg

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/getjump/pubsub/pubsub-server/pkg/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"

	pb "google.golang.org/protobuf/proto"
)

type PubSubClient struct {
	proto.PubSubClient

	clientId string
}

type Message struct {
	Data *anypb.Any
}

func NewPubSubClient(cc grpc.ClientConnInterface, clientId string) *PubSubClient {
	client := &PubSubClient{
		PubSubClient: proto.NewPubSubClient(cc),
		clientId:     clientId,
	}

	return client
}

// TODO: do better then pb.Message
func (client *PubSubClient) Publish(
	ctx context.Context, topic string, data pb.Message, validationRequired bool, publishMethod proto.PublishRequest_PublishMethod,
) error {
	body := &anypb.Any{}
	err := anypb.MarshalFrom(body, data, pb.MarshalOptions{})
	if err != nil {
		return err
	}

	message := &proto.PubSubMessage{
		Topic:              topic,
		Id:                 uuid.NewString(),
		ValidationRequired: &validationRequired,
		Body:               body,
	}

	response, err := client.PubSubClient.Publish(ctx, &proto.PublishRequest{
		PublishMethod: publishMethod,
		Message:       message,
	})
	if err != nil {
		return err
	}

	if response.PublishStatus == proto.PublishResponse_FAILED {
		return fmt.Errorf("publish failed")
	}

	return err
}

func (client *PubSubClient) PublishAsync(ctx context.Context, topic string, data pb.Message, validationRequired bool) error {
	return client.Publish(ctx, topic, data, validationRequired, proto.PublishRequest_ASYNC)
}

func (client *PubSubClient) PublishSync(ctx context.Context, topic string, data pb.Message, validationRequired bool) error {
	return client.Publish(ctx, topic, data, validationRequired, proto.PublishRequest_SYNC)
}

func (client *PubSubClient) Subscribe(ctx context.Context, topic string) (chan *Message, error) {
	sc, err := client.PubSubClient.Subscribe(ctx, &proto.SubscribeRequest{
		Topic:    topic,
		ClientId: client.clientId,
	})

	result := make(chan *Message)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				break
			}

			m, err := sc.Recv()
			if err == io.EOF {
				continue
			}

			// TODO error processing.
			if err != nil {
				log.Println(err)

				return
			}

			result <- &Message{m.Body}

			if m.ValidationRequired != nil && *m.ValidationRequired {
				go func() {
					vc, err := client.PubSubClient.ValidateRetrieval(ctx)
					// TODO error processing
					if err != nil {
						log.Println(err)

						return
					}

					err = vc.Send(&proto.ValidateRetrievalRequest{
						ClientId: client.clientId,
						Message:  m,
					})

					// TODO error processing
					if err != nil {
						log.Println(err)

						return
					}
				}()
			}
		}
	}()

	if err != nil {
		return nil, err
	}

	return result, nil
}
