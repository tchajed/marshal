package marshal

import (
	"github.com/goose-lang/std"
	"github.com/tchajed/goose/machine"
)

// Grow a slice to have at least `additional` unused bytes in the capacity.
// Runtime-check against overflow.
func reserve(b []byte, additional uint64) []byte {
	min_cap := std.SumAssumeNoOverflow(uint64(len(b)), additional)
	if uint64(cap(b)) < min_cap {
		// Amortized allocation strategy: grow slice by at least a certain factor.
		// Rust RawVec uses a factor of 2 so we follow that.
		var new_cap = uint64(cap(b)) * 2
		if new_cap < 8 {
			// Too small allocations aren't worth it, grow to 8 at least.
			new_cap = 8
		}
		if new_cap < min_cap {
			// Guard against overflow, cap(b)==0, and other nonsense.
			new_cap = min_cap
		}
		// We make a new slice with length 0 and the desired capacity.
		// Then we append `b` to that, which copies its elements without further allocation.
		return append(make([]byte, 0, new_cap), b...)
	} else {
		return b
	}
}

// Functions for the stateless decoder API
func ReadInt(b []byte) (uint64, []byte) {
	i := machine.UInt64Get(b)
	return i, b[8:]
}

// Functions for the stateless encoder API
func WriteInt(b []byte, i uint64) []byte {
	b2 := reserve(b, 8) // If go would let me shadow variables, this code would be much more readable
	off := len(b2)
	b3 := b2[:off+8] // yeah you can index into a slice to *increase* its length (up to its capacity)
	machine.UInt64Put(b3[off:], i)
	return b3
}
