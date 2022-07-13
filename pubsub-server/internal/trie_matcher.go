package internal

type trieMatcher struct {
	trieRoot *TrieNode
}

func NewTrieMatcher(trieRoot *TrieNode) *trieMatcher {
	return &trieMatcher{
		trieRoot: trieRoot,
	}
}

// TODO: Something here is completely broken need to rework a bit
func (matcher *trieMatcher) MatchTopic(topic string) []Consumer {
	return matcher.trieRoot.Find(topic)
}

func (matcher *trieMatcher) AddConsumers(topic string, consumers ...Consumer) {
	matcher.trieRoot.Insert(topic, consumers)
}
