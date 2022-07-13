package internal_test

import (
	"reflect"
	"testing"

	. "github.com/getjump/pubsub/pubsub-server/internal"
	mock_internal "github.com/getjump/pubsub/pubsub-server/internal/mock"
	"github.com/golang/mock/gomock"
)

func Test_TrieNode_Find(t *testing.T) {
	type args struct {
		topic string
	}
	type _ts struct {
		topic     string
		consumers []Consumer
	}

	ctrl := gomock.NewController(t)

	testConsumer := []Consumer{mock_internal.NewMockConsumer(ctrl)}
	test2Consumer := []Consumer{mock_internal.NewMockConsumer(ctrl), mock_internal.NewMockConsumer(ctrl)}
	tests := []struct {
		name string
		args args
		ts   []_ts
		want []Consumer
	}{
		{"Pattern equals string", args{"a.b.c"}, []_ts{
			{"a.b.c", testConsumer},
		}, testConsumer},
		{"Pattern doesn't match", args{"a.d.c"}, []_ts{
			{"a.b.c", testConsumer},
		}, nil},
		{"Question mark matches any single character", args{"a.?.c"}, []_ts{
			{"a.b.c", testConsumer},
		}, testConsumer},
		{"Wildcard matches single character", args{"a.*.c"}, []_ts{
			{"a.b.c", testConsumer},
		}, testConsumer},
		{"Wildcard matches single character but only after", args{"a.c*.c"}, []_ts{
			{"a.db.c", testConsumer},
		}, nil},
		{"Wildcard matches zero character", args{"a.b*.c"}, []_ts{
			{"a.b.c", testConsumer},
		}, testConsumer},
		{"Wildcard matches any amount of characters", args{"a.*.c"}, []_ts{
			{"a.b.c.d.c", testConsumer},
		}, testConsumer},
		{"Wildcard matches any amount of characters", args{"a.*.c"}, []_ts{
			{"a.b.c.d.c", testConsumer},
		}, testConsumer},
		{"Wildcard matches any amount of characters", args{"a*c"}, []_ts{
			{"a.d.c", testConsumer},
			{"a.b.b.c", test2Consumer},
		}, append(testConsumer, test2Consumer...)},
		{"Wildcard matches any amount of characters", args{"*"}, []_ts{
			{"a.d.c", testConsumer},
			{"a.b.b.c", test2Consumer},
		}, append(testConsumer, test2Consumer...)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := NewTrieNode()

			for _, ts := range tt.ts {
				root.Insert(ts.topic, ts.consumers)
			}

			if got := root.Find(tt.args.topic); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trieNode.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_TestTopicForPattern(t *testing.T) {
	type args struct {
		topic   string
		pattern string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Pattern equals string", args{"a.b.c", "a.b.c"}, true},
		{"Question mark matches any single character", args{"a.b.c", "a.?.c"}, true},
		{"Wildcard matches single character", args{"a.b.c", "a.*.c"}, true},
		{"Wildcard matches zero character", args{"a.b.c", "a.b*.c"}, true},
		{"Wildcard matches any amount of characters", args{"a.b.c.d.c", "a.*.c"}, true},
		{"Pattern doesn't match", args{"a.b.c", "a.d.c"}, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := TestTopicForPattern(tt.args.topic, tt.args.pattern); got != tt.want {
				t.Errorf("TestTopicForPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTrieNode(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Returns new trie node with initialized symbols map"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewTrieNode(); reflect.ValueOf(got).Elem().FieldByName("children").IsNil() {
				t.Errorf("NewTrieNode() doesn't initialize map")
			}
		})
	}
}
