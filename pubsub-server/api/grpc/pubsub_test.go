package grpc

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/getjump/pubsub/pubsub-server/internal"
	"github.com/getjump/pubsub/pubsub-server/pkg/proto"
)

func TestPubSubServer_Publish(t *testing.T) {
	type fields struct {
		UnimplementedPubSubServer proto.UnimplementedPubSubServer
		Processor                 internal.Processor
		messagesToValidate        map[clientIDType]*clientMessageChanMap
		messagesToValidateLock    sync.RWMutex
		retrySendTimeout          time.Duration
		maxRetries                uint
	}
	type args struct {
		ctx context.Context
		req *proto.PublishRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto.PublishResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &PubSubServer{
				UnimplementedPubSubServer: tt.fields.UnimplementedPubSubServer,
				Processor:                 tt.fields.Processor,
				messagesToValidate:        tt.fields.messagesToValidate,
				messagesToValidateLock:    tt.fields.messagesToValidateLock,
				retrySendTimeout:          tt.fields.retrySendTimeout,
				maxRetries:                tt.fields.maxRetries,
			}
			got, err := server.Publish(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PubSubServer.Publish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PubSubServer.Publish() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPubSubServer_retrySendIfNotRetrieved(t *testing.T) {
	type fields struct {
		UnimplementedPubSubServer proto.UnimplementedPubSubServer
		Processor                 internal.Processor
		messagesToValidate        map[clientIDType]*clientMessageChanMap
		messagesToValidateLock    sync.RWMutex
		retrySendTimeout          time.Duration
		maxRetries                uint
	}
	type args struct {
		clientId         clientIDType
		message          *proto.PubSubMessage
		doneChan         chan struct{}
		retryMessageChan chan *retryMessage
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &PubSubServer{
				UnimplementedPubSubServer: tt.fields.UnimplementedPubSubServer,
				Processor:                 tt.fields.Processor,
				messagesToValidate:        tt.fields.messagesToValidate,
				messagesToValidateLock:    tt.fields.messagesToValidateLock,
				retrySendTimeout:          tt.fields.retrySendTimeout,
				maxRetries:                tt.fields.maxRetries,
			}
			server.retrySendIfNotRetrieved(tt.args.clientId, tt.args.message, tt.args.doneChan, tt.args.retryMessageChan)
		})
	}
}

func TestPubSubServer_Subscribe(t *testing.T) {
	type fields struct {
		UnimplementedPubSubServer proto.UnimplementedPubSubServer
		Processor                 internal.Processor
		messagesToValidate        map[clientIDType]*clientMessageChanMap
		messagesToValidateLock    sync.RWMutex
		retrySendTimeout          time.Duration
		maxRetries                uint
	}
	type args struct {
		req    *proto.SubscribeRequest
		stream proto.PubSub_SubscribeServer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &PubSubServer{
				UnimplementedPubSubServer: tt.fields.UnimplementedPubSubServer,
				Processor:                 tt.fields.Processor,
				messagesToValidate:        tt.fields.messagesToValidate,
				messagesToValidateLock:    tt.fields.messagesToValidateLock,
				retrySendTimeout:          tt.fields.retrySendTimeout,
				maxRetries:                tt.fields.maxRetries,
			}
			if err := server.Subscribe(tt.args.req, tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("PubSubServer.Subscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPubSubServer_ValidateRetrieval(t *testing.T) {
	type fields struct {
		UnimplementedPubSubServer proto.UnimplementedPubSubServer
		Processor                 internal.Processor
		messagesToValidate        map[clientIDType]*clientMessageChanMap
		messagesToValidateLock    sync.RWMutex
		retrySendTimeout          time.Duration
		maxRetries                uint
	}
	type args struct {
		stream proto.PubSub_ValidateRetrievalServer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &PubSubServer{
				UnimplementedPubSubServer: tt.fields.UnimplementedPubSubServer,
				Processor:                 tt.fields.Processor,
				messagesToValidate:        tt.fields.messagesToValidate,
				messagesToValidateLock:    tt.fields.messagesToValidateLock,
				retrySendTimeout:          tt.fields.retrySendTimeout,
				maxRetries:                tt.fields.maxRetries,
			}
			if err := server.ValidateRetrieval(tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("PubSubServer.ValidateRetrieval() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
