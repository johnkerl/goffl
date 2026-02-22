// Package order provides multiplicative order and related functions for Z/nZ and F2[x]/(m).
package order

import (
	"fmt"

	"github.com/johnkerl/goffl/pkg/f2poly"
	"github.com/johnkerl/goffl/pkg/f2polyfactor"
	"github.com/johnkerl/goffl/pkg/f2polymod"
	"github.com/johnkerl/goffl/pkg/intarith"
	"github.com/johnkerl/goffl/pkg/intfactor"
	"github.com/johnkerl/goffl/pkg/intmod"
)

// ModOrderIntMod returns the multiplicative order of a in Z/mZ.
func ModOrderIntMod(am *intmod.IntMod) (int64, error) {
	a, m := am.Residue, am.Modulus()
	if intarith.Gcd(a, m) != 1 {
		return 0, fmt.Errorf("mod_order: zero or zero divisor %d mod %d", a, m)
	}
	phi := intfactor.Totient(m)
	finfo := intfactor.Factor(phi)
	phiDivisors := finfo.AllDivisors()
	rec, err := am.Recip()
	if err != nil {
		return 0, fmt.Errorf("mod_order: %w", err)
	}
	one := am.Mul(rec)

	for _, e := range phiDivisors {
		pow, err := am.Pow(e)
		if err != nil {
			return 0, fmt.Errorf("mod_order: %w", err)
		}
		if pow.Equal(one) {
			return e, nil
		}
	}
	return 0, fmt.Errorf("mod_order: coding error")
}

// ModOrderF2PolyMod returns the multiplicative order of a in F2[x]/(m).
func ModOrderF2PolyMod(am *f2polymod.F2PolyMod) (int64, error) {
	a, m := am.Residue, am.Modulus()
	if !a.Gcd(m).IsOne() {
		return 0, fmt.Errorf("mod_order: zero or zero divisor mod m")
	}
	phi := f2polyfactor.Totient(m)
	finfo := intfactor.Factor(phi)
	phiDivisors := finfo.AllDivisors()
	rec, err := am.Recip()
	if err != nil {
		return 0, fmt.Errorf("mod_order: %w", err)
	}
	one := am.Mul(rec)

	for _, e := range phiDivisors {
		pow, err := am.Pow(int(e))
		if err != nil {
			return 0, fmt.Errorf("mod_order: %w", err)
		}
		if pow.Equal(one) {
			return e, nil
		}
	}
	return 0, fmt.Errorf("mod_order: coding error")
}

func ModMaxOrderInt(m int64) (int64, error) {
	units := intmod.UnitsForModulus(m)
	var max int64
	for _, a := range units {
		ord, err := ModOrderIntMod(a)
		if err != nil {
			return 0, err
		}
		if ord > max {
			max = ord
		}
	}
	return max, nil
}

func ModMaxOrderF2Poly(m *f2poly.F2Poly) (int64, error) {
	units, err := f2polymod.UnitsForModulus(m)
	if err != nil {
		return 0, err
	}
	var max int64
	for _, a := range units {
		ord, err := ModOrderF2PolyMod(a)
		if err != nil {
			return 0, err
		}
		if ord > max {
			max = ord
		}
	}
	return max, nil
}

func OrbitIntMod(am *intmod.IntMod, bm *intmod.IntMod) []*intmod.IntMod {
	var orbit []*intmod.IntMod
	cm := intmod.New(am.Residue, am.Modulus())
	for {
		if bm == nil {
			orbit = append(orbit, intmod.New(cm.Residue, cm.Modulus()))
		} else {
			orbit = append(orbit, cm.Mul(bm))
		}
		if cm.Residue == 1 {
			break
		}
		cm = cm.Mul(am)
	}
	return orbit
}

func OrbitF2PolyMod(am *f2polymod.F2PolyMod, bm *f2polymod.F2PolyMod) []*f2polymod.F2PolyMod {
	var orbit []*f2polymod.F2PolyMod
	cm := f2polymod.New(&f2poly.F2Poly{Bits: am.Residue.Bits}, am.Modulus())
	for {
		if bm == nil {
			orbit = append(orbit, f2polymod.New(&f2poly.F2Poly{Bits: cm.Residue.Bits}, cm.Modulus()))
		} else {
			orbit = append(orbit, cm.Mul(bm))
		}
		if cm.Residue.IsOne() {
			break
		}
		cm = cm.Mul(am)
	}
	return orbit
}

// F2PolyPeriod returns the period of x in F2[x]/(m), or 0 if x is not a unit or on error.
func F2PolyPeriod(m *f2poly.F2Poly) int64 {
	x := &f2poly.F2Poly{Bits: 2}
	if !x.Gcd(m).IsOne() {
		return 0
	}
	ord, err := ModOrderF2PolyMod(f2polymod.New(x, m))
	if err != nil {
		return 0
	}
	return ord
}

func IntModGenerator(m int64) (int64, bool) {
	if m < 2 {
		panic("int_mod_generator: modulus must be >= 2")
	}
	phi := intfactor.Totient(m)
	for a := int64(1); a < m; a++ {
		if intarith.Gcd(a, m) == 1 {
			g := intmod.New(a, m)
			ord, err := ModOrderIntMod(g)
			if err == nil && ord == phi {
				return g.Residue, true
			}
		}
	}
	return 0, false
}

func F2PolyModGenerator(m *f2poly.F2Poly) (*f2poly.F2Poly, bool) {
	mdeg := m.Degree()
	if mdeg < 1 {
		panic("f2_poly_mod_generator: modulus degree must be positive")
	}
	phi := f2polyfactor.Totient(m)
	maxBits := uint64(1<<mdeg) - 1
	if mdeg >= 64 {
		maxBits = 0xFFFFFFFFFFFFFFFF
	}
	for bits := uint64(1); bits <= maxBits; bits++ {
		gRes := &f2poly.F2Poly{Bits: bits}
		if gRes.Gcd(m).IsOne() {
			g := f2polymod.New(gRes, m)
			ord, err := ModOrderF2PolyMod(g)
			if err == nil && ord == phi {
				return g.Residue, true
			}
		}
	}
	return nil, false
}

func F2PolyPrimitive(m *f2poly.F2Poly) bool {
	x := &f2poly.F2Poly{Bits: 2}
	if !m.Gcd(x).IsOne() {
		return false
	}
	rcrx := f2polymod.New(x, m)
	phi := f2polyfactor.Totient(m)
	finfo := intfactor.Factor(phi)
	mpds := finfo.MaximalProperDivisors()

	for _, mpd := range mpds {
		pow, err := rcrx.Pow(int(mpd))
		if err != nil {
			return false
		}
		if pow.IsOne() {
			return false
		}
	}
	pow, err := rcrx.Pow(int(phi))
	if err != nil {
		return false
	}
	return pow.IsOne()
}
