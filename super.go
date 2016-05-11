// Copyright (c) 2016 Melvin Eloy Irizarry-Gelpí
// Licenced under the MIT License.

package dual

import (
	"fmt"
	"math"
	"strings"
)

// A Super represents a super dual number as an ordered array of two pointers
// to Real values.
type Super [2]*Real

var (
	// Symbols for the canonical super dual real basis.
	symbSuper = [4]string{"", "σ", "τ", "στ"}
)

// Real returns the real part of z, a pointer to a Real value.
func (z *Super) Real() *Real {
	return z[0]
}

// Dual returns the dual part of z, a pointer to a Real value.
func (z *Super) Dual() *Real {
	return z[1]
}

// SetReal sets the real part of z equal to a.
func (z *Super) SetReal(a *Real) {
	z[0] = a
}

// SetDual sets the dual part of z equal to b.
func (z *Super) SetDual(b *Real) {
	z[1] = b
}

// Cartesian returns the four Cartesian components of z.
func (z *Super) Cartesian() (a, b, c, d float64) {
	a, b = z.Real().Cartesian()
	c, d = z.Dual().Cartesian()
	return
}

// String returns the string representation of a Super value.
//
// If z corresponds to the super dual real number a + bσ + cτ + dστ, then the
// string is "(a+bσ+cτ+dστ)", similar to complex128 values.
func (z *Super) String() string {
	v := make([]float64, 4)
	v[0], v[1], v[2], v[3] = z.Cartesian()
	a := make([]string, 9)
	a[0] = "("
	a[1] = fmt.Sprintf("%g", v[0])
	i := 1
	for j := 2; j < 8; j = j + 2 {
		switch {
		case math.Signbit(v[i]):
			a[j] = fmt.Sprintf("%g", v[i])
		case math.IsInf(v[i], +1):
			a[j] = "+Inf"
		default:
			a[j] = fmt.Sprintf("+%g", v[i])
		}
		a[j+1] = symbSuper[i]
		i++
	}
	a[8] = ")"
	return strings.Join(a, "")
}

// Equals returns true if z and y are equal.
func (z *Super) Equals(y *Super) bool {
	if !z.Real().Equals(y.Real()) || !z.Dual().Equals(y.Dual()) {
		return false
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Super) Copy(y *Super) *Super {
	z.SetReal(new(Real).Copy(y.Real()))
	z.SetDual(new(Real).Copy(y.Dual()))
	return z
}

// NewSuper returns a pointer to a Super value made from four given float64
// values.
func NewSuper(a, b, c, d float64) *Super {
	z := new(Super)
	z.SetReal(NewReal(a, b))
	z.SetDual(NewReal(c, d))
	return z
}

// IsInf returns true if any of the components of z are infinite.
func (z *Super) IsInf() bool {
	if z.Real().IsInf() || z.Dual().IsInf() {
		return true
	}
	return false
}

// SuperInf returns a pointer to a super dual infinity value.
func SuperInf(a, b, c, d int) *Super {
	z := new(Super)
	z.SetReal(RealInf(a, b))
	z.SetDual(RealInf(c, d))
	return z
}

// IsNaN returns true if any component of z is NaN and neither is an
// infinity.
func (z *Super) IsNaN() bool {
	if z.Real().IsInf() || z.Dual().IsInf() {
		return false
	}
	if z.Real().IsNaN() || z.Dual().IsNaN() {
		return true
	}
	return false
}

// SuperNaN returns a pointer to a super dual NaN value.
func SuperNaN() *Super {
	z := new(Super)
	z.SetReal(RealNaN())
	z.SetDual(RealNaN())
	return z
}

// Scal sets z equal to y scaled by a (with a being a Real pointer),
// and returns z.
//
// This is a special case of Mul:
// 		Scal(y, a) = Mul(y, Super{a, 0})
func (z *Super) Scal(y *Super, a *Real) *Super {
	z.SetReal(new(Real).Mul(y.Real(), a))
	z.SetDual(new(Real).Mul(y.Dual(), a))
	return z
}

// Dil sets z equal to the dilation of y by a, and returns z.
//
// This is a special case of Mul:
// 		Dil(y, a) = Mul(y, Super{Real{a, 0}, 0})
func (z *Super) Dil(y *Super, a float64) *Super {
	z.SetReal(new(Real).Scal(y.Real(), a))
	z.SetDual(new(Real).Scal(y.Dual(), a))
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Super) Neg(y *Super) *Super {
	return z.Dil(y, -1)
}

// Conj sets z equal to the conjugate of y, and returns z.
func (z *Super) Conj(y *Super) *Super {
	z.SetReal(new(Real).Conj(y.Real()))
	z.SetDual(new(Real).Neg(y.Dual()))
	return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Super) Add(x, y *Super) *Super {
	z.SetReal(new(Real).Add(x.Real(), y.Real()))
	z.SetDual(new(Real).Add(x.Dual(), y.Dual()))
	return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Super) Sub(x, y *Super) *Super {
	z.SetReal(new(Real).Sub(x.Real(), y.Real()))
	z.SetDual(new(Real).Sub(x.Dual(), y.Dual()))
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
	z.SetReal(new(Real).Mul(p.Real(), q.Real()))
	z.SetDual(new(Real).Add(
		new(Real).Mul(q.Dual(), p.Real()),
		new(Real).Mul(p.Dual(), q.Real().Conj(q.Real()))))
	return z
}

// Commutator sets z equal to the commutator of x and y, and returns z.
func (z *Super) Commutator(x, y *Super) *Super {
	return z.Sub(new(Super).Mul(x, y), new(Super).Mul(y, x))
}

// Quad returns the dual quadrance of z, a float64 value.
func (z *Super) Quad() float64 {
	a := z.Real().Real()
	return a * a
}
