package marshal

import (
	"github.com/goose-lang/std"
	"github.com/tchajed/goose/machine"
)

func compute_new_cap(old_cap uint64, min_cap uint64) uint64 {
	var new_cap = old_cap * 2
	if new_cap < min_cap {
		// Guard against overflow, and other nonsense.
		new_cap = min_cap
	}
	return new_cap
}

// Grow a slice to have at least `additional` unused bytes in the capacity.
// Runtime-check against overflow.
func reserve(b []byte, additional uint64) []byte {
	// This is less of a regular "assume" and more of a "we are okay with tearing down the entire
	// application (in a controlled way)" if it does not hold. We rely on SumAssumeNoOverflow
	// doing a regular Go panic at runtime if the condition fails.
	min_cap := std.SumAssumeNoOverflow(uint64(len(b)), additional)
	if uint64(cap(b)) < min_cap {
		// Amortized allocation strategy: grow slice by at least a certain factor.
		// Rust RawVec uses a factor of 2 so we follow that.
		new_cap := compute_new_cap(uint64(cap(b)), min_cap)
		// We make a new slice with length 0 and the desired capacity.
		// Then we append `b` to that, which copies its elements without further allocation.
		dest := make([]byte, len(b), new_cap)
		copy(dest, b)
		return dest
	} else {
		return b
	}
}

/* Functions for the stateless decoder API */

func ReadInt(b []byte) (uint64, []byte) {
	i := machine.UInt64Get(b)
	return i, b[8:]
}

// ReadBytes reads `l` bytes from b and returns (bs, rest)
func ReadBytes(b []byte, l uint64) ([]byte, []byte) {
	s := b[:l]
	return s, b[l:]
}

// Like ReadBytes, but avoids keeping the source slice [b] alive.
func ReadBytesCopy(b []byte, l uint64) ([]byte, []byte) {
	s := make([]byte, l)
	copy(s, b[:l])
	return s, b[l:]
}

func ReadBool(b []byte) (bool, []byte) {
	x := b[0] != 0
	return x, b[1:]
}

func ReadLenPrefixedBytes(b []byte) ([]byte, []byte) {
	l, b2 := ReadInt(b)
	bs, b3 := ReadBytes(b2, l)
	return bs, b3
}

/* Functions for the stateless encoder API */

// Encode i in little endian format and append it to b, returning the new slice.
func WriteInt(b []byte, i uint64) []byte {
	b2 := reserve(b, 8) // If go would let me shadow variables, this code would be much more readable
	off := len(b2)
	b3 := b2[:off+8] // yeah you can index into a slice to *increase* its length (up to its capacity)
	machine.UInt64Put(b3[off:], i)
	return b3
}

// Append data to b, returning the new slice.
func WriteBytes(b []byte, data []byte) []byte {
	return append(b, data...)
}

func WriteBool(b []byte, x bool) []byte {
	if x {
		return append(b, 1)
	} else {
		return append(b, 0)
	}
}

func WriteLenPrefixedBytes(b []byte, bs []byte) []byte {
	b2 := WriteInt(b, uint64(len(bs)))
	return WriteBytes(b2, bs)
}
