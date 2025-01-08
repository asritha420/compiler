package main

type Mappable interface {
	Hash() int
	Equal(Mappable) bool
}

type KVPair[K Mappable, V any] struct {
	key K
	value V
}

type Map[K Mappable, V any] struct {
	backing []*KVPair[K, V]
	loadFactor float32
	len int
}

func NewMap[K Mappable, V any]() *Map[K,V] {
	return &Map[K, V]{
		backing: make([]*KVPair[K,V], 8),
		loadFactor: 0.75,
		len: 0,
	}
}

func (m *Map[K, V]) findPair(key K) **KVPair[K,V] {
	backingLen := len(m.backing)
	startIdx := key.Hash() % backingLen
	currIdx := startIdx
	var found **KVPair[K,V] = nil
	for {
		kvPair := m.backing[currIdx]
		if kvPair == nil || kvPair.key.Equal(key){
			found = &m.backing[currIdx]
			break
		}
		currIdx = (currIdx + 1) % backingLen
		if currIdx == startIdx {
			break
		}
	}
	return found
}

func (m *Map[K, V]) resize() {
	pairs := m.getKVPairs()

	m.backing = make([]*KVPair[K,V], len(m.backing)*2)
	m.len = 0

	for _, pair := range pairs {
		m.Put(pair.key, pair.value)	
	}
}

func (m *Map[K, V]) getKVPairs() []*KVPair[K, V] {
	pairs := make([]*KVPair[K,V], m.len)
	i := 0
	for _, pair := range m.backing {
		if pair != nil {
			pairs[i] = pair
			i++
		}
	}
	
	return pairs
}

func (m *Map[K, V]) Put(key K, value V) {
	if float32(len(m.backing)) * m.loadFactor <= float32(m.len) {
		m.resize()
	}

	pair := m.findPair(key)

	if *pair == nil {
		m.len++
		*pair = &KVPair[K, V]{key: key, value: value} 
	}

	(*pair).value = value
}

func (m *Map[K, V]) Get(key K) (*V, bool) {
	pair := m.findPair(key)
	if pair == nil || *pair == nil {
		return nil, false
	}

	return &(*pair).value, true
}

func (m *Map[K, V]) Remove(key K) (*V, bool) {
	pair := m.findPair(key)
	if pair == nil || *pair == nil {
		return nil, false
	}

	val := &(*pair).value

	*pair = nil
	m.len--

	return val, true
}