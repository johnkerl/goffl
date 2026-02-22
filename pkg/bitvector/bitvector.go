package bitvector

import (
	"fmt"
	"github.com/johnkerl/goffl/pkg/bitarith"
)

// writeHex controls String() output (hex vs binary). Process-wide; not safe for
// concurrent use with different settings. Use SetHexOutput/SetBinaryOutput.
var writeHex bool

func SetHexOutput()    { writeHex = true }
func SetBinaryOutput() { writeHex = false }

// BitVector is a fixed-length bit vector. Bit position 0 is the LSB (rightmost).
type BitVector struct {
	numBits int
	Bits    uint64
}

func New(numBits int) (*BitVector, error) {
	if numBits <= 0 {
		return nil, fmt.Errorf("BitVector: size must be > 0; got %d", numBits)
	}
	return &BitVector{numBits: numBits}, nil
}

func (v *BitVector) NumBits() int { return v.numBits }

func (v *BitVector) String() string {
	mask := uint64(1<<v.numBits) - 1
	if v.numBits == 64 {
		mask = 0xFFFFFFFFFFFFFFFF
	}
	x := v.Bits & mask
	if writeHex {
		width := (v.numBits + 3) >> 2
		return fmt.Sprintf("%0*x", width, x)
	}
	return fmt.Sprintf("%0*b", v.numBits, x)
}

func (v *BitVector) Get(j int) (int, error) {
	if j < 0 || j >= v.numBits {
		return 0, fmt.Errorf("index %d out of bounds 0..%d", j, v.numBits-1)
	}
	return int((v.Bits >> j) & 1), nil
}

func (v *BitVector) Set(j int, val int) error {
	if j < 0 || j >= v.numBits {
		return fmt.Errorf("index %d out of bounds 0..%d", j, v.numBits-1)
	}
	if val&1 == 1 {
		v.Bits |= 1 << j
	} else {
		v.Bits &^= 1 << j
	}
	return nil
}

func (v *BitVector) ToggleElement(j int) error {
	if j < 0 || j >= v.numBits {
		return fmt.Errorf("index %d out of bounds 0..%d", j, v.numBits-1)
	}
	v.Bits ^= 1 << j
	return nil
}

// FindLeaderPos returns the position of the lowest set bit (for row-reduction), or -1 if zero.
func (v *BitVector) FindLeaderPos() int {
	return bitarith.LsbPos(v.Bits)
}
