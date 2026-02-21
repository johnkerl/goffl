package factorization

import (
	"fmt"
	"sort"
	"strings"
)

// Factorization stores a trivial factor and list of (factor, multiplicity) for integers.
type Factorization struct {
	trivialFactor *int64
	factors       []struct {
		factor int64
		mult   int
	}
}

func New() *Factorization {
	return &Factorization{factors: nil}
}

func (f *Factorization) TrivialFactor() (int64, bool) {
	if f.trivialFactor == nil {
		return 0, false
	}
	return *f.trivialFactor, true
}

func (f *Factorization) NumDistinctFactors() int { return len(f.factors) }

func (f *Factorization) NumFactors() int {
	n := 0
	for _, p := range f.factors {
		n += p.mult
	}
	return n
}

func (f *Factorization) Get(i int) (factor int64, mult int) {
	return f.factors[i].factor, f.factors[i].mult
}

func (f *Factorization) InsertTrivialFactor(n *int64) {
	if n == nil {
		return
	}
	if f.trivialFactor != nil {
		*f.trivialFactor *= *n
	} else {
		x := *n
		f.trivialFactor = &x
	}
}

func (f *Factorization) InsertFactor(newFactor int64, newMult int) {
	if newMult <= 0 {
		return
	}
	for i := range f.factors {
		if f.factors[i].factor == newFactor {
			f.factors[i].mult += newMult
			return
		}
		if newFactor < f.factors[i].factor {
			f.factors = append(f.factors, struct {
				factor int64
				mult   int
			}{})
			copy(f.factors[i+1:], f.factors[i:])
			f.factors[i] = struct {
				factor int64
				mult   int
			}{newFactor, newMult}
			return
		}
	}
	f.factors = append(f.factors, struct {
		factor int64
		mult   int
	}{newFactor, newMult})
}

func (f *Factorization) Merge(other *Factorization) {
	if t, ok := other.TrivialFactor(); ok {
		f.InsertTrivialFactor(&t)
	}
	for i := 0; i < other.NumDistinctFactors(); i++ {
		p, m := other.Get(i)
		f.InsertFactor(p, m)
	}
}

func (f *Factorization) ExpAll(e int) {
	if f.trivialFactor != nil {
		exp := int64(1)
		for i := 0; i < e; i++ {
			exp *= *f.trivialFactor
		}
		*f.trivialFactor = exp
	}
	for i := range f.factors {
		f.factors[i].mult *= e
	}
}

func (f *Factorization) NumDivisors() int {
	ndf := f.NumDistinctFactors()
	if ndf <= 0 {
		if f.trivialFactor == nil {
			panic("num_divisors: no factors have been inserted")
		}
	}
	rv := 1
	for i := 0; i < ndf; i++ {
		_, m := f.Get(i)
		rv *= m + 1
	}
	return rv
}

func (f *Factorization) KthDivisor(k int) int64 {
	ndf := f.NumDistinctFactors()
	if ndf <= 0 {
		if f.trivialFactor != nil {
			return 1
		}
		panic("kth_divisor: no factors have been inserted")
	}
	rv := int64(1)
	for i := 0; i < ndf; i++ {
		p, m := f.Get(i)
		base := m + 1
		power := k % base
		k = k / base
		for j := 0; j < power; j++ {
			rv *= p
		}
	}
	return rv
}

func (f *Factorization) AllDivisors() []int64 {
	ndf := f.NumDistinctFactors()
	if ndf <= 0 && f.trivialFactor == nil {
		panic("all_divisors: no factors have been inserted")
	}
	nd := f.NumDivisors()
	out := make([]int64, nd)
	for k := 0; k < nd; k++ {
		out[k] = f.KthDivisor(k)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func (f *Factorization) MaximalProperDivisors() []int64 {
	ndf := f.NumDistinctFactors()
	if ndf <= 0 {
		if f.trivialFactor == nil {
			panic("maximal_proper_divisors: no factors have been inserted")
		}
		return nil
	}
	n := f.Unfactor()
	out := make([]int64, ndf)
	for k := 0; k < ndf; k++ {
		p, _ := f.Get(k)
		out[k] = n / p
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func (f *Factorization) Unfactor() int64 {
	ndf := f.NumDistinctFactors()
	if ndf <= 0 {
		if f.trivialFactor == nil {
			panic("unfactor: no factors have been inserted")
		}
		return *f.trivialFactor
	}
	rv := int64(1)
	if f.trivialFactor != nil {
		rv = *f.trivialFactor
	}
	for i := 0; i < ndf; i++ {
		p, e := f.Get(i)
		for j := 0; j < e; j++ {
			rv *= p
		}
	}
	return rv
}

func (f *Factorization) String() string {
	var parts []string
	if f.trivialFactor != nil {
		parts = append(parts, fmt.Sprint(*f.trivialFactor))
	}
	for i := 0; i < f.NumDistinctFactors(); i++ {
		p, m := f.Get(i)
		s := fmt.Sprint(p)
		if m != 1 {
			s += fmt.Sprintf("^%d", m)
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, " ")
}
