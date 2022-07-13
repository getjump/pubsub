package internal_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	. "github.com/getjump/pubsub/pubsub-server/internal"
	mock_internal "github.com/getjump/pubsub/pubsub-server/internal/mock"
	"github.com/golang/mock/gomock"
)

func TestProcessor_PublishAsync(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockConsumer := mock_internal.NewMockConsumer(ctrl)
	consumers := []Consumer{mockConsumer}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	t.Run("PublishAsync should send on consumer channel", func(t *testing.T) {
		p := NewProcessor(NewTrieMatcher(NewTrieNode()), mock_internal.NewMockConsumerFactory(ctrl))

		p.AddConsumers("test", consumers...)

		testChan := make(chan ConsumerMessage)

		mockConsumer.EXPECT().GetMessageChan().Times(1).Return(testChan)

		if err := p.PublishAsync(ctx, "test", nil); err != nil {
			t.Errorf("Processor.PublishAsync() error = %v", err)
		}

		select {
		case <-testChan:
		case <-time.After(time.Millisecond * 200):
			t.Errorf("There were no send on test channel")
			return
		}
	})

	t.Run("PublishAsync should not leak goroutines", func(t *testing.T) {
		p := NewProcessor(NewTrieMatcher(NewTrieNode()), mock_internal.NewMockConsumerFactory(ctrl))

		p.AddConsumers("test", consumers...)

		testChan := make(chan ConsumerMessage)

		mockConsumer.EXPECT().GetMessageChan().Times(1).Return(testChan)

		ctx, close := context.WithTimeout(context.Background(), time.Second*5)

		if err := p.PublishAsync(ctx, "test", nil); err != nil {
			t.Errorf("Processor.PublishAsync() error = %v", err)
		}

		close()
	})
}

func TestProcessor_PublishSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockConsumer := mock_internal.NewMockConsumer(ctrl)
	consumers := []Consumer{mockConsumer}

	ctx := context.Background()

	t.Run("PublishSync should block until receive", func(t *testing.T) {
		p := NewProcessor(NewTrieMatcher(NewTrieNode()), mock_internal.NewMockConsumerFactory(ctrl))

		p.AddConsumers("test", consumers...)

		testChan := make(chan ConsumerMessage)
		mockConsumer.EXPECT().GetMessageChan().Times(1).Return(testChan)

		isBlocked := false

		go func() {
			<-time.After(time.Millisecond * 10)
			m := <-testChan
			isBlocked = true
			m.Done <- struct{}{}
		}()

		if err := p.PublishSync(ctx, "test", nil); err != nil {
			t.Errorf("Processor.PublishSync() error = %v", err)
		}

		if !isBlocked {
			t.Errorf("Processor.PublishSync() does not block")
		}
	})

	t.Run("PublishSync leaves by context signal", func(t *testing.T) {
		p := NewProcessor(NewTrieMatcher(NewTrieNode()), mock_internal.NewMockConsumerFactory(ctrl))

		p.AddConsumers("test", consumers...)

		testChan := make(chan ConsumerMessage)
		mockConsumer.EXPECT().GetMessageChan().Times(1).Return(testChan)

		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)

		testCh := make(chan struct{})

		go func() {
			<-ctx.Done()

			select {
			case <-testCh:
				return
			case <-time.After(time.Millisecond):
				t.Errorf("PublishSync still blocked")
			}
		}()

		if err := p.PublishSync(ctx, "test", nil); !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Processor.PublishSync() error = %v", err)
		}

		testCh <- struct{}{}
	})
}

func TestProcessor_Subscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockConsumer := mock_internal.NewMockConsumer(ctrl)
	ctx := context.Background()

	t.Run("Subscribe should return channel", func(t *testing.T) {
		factory := mock_internal.NewMockConsumerFactory(ctrl)
		p := NewProcessor(NewTrieMatcher(NewTrieNode()), factory)

		testChan := make(chan ConsumerMessage)

		factory.EXPECT().NewConsumer().Times(1).Return(mockConsumer)
		mockConsumer.EXPECT().GetMessageChan().Times(1).Return(testChan)

		got, _, err := p.Subscribe(ctx, "test")
		if err != nil {
			t.Errorf("Processor.Subscribe() error = %v", err)

			return
		}

		go func() {
			if !reflect.DeepEqual((<-got).Message, nil) {
				t.Errorf("Processor.Subscribe() = %v", got)
			}
		}()

		mockConsumer.EXPECT().GetMessageChan().Times(1).Return(testChan)
		p.PublishAsync(ctx, "test", nil)
	})
}
