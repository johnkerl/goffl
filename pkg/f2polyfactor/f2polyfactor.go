package f2polyfactor

import (
	"fmt"
	"goffl/pkg/bitmatrix"
	"goffl/pkg/bitvector"
	"goffl/pkg/f2poly"
)

// PolyFactorization holds factors (F2Poly) with multiplicities for polynomial factorization.
type PolyFactorization struct {
	trivial *f2poly.F2Poly
	factors []struct {
		f *f2poly.F2Poly
		m int
	}
}

func NewPolyFactorization() *PolyFactorization {
	return &PolyFactorization{factors: nil}
}

func (pf *PolyFactorization) InsertTrivialFactor(f *f2poly.F2Poly) {
	if f == nil {
		return
	}
	if pf.trivial != nil {
		pf.trivial = pf.trivial.Mul(f)
	} else {
		pf.trivial = &f2poly.F2Poly{Bits: f.Bits}
	}
}

func (pf *PolyFactorization) InsertFactor(newFactor *f2poly.F2Poly, newMult int) {
	if newMult <= 0 {
		return
	}
	for i := range pf.factors {
		if pf.factors[i].f.Bits == newFactor.Bits {
			pf.factors[i].m += newMult
			return
		}
		if newFactor.Less(pf.factors[i].f) {
			pf.factors = append(pf.factors, struct {
				f *f2poly.F2Poly
				m int
			}{})
			copy(pf.factors[i+1:], pf.factors[i:])
			pf.factors[i] = struct {
				f *f2poly.F2Poly
				m int
			}{newFactor, newMult}
			return
		}
	}
	pf.factors = append(pf.factors, struct {
		f *f2poly.F2Poly
		m int
	}{newFactor, newMult})
}

func (pf *PolyFactorization) Merge(other *PolyFactorization) {
	if other.trivial != nil {
		pf.InsertTrivialFactor(other.trivial)
	}
	for i := 0; i < other.NumDistinctFactors(); i++ {
		f, m := other.Get(i)
		pf.InsertFactor(f, m)
	}
}

func (pf *PolyFactorization) ExpAll(e int) {
	if pf.trivial != nil {
		p, _ := pf.trivial.Pow(e)
		pf.trivial = p
	}
	for i := range pf.factors {
		pf.factors[i].m *= e
	}
}

func (pf *PolyFactorization) NumDistinctFactors() int { return len(pf.factors) }

func (pf *PolyFactorization) NumFactors() int {
	n := 0
	for _, p := range pf.factors {
		n += p.m
	}
	return n
}

func (pf *PolyFactorization) Get(i int) (*f2poly.F2Poly, int) {
	return pf.factors[i].f, pf.factors[i].m
}

var oneF2 = &f2poly.F2Poly{Bits: 1}
var xF2 = &f2poly.F2Poly{Bits: 2}
var x2F2 = &f2poly.F2Poly{Bits: 4}

func Factor(f *f2poly.F2Poly) *PolyFactorization {
	finfo := NewPolyFactorization()
	if f.Degree() == 0 {
		finfo.InsertTrivialFactor(f)
		return finfo
	}
	preBerlekamp(f, finfo, true)
	return finfo
}

func preBerlekamp(f *f2poly.F2Poly, finfo *PolyFactorization, recurse bool) {
	d := f.Deriv()
	g := f.Gcd(d)

	if g.IsZero() {
		if f.IsNonzero() {
			panic("pre_berlekamp: coding error detected")
		}
		finfo.InsertFactor(f, 1)
		return
	}
	if g.IsOne() {
		berlekamp(f, finfo, recurse)
		return
	}
	if d.IsZero() {
		ok, sqroot := f.SquareRoot()
		if !ok || sqroot == nil {
			panic("pre_berlekamp: coding error detected")
		}
		sfinfo := NewPolyFactorization()
		preBerlekamp(sqroot, sfinfo, recurse)
		if f.Degree() > 0 {
			sfinfo.ExpAll(2)
		}
		finfo.Merge(sfinfo)
		return
	}
	q := f.Quo(g)
	preBerlekamp(g, finfo, recurse)
	preBerlekamp(q, finfo, recurse)
}

