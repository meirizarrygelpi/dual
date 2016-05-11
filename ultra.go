// Copyright (c) 2016 Melvin Eloy Irizarry-Gelpí
// Licenced under the MIT License.

package dual

import (
	"fmt"
	"math"
	"strings"
)

// An Ultra represents an ultra dual number as an ordered array of two pointers
// to Super values.
type Ultra [2]*Super

var (
	// Symbols for the canonical ultra dual real basis.
	symbUltra = [8]string{"", "υ₁", "υ₂", "υ₃", "υ₄", "υ₅", "υ₆", "υ₇"}
)

// Real returns the real part of z, a pointer to a Super value.
func (z *Ultra) Real() *Super {
	return z[0]
}

// Dual returns the dual part of z, a pointer to a Super value.
func (z *Ultra) Dual() *Super {
	return z[1]
}

// SetReal sets the real part of z equal to a.
func (z *Ultra) SetReal(a *Super) {
	z[0] = a
}

// SetDual sets the dual part of z equal to b.
func (z *Ultra) SetDual(b *Super) {
	z[1] = b
}

// Cartesian returns the four Cartesian components of z.
func (z *Ultra) Cartesian() (a, b, c, d, e, f, g, h float64) {
	a, b, c, d = z.Real().Cartesian()
	e, f, g, h = z.Dual().Cartesian()
	return
}

// String returns the string representation of a Ultra value.
//
// If z corresponds to the ultra dual real number a + bσ + cτ + dστ, then the
// string is "(a+bσ+cτ+dστ)", similar to complex128 values.
func (z *Ultra) String() string {
	v := make([]float64, 8)
	v[0], v[1], v[2], v[3] = z.Real().Cartesian()
	v[4], v[5], v[6], v[7] = z.Dual().Cartesian()
	a := make([]string, 17)
	a[0] = "("
	a[1] = fmt.Sprintf("%g", v[0])
	i := 1
	for j := 2; j < 16; j = j + 2 {
		switch {
		case math.Signbit(v[i]):
			a[j] = fmt.Sprintf("%g", v[i])
		case math.IsInf(v[i], +1):
			a[j] = "+Inf"
		default:
			a[j] = fmt.Sprintf("+%g", v[i])
		}
		a[j+1] = symbUltra[i]
		i++
	}
	a[16] = ")"
	return strings.Join(a, "")
}

// Equals returns true if z and y are equal.
func (z *Ultra) Equals(y *Ultra) bool {
	if !z.Real().Equals(y.Real()) || !z.Dual().Equals(y.Dual()) {
		return false
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Ultra) Copy(y *Ultra) *Ultra {
	z.SetReal(new(Super).Copy(y.Real()))
	z.SetDual(new(Super).Copy(y.Dual()))
	return z
}

// NewUltra returns a pointer to a Ultra value made from eight given float64
// values.
func NewUltra(a, b, c, d, e, f, g, h float64) *Ultra {
	z := new(Ultra)
	z.SetReal(NewSuper(a, b, c, d))
	z.SetDual(NewSuper(e, f, g, h))
	return z
}

// IsInf returns true if any of the components of z are infinite.
func (z *Ultra) IsInf() bool {
	if z.Real().IsInf() || z.Dual().IsInf() {
		return true
	}
	return false
}

// UltraInf returns a pointer to an ultra dual infinity value.
func UltraInf(a, b, c, d, e, f, g, h int) *Ultra {
	z := new(Ultra)
	z.SetReal(SuperInf(a, b, c, d))
	z.SetDual(SuperInf(e, f, g, h))
	return z
}

// IsNaN returns true if any component of z is NaN and neither is an
// infinity.
func (z *Ultra) IsNaN() bool {
	if z.Real().IsInf() || z.Dual().IsInf() {
		return false
	}
	if z.Real().IsNaN() || z.Dual().IsNaN() {
		return true
	}
	return false
}

// UltraNaN returns a pointer to an ultra dual NaN value.
func UltraNaN() *Ultra {
	z := new(Ultra)
	z.SetReal(SuperNaN())
	z.SetDual(SuperNaN())
	return z
}

// Scal sets z equal to y scaled by a (with a being a Super pointer),
// and returns z.
//
// This is a special case of Mul:
// 		Scal(y, a) = Mul(y, Ultra{a, 0})
func (z *Ultra) Scal(y *Ultra, a *Super) *Ultra {
	z.SetReal(new(Super).Mul(y.Real(), a))
	z.SetDual(new(Super).Mul(y.Dual(), a))
	return z
}

// Dil sets z equal to the dilation of y by a, and returns z.
//
// This is a special case of Mul:
// 		Dil(y, a) = Mul(y, Ultra{Super{a, 0}, 0})
func (z *Ultra) Dil(y *Ultra, a float64) *Ultra {
	z.SetReal(new(Super).Dil(y.Real(), a))
	z.SetDual(new(Super).Dil(y.Dual(), a))
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Ultra) Neg(y *Ultra) *Ultra {
	return z.Dil(y, -1)
}

// Conj sets z equal to the conjugate of y, and returns z.
func (z *Ultra) Conj(y *Ultra) *Ultra {
	z.SetReal(new(Super).Conj(y.Real()))
	z.SetDual(new(Super).Neg(y.Dual()))
	return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Ultra) Add(x, y *Ultra) *Ultra {
	z.SetReal(new(Super).Add(x.Real(), y.Real()))
	z.SetDual(new(Super).Add(x.Dual(), y.Dual()))
	return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Ultra) Sub(x, y *Ultra) *Ultra {
	z.SetReal(new(Super).Sub(x.Real(), y.Real()))
	z.SetDual(new(Super).Sub(x.Dual(), y.Dual()))
	return z
}

// Mul sets z equal to the product of x and y, and returns z.
//
// The basic multiplication rules are:
//      υ₁ * υ₂ = -υ₂ * υ₁ = υ₃
// 		υ₁ * υ₄ = -υ₄ * υ₁ = υ₅
// 		υ₂ * υ₄ = -υ₄ * υ₂ = υ₆
//		υ₂ * υ₅ = -υ₅ * υ₂ = υ₇
// 		υ₃ * υ₄ = -υ₄ * υ₃ = υ₇
// 		υ₆ * υ₁ = -υ₁ * υ₆ = υ₇
// All other products vanish. This multiplication operation is noncommutative
// and nonassociative.
func (z *Ultra) Mul(x, y *Ultra) *Ultra {
	p := new(Ultra).Copy(x)
	q := new(Ultra).Copy(y)
	z.SetReal(new(Super).Mul(p.Real(), q.Real()))
	z.SetDual(new(Super).Add(
		new(Super).Mul(q.Dual(), p.Real()),
		new(Super).Mul(p.Dual(), q.Real().Conj(q.Real())),
	))
	return z
}

// Commutator sets z equal to the commutator of x and y, and returns z.
func (z *Ultra) Commutator(x, y *Ultra) *Ultra {
	return z.Sub(new(Ultra).Mul(x, y), new(Ultra).Mul(y, x))
}

// Associator sets z equal to the associator of w, x, and y, and returns z.
func (z *Ultra) Associator(w, x, y *Ultra) *Ultra {
	return z.Sub(
		new(Ultra).Mul(new(Ultra).Mul(w, x), y),
		new(Ultra).Mul(w, new(Ultra).Mul(x, y)),
	)
}

// Quad returns the quadrance of z, a float64 value.
func (z *Ultra) Quad() float64 {
	a := z.Real().Real().Real()
	return a * a
}
