package main

import (
	"fmt"
	"strconv"

	"github.com/johnkerl/goffl/pkg/intmod"
)

// IntModNumeric implements Numeric[*intmod.IntMod, int] using goffl intmod.
// The modulus is fixed when the backend is created (e.g. from -mod flag).
type IntModNumeric struct {
	Modulus int64
}

// NewIntModNumeric creates a backend for Z/nZ. n must be positive.
func NewIntModNumeric(n int64) (*IntModNumeric, error) {
	if n <= 0 {
		return nil, fmt.Errorf("modulus must be positive, got %d", n)
	}
	return &IntModNumeric{Modulus: n}, nil
}

func (b *IntModNumeric) FromString(s string) (*intmod.IntMod, error) {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return nil, err
	}
	return intmod.New(v, b.Modulus), nil
}

func (b *IntModNumeric) ParseExponent(s string) (int, error) {
	v, err := strconv.ParseInt(s, 0, 64)
	return int(v), err
}

func (b *IntModNumeric) String(t *intmod.IntMod) string {
	return strconv.FormatInt(t.Residue, 10)
}

func (b *IntModNumeric) Add(a, c *intmod.IntMod) *intmod.IntMod      { return a.Add(c) }
func (b *IntModNumeric) Subtract(a, c *intmod.IntMod) *intmod.IntMod { return a.Sub(c) }
func (b *IntModNumeric) Multiply(a, c *intmod.IntMod) *intmod.IntMod { return a.Mul(c) }

func (b *IntModNumeric) Divide(a, c *intmod.IntMod) (*intmod.IntMod, error) {
	if c.Residue == 0 {
		return nil, fmt.Errorf("division by zero")
	}
	return a.Div(c)
}

func (b *IntModNumeric) Mod(a, c *intmod.IntMod) (*intmod.IntMod, error) {
	return nil, fmt.Errorf("modulo not defined for Z/nZ (use integer or float mode for %%)")
}

func (b *IntModNumeric) Exponentiate(base *intmod.IntMod, exp int) (*intmod.IntMod, error) {
	return base.Pow(int64(exp))
}

func (b *IntModNumeric) ToExponent(v *intmod.IntMod) (int, error) {
	return int(v.Residue), nil
}

func (b *IntModNumeric) Negate(v *intmod.IntMod) *intmod.IntMod {
	return v.Neg()
}
