package marshal

import (
	"github.com/tchajed/goose/machine"
	"github.com/tchajed/goose/machine/disk"
)

// Enc is a stateful encoder for a single disk block.
type Enc struct {
	b   disk.Block
	off *uint64
}

func NewEnc() Enc {
	return Enc{
		b:   make(disk.Block, disk.BlockSize),
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

func (enc Enc) Finish() disk.Block {
	return enc.b
}

// Dec is a stateful decoder that returns values encoded
// sequentially in a single disk block.
type Dec struct {
	b   disk.Block
	off *uint64
}

func NewDec(b disk.Block) Dec {
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
