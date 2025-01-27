package marshal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatelessInt(t *testing.T) {
	assert := assert.New(t)
	numbers := []uint64{0, 123, 1 << 58, 1 << 48}
	var data []byte
	for _, n := range numbers {
		data = WriteInt(data, n)
	}

	for i, n := range numbers {
		n2, data2 := ReadInt(data)
		data = data2
		assert.Equal(n, n2, "encode-decode index %d", i)
	}
}

func TestStatelessInt32(t *testing.T) {
	assert := assert.New(t)
	numbers := []uint32{0, 123, 1 << 22}
	var data []byte
	for _, n := range numbers {
		data = WriteInt32(data, n)
	}

	for i, n := range numbers {
		n2, data2 := ReadInt32(data)
		data = data2
		assert.Equal(n, n2, "encode-decode index %d", i)
	}
}

func TestStatelessWriteBytesInPlace(t *testing.T) {
	assert := assert.New(t)
	data := []byte{1, 2, 3, 4}
	b := make([]byte, 2, 2+4)
	// WriteBytes has enough capacity
	b = WriteBytes(b, data)
	assert.Equal([]byte{0, 0, 1, 2, 3, 4}, b)
}

func TestStatelessWriteBytesCopy(t *testing.T) {
	assert := assert.New(t)
	data := []byte{1, 2, 3, 4}
	b := make([]byte, 2, 3)
	b[0] = 10
	// not enough capacity, WriteBytes needs to copy the slice to make room
	b = WriteBytes(b, data)
	assert.Equal([]byte{10, 0, 1, 2, 3, 4}, b)
}

func TestStatelessBool(t *testing.T) {
	assert := assert.New(t)
	bools := []bool{true, false, true, false}
	var data []byte
	for _, b := range bools {
		data = WriteBool(data, b)
	}

	for i, b := range bools {
		b2, data2 := ReadBool(data)
		data = data2
		assert.Equal(b, b2, "encode-decode index %d", i)
	}
}

type things struct {
	x  uint64
	y  uint64
	ok bool
}

func readThings(b []byte) (x things, b2 []byte) {
	b2 = b
	x.x, b2 = ReadInt(b2)
	x.y, b2 = ReadInt(b2)
	x.ok, b2 = ReadBool(b2)
	return
}

func writeThings(b []byte, x things) []byte {
	b2 := b
	b2 = WriteInt(b2, x.x)
	b2 = WriteInt(b2, x.y)
	b2 = WriteBool(b2, x.ok)
	return b2
}

func TestStatelessWriteSlice(t *testing.T) {
	assert := assert.New(t)
	xs := []things{
		{2, 1, true},
		{3, 1, false},
		{0, 7, true},
	}
	b := WriteSliceLenPrefix([]byte{}, xs, writeThings)

	xs2, b_extra := ReadSliceLenPrefix(b, readThings)
	assert.Empty(b_extra)
	assert.Equal(xs, xs2)
}

func TestStatelessIntSlice(t *testing.T) {
	assert := assert.New(t)
	numbers := []uint64{0, 123, 1 << 58, 1 << 48}
	data := WriteSlice([]byte{}, numbers, WriteInt)
	result, b_extra := ReadSlice(data, 4, ReadInt)

	assert.Empty(b_extra)
	assert.Equal(numbers, result)
}
