package f2poly

import (
	"fmt"
	"github.com/johnkerl/goffl/pkg/bitarith"
	"math/rand"
)

var formatHex = true

func SetHexOutput()    { formatHex = true }
func SetBinaryOutput() { formatHex = false }

func bitDegree(bits uint64) int {
	d := bitarith.MsbPos(bits)
	if d == -1 {
		return 0
	}
	return d
}

func bitMul(this, that uint64) uint64 {
	a, b, c := this, that, uint64(0)
	ashift := a
	degThat := bitDegree(that)
	for j := 0; j <= degThat; j++ {
		if ((b >> j) & 1) == 1 {
			c ^= ashift
		}
		ashift <<= 1
	}
	return c
}

func iquotAndRem(this, that uint64) (quot, rem uint64, err error) {
	if that == 0 {
		return 0, 0, fmt.Errorf("division by zero")
	}
	divisorL1Pos := bitDegree(that)
	if this == 0 {
		return 0, 0, nil
	}
	dividendL1Pos := bitDegree(this)
	l1Diff := dividendL1Pos - divisorL1Pos
	if l1Diff < 0 {
		return 0, this, nil
	}
	shiftDivisor := that << l1Diff
	quot, rem = 0, this
	checkPos := dividendL1Pos
	quotPos := l1Diff
	for checkPos >= divisorL1Pos {
		if ((rem >> checkPos) & 1) == 1 {
			rem ^= shiftDivisor
			quot |= 1 << quotPos
		}
		shiftDivisor >>= 1
		checkPos--
		quotPos--
	}
	return quot, rem, nil
}

// F2Poly is a polynomial over GF(2). Coefficients are bits (bit j = coefficient of x^j).
type F2Poly struct {
	Bits uint64
}

func New(bits uint64) *F2Poly { return &F2Poly{Bits: bits} }

func NewFromHex(s string) (*F2Poly, error) {
	var x uint64
	_, err := fmt.Sscanf(s, "%x", &x)
	if err != nil {
		return nil, err
	}
	return &F2Poly{Bits: x}, nil
}

func (f *F2Poly) String() string {
	if formatHex {
		return fmt.Sprintf("%x", f.Bits)
	}
	return fmt.Sprintf("%b", f.Bits)
}

func (f *F2Poly) Equal(other *F2Poly) bool { return f.Bits == other.Bits }

func (f *F2Poly) IsZero() bool    { return f.Bits == 0 }
func (f *F2Poly) IsNonzero() bool { return f.Bits != 0 }
func (f *F2Poly) IsOne() bool     { return f.Bits == 1 }

func (f *F2Poly) Degree() int { return bitDegree(f.Bits) }

func (f *F2Poly) Get(j int) int {
	if (f.Bits>>j)&1 == 1 {
		return 1
	}
	return 0
}

func (f *F2Poly) Set(j int, v int) {
	if v&1 == 1 {
		f.Bits |= 1 << j
	} else {
		f.Bits &^= 1 << j
	}
}

func (f *F2Poly) Add(other *F2Poly) *F2Poly {
	return &F2Poly{Bits: f.Bits ^ other.Bits}
}

func (f *F2Poly) Sub(other *F2Poly) *F2Poly {
	return &F2Poly{Bits: f.Bits ^ other.Bits}
}

func (f *F2Poly) Neg() *F2Poly { return &F2Poly{Bits: f.Bits} }

func (f *F2Poly) Mul(other *F2Poly) *F2Poly {
	return &F2Poly{Bits: bitMul(f.Bits, other.Bits)}
}

func (f *F2Poly) QuoRem(other *F2Poly) (q, r *F2Poly, err error) {
	quot, rem, err := iquotAndRem(f.Bits, other.Bits)
	if err != nil {
		return nil, nil, err
	}
	return &F2Poly{Bits: quot}, &F2Poly{Bits: rem}, nil
}

func (f *F2Poly) Quo(other *F2Poly) *F2Poly {
	q, _, _ := f.QuoRem(other)
	return q
}

func (f *F2Poly) Mod(other *F2Poly) *F2Poly {
	_, r, _ := f.QuoRem(other)
	return r
}

func (f *F2Poly) Pow(e int) (*F2Poly, error) {
	if f.Bits == 0 {
		if e == 0 {
			return nil, fmt.Errorf("0**0 undefined")
		}
		if e < 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return &F2Poly{Bits: 0}, nil
	}
	if e < 0 {
		return nil, fmt.Errorf("negative exponents disallowed")
	}
	rv := &F2Poly{Bits: 1}
	xp := &F2Poly{Bits: f.Bits}
	for e != 0 {
		if e&1 == 1 {
			rv = rv.Mul(xp)
		}
		e >>= 1
		xp = xp.Mul(xp)
	}
	return rv, nil
}

func (f *F2Poly) Less(other *F2Poly) bool { return f.Bits < other.Bits }

func (f *F2Poly) Gcd(other *F2Poly) *F2Poly {
	if f.Bits == 0 {
		return &F2Poly{Bits: other.Bits}
	}
	if other.Bits == 0 {
		return &F2Poly{Bits: f.Bits}
	}
	c, d := f.Bits, other.Bits
	for {
		_, r, _ := iquotAndRem(c, d)
		if r == 0 {
			break
		}
		c, d = d, r
	}
	return &F2Poly{Bits: d}
}

func (f *F2Poly) Lcm(other *F2Poly) *F2Poly {
	return f.Mul(other).Quo(f.Gcd(other))
}

func (f *F2Poly) ExtGcd(other *F2Poly) (g, s, t *F2Poly) {
	if f.Bits == 0 {
		return &F2Poly{Bits: other.Bits}, &F2Poly{Bits: 0}, &F2Poly{Bits: 1}
	}
	if other.Bits == 0 {
		return &F2Poly{Bits: f.Bits}, &F2Poly{Bits: 1}, &F2Poly{Bits: 0}
	}
	sprime, tVal := uint64(1), uint64(1)
	sVal, tprime := uint64(0), uint64(0)
	c, d := f.Bits, other.Bits
	for {
		q, r, _ := iquotAndRem(c, d)
		if r == 0 {
			break
		}
		c, d = d, r
		sprime, sVal = sVal, sprime^bitMul(q, sVal)
		tprime, tVal = tVal, tprime^bitMul(q, tVal)
	}
	return &F2Poly{Bits: d}, &F2Poly{Bits: sVal}, &F2Poly{Bits: tVal}
}

func (f *F2Poly) Deriv() *F2Poly {
	mask := uint64(0x55555555)
	for mask < f.Bits {
		mask = (mask << 32) | 0x55555555
	}
	bits := (f.Bits >> 1) & mask
	return &F2Poly{Bits: bits}
}

func (f *F2Poly) SquareRoot() (ok bool, sq *F2Poly) {
	deg := f.Degree()
	sqrootBits := uint64(0)
	inbit := uint64(1)
	outbit := uint64(1)
	si := 0
	for si <= deg {
		if (f.Bits & inbit) != 0 {
			sqrootBits |= outbit
		}
		inbit <<= 1
		if (f.Bits & inbit) != 0 {
			return false, nil
		}
		inbit <<= 1
		outbit <<= 1
		si += 2
	}
	return true, &F2Poly{Bits: sqrootBits}
}

func Random(degree int) *F2Poly {
	msb := uint64(1 << degree)
	return &F2Poly{Bits: msb | uint64(rand.Uint32())%msb}
}
