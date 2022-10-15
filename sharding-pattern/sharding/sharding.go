package sharding

import (
	"crypto/sha1"
	"log"
	"sync"
)

type Shard[T any] struct {
	sync.RWMutex
	m map[string]T
}

type ShardedMap[T any] []*Shard[T]

func NewShardedMap[T any](nshards int) ShardedMap[T] {
	shards := make([]*Shard[T], nshards)

	for i := 0; i < nshards; i++ {
		shard := make(map[string]T)
		shards[i] = &Shard[T]{m: shard}
	}

	return shards
}

func (m ShardedMap[T]) getShardIndex(key string) int {
	checksum := sha1.Sum([]byte(key))
	log.Printf("Checksum: %v\n", checksum)
	hash := int(checksum[13])<<8 | int(checksum[17])
	log.Printf("hash: %v\n", hash)
	return hash % len(m)
}

func (m ShardedMap[T]) getShard(key string) *Shard[T] {
	index := m.getShardIndex(key)
	return m[index]
}

func (m ShardedMap[T]) Get(key string) T {
	shard := m.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	return shard.m[key]
}

func (m ShardedMap[T]) Set(key string, value T) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	shard.m[key] = value
}

func (m ShardedMap[T]) Keys() []string {
	keys := make([]string, 0)

	mutex := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(len(m))

	for _, shard := range m {
		go func(s *Shard[T]) {
			s.RLock()

			for key := range s.m {
				mutex.Lock()
				keys = append(keys, key)
				mutex.Unlock()
			}

			s.RUnlock()
			wg.Done()
		}(shard)
	}

	wg.Wait()

	return keys
}
