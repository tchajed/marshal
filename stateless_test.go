package marshal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
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
