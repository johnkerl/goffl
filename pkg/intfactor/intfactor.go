// Package intfactor provides integer factorization and Euler totient for int64.
package intfactor

import (
	"github.com/johnkerl/goffl/pkg/factorization"
	"github.com/johnkerl/goffl/pkg/intarith"
)

func Factor(n int64) *factorization.Factorization {
	finfo := factorization.New()
	if n >= -1 && n <= 1 {
		finfo.InsertTrivialFactor(&n)
		return finfo
	}
	if n < 0 {
		minusOne := int64(-1)
		finfo.InsertTrivialFactor(&minusOne)
		n = -n
	}
	p := int64(2)
	for n > 1 {
		multiplicity := 0
		for n%p == 0 {
			multiplicity++
			n /= p
		}
		if multiplicity > 0 {
			finfo.InsertFactor(p, multiplicity)
		}
		if p > 2 {
			p += 2
		} else {
			p += 1
		}
	}
	return finfo
}

func SlowTotient(n int64) int64 {
	var count int64
	for a := int64(1); a < n; a++ {
		if intarith.Gcd(a, n) == 1 {
			count++
		}
	}
	return count
}

func Totient(n int64) int64 {
	finfo := Factor(n)
	rv := n
	for i := 0; i < finfo.NumDistinctFactors(); i++ {
		p, e := finfo.Get(i)
		for j := 0; j < e; j++ {
			rv = rv * (p - 1) / p
		}
	}
	return rv
}
