package internal_test

import (
	"context"
	"reflect"
	"testing"

	. "github.com/getjump/pubsub/pubsub-server/internal"
	mock_internal "github.com/getjump/pubsub/pubsub-server/internal/mock"
	"github.com/golang/mock/gomock"
)

func TestSimpleStringConsumerMatcher_MatchTopic(t *testing.T) {
	type args struct {
		topic string
	}

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	consumers := []Consumer{mock_internal.NewMockConsumer(ctrl)}

	tests := []struct {
		name string
		args args
		want []Consumer
	}{
		{"Simple Matcher matches string", args{"test"}, consumers},
		{"Simple Matcher doesn't match string", args{"test2"}, nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			matcher := NewSimpleStringConsumerMatcher(ctx)
			matcher.AddConsumers(tt.args.topic, tt.want...)
			if got := matcher.MatchTopic(tt.args.topic); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleStringConsumerMatcher.MatchTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleStringConsumerMatcher_AddConsumers(t *testing.T) {
	type args struct {
		topic     string
		consumers []Consumer
	}

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	mockConsumer := mock_internal.NewMockConsumer(ctrl)
	consumers := []Consumer{mockConsumer}

	tests := []struct {
		name string
		args args
	}{
		{"MatchTopic returns added consumers", args{"test", consumers}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockConsumer.EXPECT().GetID().AnyTimes().Return(uint64(1))
			matcher := NewSimpleStringConsumerMatcher(ctx)
			matcher.AddConsumers(tt.args.topic, tt.args.consumers...)
			if got := matcher.MatchTopic(tt.args.topic); !reflect.DeepEqual(got, tt.args.consumers) {
				t.Errorf("SimpleStringConsumerMatcher.MatchTopic() = %v, want %v", got, tt.args.consumers)
			}
		})
	}
}
