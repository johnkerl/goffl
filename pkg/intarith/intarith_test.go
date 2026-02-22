package intarith

import "testing"

func TestGcd(t *testing.T) {
	if got := Gcd(0, 0); got != 0 {
		t.Errorf("Gcd(0,0) = %d, want 0", got)
	}
	if got := Gcd(24, 60); got != 12 {
		t.Errorf("Gcd(24,60) = %d, want 12", got)
	}
	if got := Gcd(24, 65); got != 1 {
		t.Errorf("Gcd(24,65) = %d, want 1", got)
	}
}

func TestExtGcd(t *testing.T) {
	d, m, n := ExtGcd(1, 1)
	if d != 1 || d != 1*m+1*n {
		t.Errorf("ExtGcd(1,1) = (%d,%d,%d)", d, m, n)
	}
	d, m, n = ExtGcd(24, 65)
	if d != 1 || 24*m+65*n != 1 {
		t.Errorf("ExtGcd(24,65) = (%d,%d,%d), want d=1 and 24*m+65*n=1", d, m, n)
	}
	d, m, n = ExtGcd(24, 60)
	if d != 12 || m != -2 || n != 1 {
		t.Errorf("ExtGcd(24,60) = (%d,%d,%d), want (12,-2,1)", d, m, n)
	}
	// b == 0 and a == 0
	d, m, n = ExtGcd(7, 0)
	if d != 7 || m != 1 || n != 0 {
		t.Errorf("ExtGcd(7,0) = (%d,%d,%d), want (7,1,0)", d, m, n)
	}
	d, m, n = ExtGcd(0, 11)
	if d != 11 || m != 0 || n != 1 {
		t.Errorf("ExtGcd(0,11) = (%d,%d,%d), want (11,0,1)", d, m, n)
	}
	d, m, n = ExtGcd(0, 0)
	if d != 0 || m != 0 || n != 0 {
		t.Errorf("ExtGcd(0,0) = (%d,%d,%d), want (0,0,0)", d, m, n)
	}
}

func TestEulerPhi(t *testing.T) {
	tests := []struct {
		n, want int64
	}{
		{0, 0}, {1, 0}, {2, 1}, {5, 4}, {7, 6},
	}
	for _, tt := range tests {
		if got := EulerPhi(tt.n); got != tt.want {
			t.Errorf("EulerPhi(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestIntModExp(t *testing.T) {
	if got := IntModExp(2, 0, 11); got != 1 {
		t.Errorf("IntModExp(2,0,11) = %d, want 1", got)
	}
	if got := IntModExp(2, 10, 11); got != 1 {
		t.Errorf("IntModExp(2,10,11) = %d, want 1", got)
	}
	if got := IntModExp(2, -1, 11); got != 6 {
		t.Errorf("IntModExp(2,-1,11) = %d, want 6", got)
	}
}

func TestLcm(t *testing.T) {
	if got := Lcm(4, 6); got != 12 {
		t.Errorf("Lcm(4,6) = %d, want 12", got)
	}
	if got := Lcm(0, 6); got != 0 {
		t.Errorf("Lcm(0,6) = %d, want 0", got)
	}
	if got := Lcm(4, 0); got != 0 {
		t.Errorf("Lcm(4,0) = %d, want 0", got)
	}
}
