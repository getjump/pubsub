package internal

import (
	"context"
	"log"
	"sync"

	"github.com/emirpasic/gods/trees/redblacktree"
)

// SimpleStringConsumerMatcher accepts simple string topic and matches that to consumers
// TODO: Implement custom Matcher via Trie/... to accept wildcard everywhere in topic
type SimpleStringConsumerMatcher struct {
	lock            sync.RWMutex
	topicToConsumer map[string]*redblacktree.Tree

	consumerToTopic map[Consumer][]string
}

func NewSimpleStringConsumerMatcher(ctx context.Context) *SimpleStringConsumerMatcher {
	matcher := SimpleStringConsumerMatcher{
		topicToConsumer: make(map[string]*redblacktree.Tree),
		consumerToTopic: make(map[Consumer][]string),
	}

	return &matcher
}

func (matcher *SimpleStringConsumerMatcher) MatchTopic(topic string) []Consumer {
	matcher.lock.RLock()
	defer matcher.lock.RUnlock()

	query, ok := matcher.topicToConsumer[topic]

	if !ok {
		return nil
	}

	var result []Consumer

	iterator := query.Iterator()

	for iterator.Next() {
		c, ok := (iterator.Key()).(Consumer)

		if !ok {
			log.Fatalf("wrong type")

			continue
		}

		result = append(result, c)
	}

	return result
}

func consumerComparator(a, b interface{}) int {
	c1 := a.(Consumer)
	c2 := b.(Consumer)

	if c1.GetID() > c2.GetID() {
		return 1
	} else if c1.GetID() < c2.GetID() {
		return -1
	}

	return 0
}

func (matcher *SimpleStringConsumerMatcher) AddConsumers(topic string, consumers ...Consumer) {
	matcher.lock.Lock()
	defer matcher.lock.Unlock()

	tree := matcher.topicToConsumer[topic]

	if tree == nil {
		matcher.topicToConsumer[topic] = redblacktree.NewWith(consumerComparator)
		tree = matcher.topicToConsumer[topic]
	}

	for _, consumer := range consumers {
		matcher.consumerToTopic[consumer] = append(matcher.consumerToTopic[consumer], topic)
		tree.Put(consumer, struct{}{})
	}
}

func (matcher *SimpleStringConsumerMatcher) RemoveConsumerFromTopic(topic string, consumer Consumer) {
	matcher.lock.Lock()
	defer matcher.lock.Unlock()

	delete(matcher.consumerToTopic, consumer)

	tree := matcher.topicToConsumer[topic]
	tree.Remove(consumer)
}
