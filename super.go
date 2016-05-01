package dual

import (
	"fmt"
	"math"
	"strings"
)

// A Super represents a super dual real number as an ordered array of four
// float64 values
type Super [4]float64

var (
	// Symbols for the canonical super dual real basis.
	symbSuper = [4]string{"", "σ", "τ", "στ"}
)

// String returns the string representation of a Super value. If z corresponds
// to the super dual real number a + bσ + cτ + dστ, then the string is
// "(a+bσ+cτ+dστ)", similar to complex128 values.
func (z *Super) String() string {
	a := make([]string, 9)
	a[0] = "("
	a[1] = fmt.Sprintf("%g", z[0])
	i := 1
	for j := 2; j < 8; j = j + 2 {
		switch {
		case math.Signbit(z[i]):
			a[j] = fmt.Sprintf("%g", z[i])
		case math.IsInf(z[i], +1):
			a[j] = "+Inf"
		default:
			a[j] = fmt.Sprintf("+%g", z[i])
		}
		a[j+1] = symbSuper[i]
		i++
	}
	a[8] = ")"
	return strings.Join(a, "")
}

// Equals returns true if z and y are equal.
func (z *Super) Equals(y *Super) bool {
	for i := range z {
		if notEquals(z[i], y[i]) {
			return false
		}
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Super) Copy(y *Super) *Super {
	for i, v := range y {
		z[i] = v
	}
	return z
}

// NewSuper returns a pointer to a Super value made from four given float64
// values.
func NewSuper(a, b, c, d float64) *Super {
	z := new(Super)
	z[0] = a
	z[1] = b
	z[2] = c
	z[3] = d
	return z
}

// IsSuperInf returns true if any of the components of z are infinite.
func (z *Super) IsSuperInf() bool {
	for _, v := range z {
		if math.IsInf(v, 0) {
			return true
		}
	}
	return false
}

// SuperInf returns a pointer to a super dual real infinity value.
func SuperInf(a, b, c, d int) *Super {
	z := new(Super)
	z[0] = math.Inf(a)
	z[1] = math.Inf(b)
	z[2] = math.Inf(c)
	z[3] = math.Inf(d)
	return z
}

// IsSuperNaN returns true if any component of z is NaN and neither is an
// infinity.
func (z *Super) IsSuperNaN() bool {
	for _, v := range z {
		if math.IsInf(v, 0) {
			return false
		}
	}
	for _, v := range z {
		if math.IsNaN(v) {
			return true
		}
	}
	return false
}

// SuperNaN returns a pointer to a super dual real NaN value.
func SuperNaN() *Super {
	nan := math.NaN()
	return &Super{nan, nan, nan, nan}
}

// Scal sets z equal to y scaled by a, and returns z.
func (z *Super) Scal(y *Super, a float64) *Super {
	for i, v := range y {
		z[i] = a * v
	}
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Super) Neg(y *Super) *Super {
	return z.Scal(y, -1)
}

// DualConj sets z equal to the super dual conjugate of y, and returns z.
func (z *Super) DualConj(y *Super) *Super {
	z[0] = +y[0]
	z[1] = -y[1]
	z[2] = -y[2]
	z[3] = -y[3]
	return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Super) Add(x, y *Super) *Super {
	for i, v := range x {
		z[i] = v + y[i]
	}
	return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Super) Sub(x, y *Super) *Super {
	for i, v := range x {
		z[i] = v - y[i]
	}
	return z
}

// Mul sets z equal to the product of x and y, and returns z.
//
// The basic multiplication rules are:
//      σ * σ = τ * τ = 0
//      σ * τ = -τ * σ = στ
//      στ * στ = 0
//      σ * στ = στ * σ = 0
//      τ * στ = στ * τ = 0
// This multiplication operation is noncommutative but associative.
func (z *Super) Mul(x, y *Super) *Super {
	p := new(Super).Copy(x)
	q := new(Super).Copy(y)
	z[0] = (p[0] * q[0])
	z[1] = (p[0] * q[1]) + (p[1] * q[0])
	z[2] = (p[0] * q[2]) + (p[2] * q[0])
	z[3] = (p[0] * q[3]) + (p[1] * q[2]) - (p[2] * q[1]) + (p[3] * q[0])
	return z
}

// Commutator sets z equal to the commutator of x and y, and returns z.
func (z *Super) Commutator(x, y *Super) *Super {
	return z.Sub(new(Super).Mul(x, y), new(Super).Mul(y, x))
}

// DualQuad returns the super dual quadrance of z, a float64 value.
func (z *Super) DualQuad() float64 {
	return z[0] * z[0]
}

// IsZeroDiv returns true if z is a zero divisor. This is equivalent to
// z being nilpotent (i.e. z² = 0).
func (z *Super) IsZeroDiv() bool {
	return !notEquals(z[0], 0)
}

// Inv sets z equal to the inverse of y, and returns z. If y is a zero divisor,
// then Inv panics.
func (z *Super) Inv(y *Super) *Super {
	if y.IsZeroDiv() {
		panic("zero divisor")
	}
	return z.Scal(new(Super).DualConj(y), 1/y.DualQuad())
}

// Quo sets z equal to the quotient of x and y, and returns z. If y is a zero
// divisor, then Quo panics.
func (z *Super) Quo(x, y *Super) *Super {
	if y.IsZeroDiv() {
		panic("zero divisor denominator")
	}
	return z.Scal(new(Super).Mul(x, new(Super).DualConj(y)), 1/y.DualQuad())
}