func berlekamp(f *f2poly.F2Poly, finfo *PolyFactorization, recurse bool) {
	n := f.Degree()
	x2modf := x2F2.Mod(f)
	x2i := &f2poly.F2Poly{Bits: 1}

	if n < 2 {
		finfo.InsertFactor(f, 1)
		return
	}

	bi, err := bitmatrix.New(n, n)
	if err != nil {
		panic(err)
	}

	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			bi.Rows[n-1-i].Set(n-1-j, x2i.Get(i))
		}
		x2i = x2i.Mul(x2modf).Mod(f)
	}
	for i := 0; i < n; i++ {
		bi.Rows[i].ToggleElement(i)
	}

	bi.RowEchelonForm()
	rank := bi.RankRR()
	dimker := n - rank

	if dimker == 1 {
		finfo.InsertFactor(f, 1)
		return
	}

	nullspaceBasis, err := bi.KernelBasis()
	if err != nil || nullspaceBasis == nil {
		panic("coding error detected: kernel_basis")
	}
	if nullspaceBasis.NumRows() != dimker {
		panic("coding error detected: kernel_basis")
	}

	for row := 0; row < dimker; row++ {
		h := f2polyFromVector(nullspaceBasis.Row(row), n)
		hc := h.Add(oneF2)

		check1 := h.Mul(h).Mod(f)
		check2 := hc.Mul(hc).Mod(f)
		if !h.Equal(check1) || !hc.Equal(check2) {
			panic("coding error detected: h^2 check")
		}

		f1 := f.Gcd(h)
		f2 := f.Gcd(hc)

		if f1.IsOne() || f2.IsOne() {
			continue
		}

		if dimker == 2 {
			finfo.InsertFactor(f1, 1)
			finfo.InsertFactor(f2, 1)
		} else if !recurse {
			finfo.InsertFactor(f1, 1)
			finfo.InsertFactor(f2, 1)
		} else {
			preBerlekamp(f1, finfo, recurse)
			preBerlekamp(f2, finfo, recurse)
		}
		return
	}
	panic("coding error detected: berlekamp")
}

func f2polyFromVector(v *bitvector.BitVector, n int) *f2poly.F2Poly {
	f := &f2poly.F2Poly{Bits: 0}
	for i := 0; i < n; i++ {
		val, _ := v.Get(i)
		f.Set(i, val)
	}
	return f
}

func Irr(f *f2poly.F2Poly) bool {
	if f.Degree() == 0 {
		return false
	}
	if f.Degree() == 1 {
		return true
	}
	finfo := NewPolyFactorization()
	preBerlekamp(f, finfo, false)
	return finfo.NumFactors() == 1
}

func LowestIrr(degree int) (*f2poly.F2Poly, error) {
	if degree < 1 {
		return nil, fmt.Errorf("lowest_irr: degree must be positive; got %d", degree)
	}
	rv := &f2poly.F2Poly{Bits: (1 << degree) | 1}
	for rv.Degree() == degree {
		if Irr(rv) {
			return rv, nil
		}
		rv.Bits += 2
	}
	return nil, fmt.Errorf("lowest_irr: coding error detected")
}

func RandomIrr(degree int) (*f2poly.F2Poly, error) {
	if degree < 1 {
		return nil, fmt.Errorf("random_irr: degree must be positive; got %d", degree)
	}
	for {
		rv := f2poly.Random(degree)
		rv.Bits |= 1
		if Irr(rv) {
			return rv, nil
		}
	}
}

func Totient(f *f2poly.F2Poly) int64 {
	finfo := Factor(f)
	rv := int64(1)
	nf := finfo.NumDistinctFactors()
	for i := 0; i < nf; i++ {
		fi, ei := finfo.Get(i)
		di := fi.Degree()
		rv *= (1 << (di * (ei - 1))) * ((1 << di) - 1)
	}
	return rv
}
