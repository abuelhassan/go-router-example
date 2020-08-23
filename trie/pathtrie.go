package trie

// Trier is an interface for a trie.
type Trier interface {
	Get(key string) interface{}
	Put(key string, value interface{})
}

// PathTrie is an implementation of Trier.
// A node's key is a path e.g. "/a/b/c".
type PathTrie struct {
	value    interface{}
	children map[string]*PathTrie
}

// NewPathTrie returns a new instance of PathTrie
func NewPathTrie() Trier {
	return &PathTrie{}
}

// Get
// Receives key in a path format e.g. "/a/b/c".
// Returns the value stored in this path.
func (t *PathTrie) Get(key string) interface{} {
	n := t
	for st := 0; st != -1; {
		seg, nxt := nextPathSegment(key, st)
		n = n.children[seg]
		if n == nil {
			return nil
		}
		st = nxt
	}
	return n.value
}

// Put
// Receives key in a path format e.g. "/a/b/c" and value
// Sets value to the given key.
func (t *PathTrie) Put(key string, value interface{}) {
	n := t
	for st := 0; st != -1; {
		if n.children == nil {
			n.children = map[string]*PathTrie{}
		}
		seg, nxt := nextPathSegment(key, st)
		ch := n.children[seg]
		if ch == nil {
			ch = &PathTrie{}
			n.children[seg] = ch
		}
		n, st = ch, nxt
	}
	n.value = value
}

// nextPathSegment
// Receives path and start index of the current segment.
// Returns the current segment and index of the next segment.
func nextPathSegment(p string, st int) (string, int) {
	const sep = '/'

	if st < 0 || st > len(p) {
		return "", -1
	}

	for i := st + 1; i < len(p); i++ {
		if p[i] == sep {
			return p[st:i], i
		}
	}
	return p[st:], -1
}
