package f2polymod

import (
	"fmt"
	"goffl/pkg/f2poly"
)

// F2PolyMod is a residue class: F2[x] mod m(x).
type F2PolyMod struct {
	Residue *f2poly.F2Poly
	modulus *f2poly.F2Poly
}

func New(residue *f2poly.F2Poly, modulus *f2poly.F2Poly) *F2PolyMod {
	if residue == nil {
		residue = &f2poly.F2Poly{Bits: 0}
	}
	if modulus == nil {
		modulus = &f2poly.F2Poly{Bits: 1}
	}
	r := residue.Mod(modulus)
	return &F2PolyMod{Residue: r, modulus: modulus}
}

func NewFromInts(residueBits uint64, modulusBits uint64) *F2PolyMod {
	return New(&f2poly.F2Poly{Bits: residueBits}, &f2poly.F2Poly{Bits: modulusBits})
}

func (a *F2PolyMod) Modulus() *f2poly.F2Poly { return a.modulus }

func (a *F2PolyMod) IsZero() bool { return a.Residue.IsZero() }
func (a *F2PolyMod) IsOne() bool  { return a.Residue.IsOne() }

func (a *F2PolyMod) Add(other *F2PolyMod) *F2PolyMod {
	r := a.Residue.Add(other.Residue).Mod(a.modulus)
	return &F2PolyMod{Residue: r, modulus: a.modulus}
}

func (a *F2PolyMod) Sub(other *F2PolyMod) *F2PolyMod {
	r := a.Residue.Sub(other.Residue).Mod(a.modulus)
	return &F2PolyMod{Residue: r, modulus: a.modulus}
}

func (a *F2PolyMod) Neg() *F2PolyMod {
	r := a.Residue.Neg().Mod(a.modulus)
	return &F2PolyMod{Residue: r, modulus: a.modulus}
}

func (a *F2PolyMod) Mul(other *F2PolyMod) *F2PolyMod {
	r := a.Residue.Mul(other.Residue).Mod(a.modulus)
	return &F2PolyMod{Residue: r, modulus: a.modulus}
}

func (a *F2PolyMod) Recip() (*F2PolyMod, error) {
	g, s, _ := a.Residue.ExtGcd(a.modulus)
	if !g.IsOne() {
		return nil, fmt.Errorf("recip: division by zero")
	}
	return &F2PolyMod{Residue: s, modulus: a.modulus}, nil
}

func (a *F2PolyMod) Div(other *F2PolyMod) (*F2PolyMod, error) {
	rec, err := other.Recip()
	if err != nil {
		return nil, err
	}
	return a.Mul(rec), nil
}

func (a *F2PolyMod) Pow(e int) (*F2PolyMod, error) {
	if a.Residue.IsZero() {
		if e == 0 {
			return nil, fmt.Errorf("0**0 undefined")
		}
		if e < 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return &F2PolyMod{Residue: &f2poly.F2Poly{Bits: 0}, modulus: a.modulus}, nil
	}
	rv := &F2PolyMod{Residue: &f2poly.F2Poly{Bits: 1}, modulus: a.modulus}
	xp := &F2PolyMod{Residue: &f2poly.F2Poly{Bits: a.Residue.Bits}, modulus: a.modulus}
	if e < 0 {
		var err error
		xp, err = xp.Recip()
		if err != nil {
			return nil, err
		}
		e = -e
	}
	for e != 0 {
		if e&1 == 1 {
			rv = rv.Mul(xp)
		}
		e >>= 1
		xp = xp.Mul(xp)
	}
	return rv, nil
}

func (a *F2PolyMod) Equal(other *F2PolyMod) bool {
	return a.Residue.Equal(other.Residue) && a.modulus.Equal(other.modulus)
}

func ElementsForModulus(m *f2poly.F2Poly) []*F2PolyMod {
	maxBits := uint64(1<<m.Degree()) - 1
	if m.Degree() >= 64 {
		maxBits = 0xFFFFFFFFFFFFFFFF
	}
	out := make([]*F2PolyMod, 0, maxBits+1)
	for a := uint64(0); a <= maxBits; a++ {
		out = append(out, New(&f2poly.F2Poly{Bits: a}, m))
	}
	return out
}

func UnitsForModulus(m *f2poly.F2Poly) []*F2PolyMod {
	maxBits := uint64(1<<m.Degree()) - 1
	if m.Degree() >= 64 {
		maxBits = 0xFFFFFFFFFFFFFFFF
	}
	var out []*F2PolyMod
	for j := uint64(1); j <= maxBits; j++ {
		g := m.Gcd(&f2poly.F2Poly{Bits: j})
		if g.IsOne() {
			out = append(out, New(&f2poly.F2Poly{Bits: j}, m))
		}
	}
	return out
}
