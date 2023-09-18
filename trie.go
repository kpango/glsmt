package glsmt

import (
	"sync"

	"github.com/dolthub/swiss"
)

type Node[V any] struct {
	parent      *Node[V]
	children    [256]*Node[V]
	extChildren *swiss.Map[rune, *Node[V]]
	value       *V
	mu          sync.RWMutex
}

type Trie[V any] interface {
	Insert(key string, value *V) bool
	Get(key string) (*V, bool)
	Delete(key string) (*V, bool)
}

type trie[V any] struct {
	root  *Node[V]
	nSize int
}

func NewTrie[V any](size int) Trie[V] {
	return &trie[V]{
		root:  newNode[V](),
		nSize: size,
	}
}

func newNode[V any]() *Node[V] {
	return &Node[V]{}
}

func (n *Node[V]) getChild(ch rune) (node *Node[V], ok bool) {
	if ch < 256 {
		n.mu.RLock()
		node = n.children[ch]
		n.mu.RUnlock()
		return node, node != nil
	}
	n.mu.RLock()
	node, ok = n.extChildren.Get(ch)
	n.mu.RUnlock()
	return node, ok && node != nil
}

func (n *Node[V]) setChild(ch rune, child *Node[V], size int) (node *Node[V], ok bool) {
	if ch < 256 {
		n.mu.Lock()
		node = n.children[ch]
		if node == nil {
			n.children[ch] = child
			node = child
			ok = true
		}
		n.mu.Unlock()
	} else {
		n.mu.Lock()
		if n.extChildren == nil {
			n.extChildren = swiss.NewMap[rune, *Node[V]](16)
		}
		node, ok = n.extChildren.Get(ch)
		if !ok && node == nil {
			n.extChildren.Put(ch, child)
			node = child
			ok = true
		}
		n.mu.Unlock()
	}
	return
}

func (t *trie[V]) Insert(key string, value *V) bool {
	return t.traverse(key, func(node *Node[V], idx int) (ok bool) {
		if idx < len(key)-1 {
			var next *Node[V]
			for i, ch := range key[idx:] {
				next, ok = node.getChild(ch)
				if !ok || next == nil {
					next, _ = node.setChild(ch, newNode[V](), t.nSize)
					for _, c := range key[idx+i+1:] {
						next, _ = next.setChild(c, newNode[V](), t.nSize)
					}
					next.value = value
					return true
				}
				node = next
			}
		}
		node.value = value
		return true
	})
}

func (t *trie[V]) Get(key string) (v *V, ok bool) {
	return v, t.traverse(key, func(node *Node[V], idx int) (ok bool) {
		if idx == len(key)-1 {
			node.mu.RLock()
			v = node.value
			node.mu.RUnlock()
			ok = v != nil
		}
		return ok
	})
}

func (t *trie[V]) Delete(key string) (v *V, ok bool) {
	return v, t.traverse(key, func(node *Node[V], idx int) (ok bool) {
		if idx == len(key)-1 {
			node.mu.Lock()
			v = node.value
			node.value = nil
			node.mu.Unlock()
			ok = v != nil
		}
		return ok
	})
}

func (t *trie[V]) traverse(key string, fn func(node *Node[V], idx int) bool) (ok bool) {
	node := t.root
	var (
		i    int
		ch   rune
		next *Node[V]
	)
	for i, ch = range key {
		next, ok = node.getChild(ch)
		if !ok || next == nil {
			return fn(node, i)
		}
		node = next
	}
	return fn(node, i)
}
