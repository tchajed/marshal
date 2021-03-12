package marshal

import (
	"github.com/tchajed/goose/machine"
)

// Enc is a stateful encoder for a statically-allocated array.
type Enc struct {
	b   []byte
	off *uint64
}

func NewEnc(sz uint64) Enc {
	return Enc{
		b:   make([]byte, sz),
		off: new(uint64),
	}
}

func (enc Enc) PutInt(x uint64) {
	off := *enc.off
	machine.UInt64Put(enc.b[off:], x)
	*enc.off += 8
}

func (enc Enc) PutInt32(x uint32) {
	off := *enc.off
	machine.UInt32Put(enc.b[off:], x)
	*enc.off += 4
}

func (enc Enc) PutInts(xs []uint64) {
	for _, x := range xs {
		enc.PutInt(x)
	}
}

func (enc Enc) PutBytes(b []byte) {
	off := *enc.off
	n := uint64(copy(enc.b[off:], b))
	*enc.off += n // should be len(b) (unless too much data was provided)
}

func (enc Enc) PutBool(b bool) {
	off := *enc.off
	if !b {
		enc.b[off] = 0
	} else {
		enc.b[off] = 1
	}
	*enc.off += 1
}

func (enc Enc) Finish() []byte {
	return enc.b
}

// Dec is a stateful decoder that returns values encoded
// sequentially in a single slice.
type Dec struct {
	b   []byte
	off *uint64
}

func NewDec(b []byte) Dec {
	return Dec{b: b, off: new(uint64)}
}

func (dec Dec) GetInt() uint64 {
	off := *dec.off
	*dec.off += 8
	return machine.UInt64Get(dec.b[off:])
}

func (dec Dec) GetInt32() uint32 {
	off := *dec.off
	*dec.off += 4
	return machine.UInt32Get(dec.b[off:])
}

func (dec Dec) GetInts(num uint64) []uint64 {
	var xs []uint64
	for i := uint64(0); i < num; i++ {
		xs = append(xs, dec.GetInt())
	}
	return xs
}

func (dec Dec) GetBytes(num uint64) []byte {
	off := *dec.off
	b := dec.b[off : off+num]
	*dec.off += num
	return b
}

func (dec Dec) GetBool() bool {
	off := *dec.off
	*dec.off += 1
	if dec.b[off] == 0 {
		return false
	} else {
		return true
	}
}
