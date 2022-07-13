package internal_test

import (
	"reflect"
	"testing"

	. "github.com/getjump/pubsub/pubsub-server/internal"
	mock_internal "github.com/getjump/pubsub/pubsub-server/internal/mock"
	"github.com/golang/mock/gomock"
)

func Test_trieMatcher_AddConsumers(t *testing.T) {
	type fields struct {
		trieRoot *TrieNode
	}
	type args struct {
		topic     string
		consumers []Consumer
	}

	ctrl := gomock.NewController(t)
	consumers := []Consumer{mock_internal.NewMockConsumer(ctrl)}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"TrieMatcher should call Insert on TrieNode", fields{NewTrieNode()}, args{"test", consumers}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			matcher := NewTrieMatcher(tt.fields.trieRoot)
			matcher.AddConsumers(tt.args.topic, tt.args.consumers...)

			if got := tt.fields.trieRoot.Find("test"); !reflect.DeepEqual(got, consumers) {
				t.Errorf("trieMatcher.AddConsumers() = %v, want %v", got, consumers)
			}
		})
	}
}

func Test_trieMatcher_MatchTopic(t *testing.T) {
	type fields struct {
		trieRoot *TrieNode
	}
	type args struct {
		topic string
	}
	ctrl := gomock.NewController(t)
	consumers := []Consumer{mock_internal.NewMockConsumer(ctrl)}

	trie := NewTrieNode()
	trie.Insert("test", consumers)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Consumer
	}{
		{"MatchTopic should call Find on trieRoot", fields{trie}, args{"test"}, consumers},
		{"MatchTopic should call Find on trieRoot", fields{trie}, args{"test2"}, nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			matcher := NewTrieMatcher(tt.fields.trieRoot)
			if got := matcher.MatchTopic(tt.args.topic); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trieMatcher.MatchTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTrieMatcher(t *testing.T) {
	t.Run("NewTrieMatcher() returns trieMatcher", func(t *testing.T) {
		t.Parallel()
		if got := NewTrieMatcher(NewTrieNode()); reflect.TypeOf(got).String() == "" {
			t.Errorf("NewTrieMatcher() = %v", got)
		}
	})
}
