package intfactor

import (
	"reflect"
	"testing"
)

func TestFactor(t *testing.T) {
	finfo := Factor(72)
	if got := finfo.Unfactor(); got != 72 {
		t.Errorf("Unfactor() = %d, want 72", got)
	}
	divs := finfo.AllDivisors()
	want := []int64{1, 2, 3, 4, 6, 8, 9, 12, 18, 24, 36, 72}
	if !reflect.DeepEqual(divs, want) {
		t.Errorf("AllDivisors() = %v, want %v", divs, want)
	}
}

func TestTotient(t *testing.T) {
	if got := Totient(1); got != 1 {
		t.Errorf("Totient(1) = %d, want 1", got)
	}
	if got := Totient(7); got != 6 {
		t.Errorf("Totient(7) = %d, want 6", got)
	}
	if got := Totient(10); got != 4 {
		t.Errorf("Totient(10) = %d, want 4", got)
	}
}
