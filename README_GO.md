
- **Bit arithmetic**: `bitarith` (Msb, Lsb, Ones, MsbPos, LsbPos, etc.), `bitvector`, `bitmatrix` (row echelon, kernel basis over GF(2)).
- **Integer arithmetic**: `intarith` (Gcd, ExtGcd, Lcm, EulerPhi, IntModExp), `intmod` (integers mod *n*), `factorization`, `intfactor` (trial division, Totient).
- **Polynomials over GF(2)**: `f2poly` (bits as coefficients), `f2polymod` (quotient ring), `f2polyfactor` (Berlekamp factorization, Irr, Totient).
- **Orders**: `order` (ModOrderIntMod, ModOrderF2PolyMod, orbit, period, generators, F2PolyPrimitive).

## Install

From the project root:

```bash
go build ./...
go test ./...
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/johnkerl/goffl/pkg/intarith"
	"github.com/johnkerl/goffl/pkg/intmod"
	"github.com/johnkerl/goffl/pkg/intfactor"
	"github.com/johnkerl/goffl/pkg/f2poly"
	"github.com/johnkerl/goffl/pkg/f2polymod"
	"github.com/johnkerl/goffl/pkg/f2polyfactor"
	"github.com/johnkerl/goffl/pkg/order"
)

func main() {
	// Integer arithmetic
	fmt.Println(intarith.Gcd(24, 60)) // 12
	d, m, n := intarith.ExtGcd(24, 65)
	fmt.Println(d, 24*m+65*n == 1) // 1 true
	a := intmod.New(2, 11)
	fmt.Println(a.Pow(10).Residue) // 1

	// Integer factorization
	finfo := intfactor.Factor(72)
	fmt.Println(finfo.Unfactor())       // 72
	fmt.Println(finfo.AllDivisors())    // [1 2 3 4 6 8 9 12 18 24 36 72]

	// F2 polynomials (hex or int: x^4 + x + 1 = 0x13)
	f := f2poly.New(0x13)
	fmt.Println(f.Degree()) // 4
	_ = f2polyfactor.Factor(f) // irreducible
	fmt.Println(f2polyfactor.Irr(f)) // true

	// Multiplicative order in F2[x]/(m)
	m := f2poly.New(0x11b) // AES field poly
	x := f2polymod.NewFromInts(2, m.Bits)
	ord, _ := order.ModOrderF2PolyMod(x)
	fmt.Println(ord) // 255
}
```

## License

See [LICENSE.txt](LICENSE.txt). Redistribution and use in source and binary forms, with or without modification, are permitted under the terms of the BSD-style license.
