package intarith

import "sync"

// Integer arithmetic: gcd, extended gcd, lcm, totient, modular exponentiation.

var eulerPhiCache = make(map[int64]int64)
var eulerPhiMu sync.Mutex

func Gcd(a, b int64) int64 {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	for {
		r := a % b
		if r == 0 {
			break
		}
		a, b = b, r
	}
	if b < 0 {
		b = -b
	}
	return b
}

// ExtGcd returns (d, m, n) with d = a*m + b*n (Blankinship's algorithm).
func ExtGcd(a, b int64) (d, m, n int64) {
	mprime, n := int64(1), int64(1)
	m, nprime := int64(0), int64(0)
	c, d := a, b
	for {
		q, r := c/d, c%d
		if r == 0 {
			break
		}
		c, d = d, r
		t := mprime
		mprime = m
		m = t - q*m
		t = nprime
		nprime = n
		n = t - q*n
	}
	return d, m, n
}

func Lcm(a, b int64) int64 {
	return a * b / Gcd(a, b)
}

func EulerPhi(n int64) int64 {
	eulerPhiMu.Lock()
	defer eulerPhiMu.Unlock()
	if v, ok := eulerPhiCache[n]; ok {
		return v
	}
	if n <= 1 {
		return 0
	}
	phi := int64(0)
	for i := int64(1); i < n; i++ {
		if Gcd(n, i) == 1 {
			phi++
		}
	}
	eulerPhiCache[n] = phi
	return phi
}

func IntExp(x, e int64) (int64, error) {
	if e < 0 {
		return 0, &NegativeExponentError{E: e}
	}
	xp := x
	rv := int64(1)
	for e != 0 {
		if e&1 == 1 {
			rv = rv * xp
		}
		e >>= 1
		xp = xp * xp
	}
	return rv, nil
}

func IntModExp(x, e, m int64) int64 {
	if e < 0 {
		e = -e
		x = IntModRecip(x, m)
	}
	xp := x
	rv := int64(1)
	for e != 0 {
		if e&1 == 1 {
			rv = (rv * xp) % m
		}
		e >>= 1
		xp = (xp * xp) % m
	}
	return rv
}

func IntModRecip(x, m int64) int64 {
	if Gcd(x, m) != 1 {
		panic("intmodrecip: impossible inverse")
	}
	phi := EulerPhi(m)
	return IntModExp(x, phi-1, m)
}

func Factorial(n int64) (int64, error) {
	if n < 0 {
		return 0, &NegativeInputError{}
	}
	if n < 2 {
		return 1, nil
	}
	rv := int64(1)
	for k := int64(2); k <= n; k++ {
		rv *= k
	}
	return rv, nil
}

type NegativeExponentError struct{ E int64 }

func (e *NegativeExponentError) Error() string { return "intexp: negative exponent disallowed" }

type NegativeInputError struct{}

func (e *NegativeInputError) Error() string { return "factorial: negative input disallowed" }
