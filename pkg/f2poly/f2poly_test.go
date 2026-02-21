package f2poly_test

import (
	"github.com/johnkerl/goffl/pkg/f2poly"
	"github.com/johnkerl/goffl/pkg/f2polyfactor"
	"testing"
)

func TestDegree(t *testing.T) {
	if got := f2poly.New(0).Degree(); got != 0 {
		t.Errorf("F2Poly(0).Degree() = %d, want 0", got)
	}
	if got := f2poly.New(1).Degree(); got != 0 {
		t.Errorf("F2Poly(1).Degree() = %d, want 0", got)
	}
	if got := f2poly.New(2).Degree(); got != 1 {
		t.Errorf("F2Poly(2).Degree() = %d, want 1", got)
	}
	if got := f2poly.New(0x13).Degree(); got != 4 {
		t.Errorf("F2Poly(0x13).Degree() = %d, want 4", got)
	}
}

func TestArithmetic(t *testing.T) {
	a := f2poly.New(0x13)
	b := f2poly.New(0x0B)
	if got := a.Add(b).Bits; got != 0x13^0x0B {
		t.Errorf("(a+b).Bits = 0x%x, want 0x%x", got, 0x13^0x0B)
	}
	if !a.Mul(f2poly.New(1)).Equal(a) {
		t.Error("a*1 != a")
	}
	q := a.Quo(b)
	r := a.Mod(b)
	if q.Mul(b).Add(r).Bits != a.Bits {
		t.Error("q*b + r != a")
	}
}

func TestGcd(t *testing.T) {
	a := f2poly.New(0x13)
	b := f2poly.New(0x0B)
	g := a.Gcd(b)
	if !g.IsOne() && !g.Equal(a) && !g.Equal(b) {
		t.Error("gcd should be 1, a, or b")
	}
}

func TestIrr(t *testing.T) {
	// x+1 (0x3) has degree 1, so irreducible
	if !f2polyfactor.Irr(f2poly.New(0x3)) {
		t.Error("F2Poly(0x3) should be irreducible")
	}
	// x^2 (0x4) is reducible
	if f2polyfactor.Irr(f2poly.New(0x4)) {
		t.Error("F2Poly(0x4) should not be irreducible")
	}
}

func TestFactor(t *testing.T) {
	// x^2+x+1 (0x7) is irreducible
	finfo := f2polyfactor.Factor(f2poly.New(0x7))
	if finfo.NumFactors() < 1 {
		t.Error("factor(0x7) should have at least one factor")
	}
}
