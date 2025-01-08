package utils

import (
	"testing"
)

type num struct {
	val float64
}

func (n num) Equal(other num) bool {
	return n.val == other.val
}

func (n num) Hash() int {
	return int(n.val)
}

func initMap(m *Map[num, int]) {
	m.Put(num{3}, 54)
	m.Put(num{16}, 91) //collision between 16 and 0
	m.Put(num{0}, 85)
	m.Put(num{7}, 2000)
}

func TestPut(t *testing.T) {
	simpleMap := NewMap[num, int]()
	initMap(simpleMap)
	initMap(simpleMap) //test putting twice wont do anything

	AssertEqual(t, "m.len", 4, simpleMap.len)

	AssertEqual(t, "m.backing[0].key", num{16}, simpleMap.backing[0].key)
	AssertEqual(t, "m.backing[0].val", 91, simpleMap.backing[0].value)

	AssertEqual(t, "m.backing[0].next.key", num{0}, simpleMap.backing[0].next.key)
	AssertEqual(t, "m.backing[0].next.val", 85, simpleMap.backing[0].next.value)

	AssertEqual(t, "m.backing[3].key", num{3}, simpleMap.backing[3].key)
	AssertEqual(t, "m.backing[3].val", 54, simpleMap.backing[3].value)

	AssertEqual(t, "m.backing[7].key", num{7}, simpleMap.backing[7].key)
	AssertEqual(t, "m.backing[7].val", 2000, simpleMap.backing[7].value)

	// test change value
	simpleMap.Put(num{7}, 2023)

	AssertEqual(t, "m.len", 4, simpleMap.len)

	AssertEqual(t, "m.backing[7].key", num{7}, simpleMap.backing[7].key)
	AssertEqual(t, "m.backing[7].val", 2023, simpleMap.backing[7].value)
}

func TestGet(t *testing.T) {
	simpleMap := NewMap[num, int]()
	initMap(simpleMap)

	val, ok := simpleMap.Get(num{16})
	AssertEqual(t, "ok", true, ok)
	AssertEqual(t, "val", 91, *val)

	val, ok = simpleMap.Get(num{0})
	AssertEqual(t, "ok", true, ok)
	AssertEqual(t, "val", 85, *val)

	val, ok = simpleMap.Get(num{7})
	AssertEqual(t, "ok", true, ok)
	AssertEqual(t, "val", 2000, *val)

	val, ok = simpleMap.Get(num{40})
	AssertEqual(t, "ok", false, ok)
	AssertEqual(t, "val", nil, val)
}

func TestRemove(t *testing.T) {
	simpleMap := NewMap[num, int]()
	initMap(simpleMap)

	val, ok := simpleMap.Remove(num{16})
	AssertEqual(t, "16 ok", true, ok)
	AssertEqual(t, "16 val", 91, *val)
	AssertEqual(t, "16 m.len", 3, simpleMap.len)

	val, ok = simpleMap.Remove(num{0})
	AssertEqual(t, "0 ok", true, ok)
	AssertEqual(t, "0 val", 85, *val)
	AssertEqual(t, "0 m.len", 2, simpleMap.len)

	val, ok = simpleMap.Remove(num{40})
	AssertEqual(t, "40 ok", false, ok)
	AssertEqual(t, "40 val", nil, val)
	AssertEqual(t, "40 m.len", 2, simpleMap.len)

	val, ok = simpleMap.Remove(num{7})
	AssertEqual(t, "7 ok", true, ok)
	AssertEqual(t, "7 val", 2000, *val)
	AssertEqual(t, "7 m.len", 1, simpleMap.len)
}

func TestResize(t *testing.T) {
	m := NewMap[num, any]()
	for i := 0; i < 64; i++ {
		m.Put(num{float64(i)}, i)
	}

	AssertEqual(t, "m.len", 64, m.len)
	AssertEqual(t, "len(m.backing)", 128, len(m.backing))
}

func TestLoadFactor(t *testing.T) {
	m := NewMap[num, any]()
	m.loadFactor = 2
	for i := 0; i < 64; i++ {
		m.Put(num{float64(i)}, i)
	}

	AssertEqual(t, "m.len", 64, m.len)
	AssertEqual(t, "len(m.backing)", 32, len(m.backing))

	m.Put(num{float64(5)}, 13)

	AssertEqual(t, "m.len", 64, m.len)
	AssertEqual(t, "len(m.backing)", 64, len(m.backing))

	m.Put(num{float64(160)}, 13)

	AssertEqual(t, "m.len", 65, m.len)
	AssertEqual(t, "len(m.backing)", 64, len(m.backing))

	m.loadFactor = 0.5

	m.Put(num{5}, 43)

	AssertEqual(t, "m.len", 65, m.len)
	AssertEqual(t, "len(m.backing)", 256, len(m.backing))
}
