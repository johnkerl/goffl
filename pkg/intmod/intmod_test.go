package intmod

import (
	"github.com/johnkerl/goffl/pkg/intarith"
	"testing"
)

func TestBasic(t *testing.T) {
	a := New(2, 11)
	b := New(3, 11)
	if got := a.Add(b).Residue; got != 5 {
		t.Errorf("(a+b).Residue = %d, want 5", got)
	}
	if got := a.Mul(b).Residue; got != 6 {
		t.Errorf("(a*b).Residue = %d, want 6", got)
	}
	if got := a.Pow(10).Residue; got != 1 {
		t.Errorf("(a^10).Residue = %d, want 1", got)
	}
}

func TestRecip(t *testing.T) {
	a := New(2, 11)
	r := a.Recip()
	if got := a.Mul(r).Residue; got != 1 {
		t.Errorf("(a*recip).Residue = %d, want 1", got)
	}
}

func TestUnits(t *testing.T) {
	units := UnitsForModulus(10)
	want := intarith.EulerPhi(10)
	if int64(len(units)) != want {
		t.Errorf("len(UnitsForModulus(10)) = %d, want %d", len(units), want)
	}
}
