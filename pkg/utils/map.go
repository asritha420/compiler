package utils

type Hashable interface {
	Hash() int
}

type Comparable[K any] interface {
	Equal(K) bool
}

/*
Interface that must be implemented by the key of a map.
*/
type Mappable[K any] interface {
	Hashable
	Comparable[K]
}

type mapItem[K Mappable[K], V any] struct {
	key   K
	value V
	next  *mapItem[K, V]
}

/*
Simple hash map that can have any type that implements Mappable as a key.
*/
type Map[K Mappable[K], V any] struct {
	backing    []*mapItem[K, V]
	loadFactor float32
	len        int
}

func NewMap[K Mappable[K], V any]() *Map[K, V] {
	return &Map[K, V]{
		backing:    make([]*mapItem[K, V], 8),
		loadFactor: 2,
		len:        0,
	}
}

func (m *Map[K, V]) findPair(key K) **mapItem[K, V] {
	idx := key.Hash() % len(m.backing)
	curr := &m.backing[idx]
	for !(*curr == nil || (*curr).key.Equal(key)) {
		curr = &(*curr).next
	}
	return curr
}

func (m *Map[K, V]) resize() {
	keys, vals := m.GetAll()

	m.backing = make([]*mapItem[K, V], len(m.backing)*2)
	m.len = 0

	m.PutAll(keys, vals)
}

func (m *Map[K, V]) Len() int {
	return m.len
}

/*
Puts or replaces a value with the associated key. This will call resize at the beginning iff the load factor has been reached (even if it is just modifying an existing key).
*/
func (m *Map[K, V]) Put(key K, value V) {
	if float32(len(m.backing))*m.loadFactor <= float32(m.len) {
		m.resize()
	}

	pair := m.findPair(key)

	if *pair == nil {
		m.len++
		*pair = &mapItem[K, V]{key: key, value: value}
	}

	(*pair).value = value
}

func (m *Map[K, V]) PutAll(keys []K, vals []V) bool {
	if len(keys) != len(vals) {
		return false
	}

	for i, key := range keys {
		m.Put(key, vals[i])
	}

	return true
}

func (m Map[K, V]) Get(key K) (*V, bool) {
	pair := m.findPair(key)
	if pair == nil || *pair == nil {
		return nil, false
	}

	return &(*pair).value, true
}

func (m Map[K, V]) GetAllKeys() []K {
	keys := make([]K, m.len)

	i := 0
	for _, pair := range m.backing {
		for pair != nil {
			keys[i] = pair.key
			pair = pair.next
			i++
		}
	}

	return keys
}

func (m Map[K, V]) GetKeysHash() int {
	sum := 0
	for _, key := range m.GetAllKeys() {
		sum += key.Hash()
	}
	return sum
}

func (m Map[K, V]) GetAllVals() []V {
	vals := make([]V, m.len)

	i := 0
	for _, pair := range m.backing {
		for pair != nil {
			vals[i] = pair.value
			pair = pair.next
			i++
		}
	}

	return vals
}

func (m Map[K, V]) GetAll() ([]K, []V) {
	keys := make([]K, m.len)
	vals := make([]V, m.len)

	i := 0
	for _, pair := range m.backing {
		for pair != nil {
			keys[i] = pair.key
			vals[i] = pair.value
			pair = pair.next
			i++
		}
	}

	return keys, vals
}

func (m *Map[K, V]) Remove(key K) (*V, bool) {
	pair := m.findPair(key)
	if *pair == nil {
		return nil, false
	}

	val := &(*pair).value

	*pair = (*pair).next
	m.len--

	return val, true
}

func (m Map[K, V]) KeysEqual(other Map[K, V]) bool {
	if m.len != other.len {
		return false
	}

	for _, key := range m.GetAllKeys() {
		if _, ok := other.Get(key); !ok {
			return false
		}
	}

	return true
}
