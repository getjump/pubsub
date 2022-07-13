package internal

import "sync"

const (
	WILDCARD      = '*'
	QUESTION_MARK = '?'
)

type TrieNode struct {
	lock     sync.RWMutex
	children map[rune]*TrieNode

	consumers []Consumer
	char      rune
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
	}
}

func (n *TrieNode) Insert(topic string, consumers []Consumer) {
	node := n

	for _, r := range topic {
		node.lock.Lock()
		if _, ok := node.children[r]; !ok {
			node.children[r] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node.lock.Unlock()

		on := node

		on.lock.RLock()
		node = node.children[r]
		on.lock.RUnlock()
	}

	node.lock.Lock()
	node.consumers = consumers
	node.lock.Unlock()
}

// testTopicForPattern returns true if pattern matches topic
// Interestingly there are a leetcode problem for this https://leetcode.com/problems/wildcard-matching/
func TestTopicForPattern(topic string, pattern string) bool {
	n := len(pattern)
	m := len(topic)

	i := 0
	j := 0

	wildcardPosition := -1
	afterWildcardPosition := -1

	for i < m {
		switch {
		case j < n && ([]rune(topic)[i] == []rune(pattern)[j] || []rune(pattern)[j] == QUESTION_MARK):
			i++
			j++
		case j < n && ([]rune(pattern)[j] == WILDCARD):
			wildcardPosition = j
			afterWildcardPosition = i
			j++
		case wildcardPosition != -1:
			j = wildcardPosition + 1
			afterWildcardPosition++
			i = afterWildcardPosition
		default:
			return false
		}
	}

	for j < n && []rune(pattern)[j] == WILDCARD {
		j++
	}

	return j == n
}

// Find matches of pattern to topic for Consumer
// TODO: Optimize simple cases, like *, or *..., or ...*
// FIXME: Arg should be topic and we should recover consumer pattern from trie.
func (n *TrieNode) Find(pattern string) []Consumer {
	type llt struct {
		node  *TrieNode
		count int
		topic string
		next  *llt

		wildcardPosition int
	}

	stack := &llt{
		count:            0,
		topic:            "",
		node:             n,
		wildcardPosition: -1,
	}

	tail := stack
	appendToTail := func(node *TrieNode, count int, topic string, wildcardPosition int) {
		tail.next = &llt{
			node:             node,
			count:            count,
			topic:            topic,
			wildcardPosition: wildcardPosition,
		}
		tail = tail.next
	}

	var result []Consumer

	for stack != nil {
		node, count, topic, wildcardPosition := stack.node, stack.count, stack.topic, stack.wildcardPosition
		node.lock.RLock()

		// TODO: Probably no need to TestTopicForPattern, can we do better?
		if len(node.children) == 0 && node.consumers != nil && TestTopicForPattern(topic, pattern) {
			result = append(result, node.consumers...)
			stack = stack.next

			node.lock.RUnlock()

			continue
		}

		// TODO: greedy match heuristic, probably something like TestTopicForPattern analogue will do?
		// example: (a.*.b) <- discard paths where doesn't match latest .b
		// TODO: simplify
		var r rune

		if count < len(pattern) {
			r = []rune(pattern)[count]
		}

		if cn, ok := node.children[r]; ok && count < len(pattern) {
			appendToTail(cn, count+1, topic+string(r), wildcardPosition)
		} else if (r == QUESTION_MARK || r == WILDCARD) && count < len(pattern) {
			for cr, cn := range node.children {
				appendToTail(cn, count+1, topic+string(cr), wildcardPosition)
				if r == WILDCARD {
					tail.wildcardPosition = count
				}
			}
		} else if wildcardPosition != -1 {
			for cr, cn := range node.children {
				appendToTail(cn, wildcardPosition+1, topic+string(cr), wildcardPosition)
			}
		}

		node.lock.RUnlock()

		stack = stack.next
	}

	return result
}
