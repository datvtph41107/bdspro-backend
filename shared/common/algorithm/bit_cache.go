package _algorithm

import (
	"encoding/binary"
	"os"
)

type BitCache struct {
	data []uint64
}

func NewBitCache(size int) *BitCache {
	return &BitCache{
		data: make([]uint64, (size+63)/64),
	}
}

func (b *BitCache) Set(i int, value bool) {
	if value {
		b.data[i>>6] |= 1 << (i & 63)
	} else {
		b.data[i>>6] &= ^(1 << (i & 63))
	}
}

func (b *BitCache) Get(i int) bool {
	return (b.data[i>>6] & (1 << (i & 63))) != 0
}

func (b *BitCache) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return binary.Write(f, binary.LittleEndian, b.data)
}

func (b *BitCache) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return binary.Read(f, binary.LittleEndian, &b.data)
}
