// Package cache implements a cache. This cache is simple: a map with a mutex.
// There is no fancy expunge algorithm, it just randomly evicts elements when
// it gets full.
// To improve concurrency the cache is shared 256 ways. The hash key is created
// with hash.fnv and the first 8 bits are used to select the cache shared.
// TODO(...): The number 256 is randomly choosen.
package cache

import (
	"hash/fnv"
	"sync"
)

// onEvict?

// Hash return the hash of what.
func Hash(what []byte) uint32 {
	h := fnv.New32()
	h.Write(what)
	return h.Sum32()
}

const shardSize = 256

// Cache is a ...
type Cache struct {
	shards [shardSize]*shard
}

// shard is a cache with random eviction.
type shard struct {
	items map[uint32]interface{}
	size  int

	sync.RWMutex
}

func New(size int) *Cache {
	ssize := size / shardSize
	if ssize < 512 {
		ssize = 512
	}

	c := &Cache{}

	// Initialize all the shards
	for i := 0; i < shardSize; i++ {
		c.shards[i] = newShard(ssize)
	}
	return c
}

func (c *Cache) Add(key uint32, el interface{}) {
	shard := key & (shardSize - 1)
	c.shards[shard].Add(key, el)
}

func (c *Cache) Get(key uint32) (interface{}, bool) {
	shard := key & (shardSize - 1)
	return c.shards[shard].Get(key)
}

func (c *Cache) Remove(key uint32) {
	shard := key & (shardSize - 1)
	c.shards[shard].Remove(key)
}

func (c *Cache) Len() int {
	l := 0
	for _, s := range c.shards {
		l += s.Len()
	}
	return l
}

// newShard returns a new shard with the specified size.
func newShard(size int) *shard { return &shard{items: make(map[uint32]interface{}), size: size} }

// Add adds element indexed by key into the cache. Any existing element is overwritten
func (s *shard) Add(key uint32, el interface{}) {
	l := s.Len()
	if l+1 > s.size {
		s.Evict()
	}

	// Now our locking.
	s.Lock()
	defer s.Unlock()
	s.items[key] = el
}

// Remove removes the element indexed by key from the cache.
func (s *shard) Remove(key uint32) {
	s.Lock()
	defer s.Unlock()
	delete(s.items, key)
}

// Evict removes a random element from the cache.
func (s *shard) Evict() {
	s.Lock()
	defer s.Unlock()

	key := -1
	for k := range s.items {
		key = int(k)
		break
	}
	if key == -1 {
		// empty cache
		return
	}
	delete(s.items, uint32(key))
}

// Get looks up the element indexed under key.
func (s *shard) Get(key uint32) (interface{}, bool) {
	s.RLock()
	defer s.RUnlock()

	el, found := s.items[key]
	return el, found
}

// Len returns the current length of the cache.
func (s *shard) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.items)
}
