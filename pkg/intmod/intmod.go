package intmod

import (
	"github.com/johnkerl/goffl/pkg/intarith"
	"math/rand"
)

// IntMod is a residue class: integer mod m.
type IntMod struct {
	Residue int64
	modulus int64
}

func New(residue, modulus int64) *IntMod {
	m := modulus
	if m < 0 {
		m = -m
	}
	r := residue % m
	if r < 0 {
		r += m
	}
	return &IntMod{Residue: r, modulus: m}
}

func (a *IntMod) Modulus() int64 { return a.modulus }

func (a *IntMod) Recip() *IntMod {
	return New(intarith.IntModRecip(a.Residue, a.modulus), a.modulus)
}

func (a *IntMod) Add(other *IntMod) *IntMod {
	if a.modulus != other.modulus {
		panic("modulus mismatch")
	}
	r := (a.Residue + other.Residue) % a.modulus
	if r < 0 {
		r += a.modulus
	}
	return &IntMod{Residue: r, modulus: a.modulus}
}

func (a *IntMod) Sub(other *IntMod) *IntMod {
	if a.modulus != other.modulus {
		panic("modulus mismatch")
	}
	r := (a.Residue - other.Residue) % a.modulus
	if r < 0 {
		r += a.modulus
	}
	return &IntMod{Residue: r, modulus: a.modulus}
}

func (a *IntMod) Neg() *IntMod {
	r := (-a.Residue) % a.modulus
	if r < 0 {
		r += a.modulus
	}
	return &IntMod{Residue: r, modulus: a.modulus}
}

func (a *IntMod) Mul(other *IntMod) *IntMod {
	if a.modulus != other.modulus {
		panic("modulus mismatch")
	}
	r := (a.Residue * other.Residue) % a.modulus
	if r < 0 {
		r += a.modulus
	}
	return &IntMod{Residue: r, modulus: a.modulus}
}

func (a *IntMod) Div(other *IntMod) *IntMod {
	return a.Mul(other.Recip())
}

func (a *IntMod) Pow(e int64) *IntMod {
	if a.Residue == 0 {
		if e == 0 {
			panic("0**0 undefined")
		}
		if e < 0 {
			panic("division by zero")
		}
		return New(0, a.modulus)
	}
	rv := New(1, a.modulus)
	xp := New(a.Residue, a.modulus)
	if e < 0 {
		xp = xp.Recip()
		e = -e
	}
	for e != 0 {
		if e&1 == 1 {
			rv = rv.Mul(xp)
		}
		e >>= 1
		xp = xp.Mul(xp)
	}
	return rv
}

func (a *IntMod) Equal(other *IntMod) bool {
	return a.Residue == other.Residue && a.modulus == other.modulus
}

func Random(m int64) *IntMod {
	return New(rand.Int63n(m), m)
}

func ElementsForModulus(m int64) []*IntMod {
	out := make([]*IntMod, m)
	for a := int64(0); a < m; a++ {
		out[a] = New(a, m)
	}
	return out
}

func UnitsForModulus(m int64) []*IntMod {
	var out []*IntMod
	for a := int64(0); a < m; a++ {
		if intarith.Gcd(a, m) == 1 {
			out = append(out, New(a, m))
		}
	}
	return out
}
