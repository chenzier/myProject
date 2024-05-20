package blond

import "fmt"

type bitmap struct {
	keys []byte

	len int
}

func NewBitMap() *bitmap {

	return &bitmap{keys: make([]byte, 0), len: 0}

}

func (b *bitmap) has(v int) bool {

	k := v / 8

	kv := byte(v % 8)

	if k > len(b.keys) { //todo not exist

		return false

	}

	if b.keys[k]&(1<<kv) != 0 {

		return true

	}

	return false

}

func (b *bitmap) set(v int) {

	k := v / 8

	kv := byte(v % 8)

	for b.len <= k {

		b.keys = append(b.keys, 0)

		b.len++

	}

	b.keys[k] = b.keys[k] | (1 << kv)

}

func (b *bitmap) length() int {

	return b.len

}

func (b *bitmap) print() {

	for _, v := range b.keys {
		fmt.Printf("%08b\n", v)

	}

}
