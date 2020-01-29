package marshal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tchajed/goose/machine/disk"
)

func TestUInt64(t *testing.T) {
	assert := assert.New(t)
	numbers := []uint64{0, 123, 1 << 58, 1 << 48}
	enc := NewEnc()
	for _, n := range numbers {
		enc.PutInt(n)
	}
	b := enc.Finish()

	dec := NewDec(b)
	for i, n := range numbers {
		assert.Equal(n, dec.GetInt(), "encode-decode index %d", i)
	}
}

func TestUInt32(t *testing.T) {
	assert := assert.New(t)
	numbers := []uint32{0, 123, 1<<0 | 1<<15 | 1<<31, 1 << 22, 1<<32 - 1}
	enc := NewEnc()
	for _, n := range numbers {
		enc.PutInt32(n)
	}
	b := enc.Finish()

	dec := NewDec(b)
	for i, n := range numbers {
		assert.Equal(n, dec.GetInt32(), "encode-decode index %d", i)
	}
}

type data struct {
	a uint64
	b uint32
	c uint64
	d uint32
}

func (s data) encode() disk.Block {
	enc := NewEnc()
	enc.PutInt(s.a)
	enc.PutInt32(s.b)
	enc.PutInt(s.c)
	enc.PutInt32(s.d)
	return enc.Finish()
}

func decodeData(b disk.Block) data {
	dec := NewDec(b)
	// note that this works correctly, but there's a subtle difference between
	// Go's program order initialization and Goose's declaration order (in this
	// case these are the same, as one would expect from idiomatic code)
	return data{
		a: dec.GetInt(),
		b: dec.GetInt32(),
		c: dec.GetInt(),
		d: dec.GetInt32(),
	}
}

func TestMixed(t *testing.T) {
	assert := assert.New(t)
	s := data{4, 1, 17, 25}
	assert.Equal(s, decodeData(s.encode()))
}

func TestFillBlock(t *testing.T) {
	assert := assert.New(t)
	var numbers []uint64
	for i := 0; i < 4096/8; i++ {
		numbers = append(numbers, 1)
	}
	enc := NewEnc()
	for _, n := range numbers {
		enc.PutInt(n)
	}
	b := enc.Finish()

	dec := NewDec(b)
	for i, n := range numbers {
		assert.Equal(n, dec.GetInt(), "encode-decode index %d", i)
	}
}

func TestInts(t *testing.T) {
	assert := assert.New(t)
	numbers := []uint64{2, 4, 10, 23}
	enc := NewEnc()
	enc.PutInts(numbers)
	b := enc.Finish()

	dec := NewDec(b)
	assert.Equal(numbers, dec.GetInts(uint64(len(numbers))))
}
