package marshal

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
