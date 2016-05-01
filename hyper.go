package dual

import (
	"fmt"
	"math"
	"strings"
)

// A Hyper represents a hyper dual real number as an ordered array of four
// float64 values
type Hyper [4]float64

var (
	symbH = [4]string{"", "ε", "η", "εη"}
)

// String returns the string representation of a Hyper value. If z corresponds
// to the hyper dual real number a + bε + cη + dεη, then the string is
// "(a+bε+cη+dεη)", similar to complex128 values.
func (z *Hyper) String() string {
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
		a[j+1] = symbH[i]
		i++
	}
	a[8] = ")"
	return strings.Join(a, "")
}

// Equals returns true if z and y are equal.
func (z *Hyper) Equals(y *Hyper) bool {
	for i := range z {
		if notEquals(z[i], y[i]) {
			return false
		}
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Hyper) Copy(y *Hyper) *Hyper {
	for i, v := range y {
		z[i] = v
	}
	return z
}

// NewHyper returns a pointer to a Hyper value made from four given float64
// values.
func NewHyper(a, b, c, d float64) *Hyper {
	z := new(Hyper)
	z[0] = a
	z[1] = b
	z[2] = c
	z[3] = d
	return z
}

// IsHyperInf returns true if any of the components of z are infinite.
func (z *Hyper) IsHyperInf() bool {
	for _, v := range z {
		if math.IsInf(v, 0) {
			return true
		}
	}
	return false
}

// HyperInf returns a pointer to a hyper dual real infinity value.
func HyperInf(a, b, c, d int) *Hyper {
	z := new(Hyper)
	z[0] = math.Inf(a)
	z[1] = math.Inf(b)
	z[2] = math.Inf(c)
	z[3] = math.Inf(d)
	return z
}

// IsHyperNaN returns true if any component of z is NaN and neither is an
// infinity.
func (z *Hyper) IsHyperNaN() bool {
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

// HyperNaN returns a pointer to a hyper dual real NaN value.
func HyperNaN() *Hyper {
	nan := math.NaN()
	return &Hyper{nan, nan, nan, nan}
}

// Scal sets z equal to y scaled by a, and returns z.
func (z *Hyper) Scal(y *Hyper, a float64) *Hyper {
	for i, v := range y {
		z[i] = a * v
	}
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Hyper) Neg(y *Hyper) *Hyper {
	return z.Scal(y, -1)
}

// DConj sets z equal to the hyper dual conjugate of y, and returns z.
func (z *Hyper) DConj(y *Hyper) *Hyper {
	z[0] = +y[0]
	z[1] = -y[1]
	z[2] = -y[2]
	z[3] = -y[3]
	return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Hyper) Add(x, y *Hyper) *Hyper {
	for i, v := range x {
		z[i] = v + y[i]
	}
	return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Hyper) Sub(x, y *Hyper) *Hyper {
	for i, v := range x {
		z[i] = v - y[i]
	}
	return z
}

// Mul sets z equal to the product of x and y, and returns z.
//
// The basic rules are:
//      ε * ε = η * η = 0
//      ε * η = εη
//      η * ε = -εη
//      εη * εη = 0
//      ε * εη = εη * ε = 0
//      η * εη = εη * η = 0
// Note that this multiplication operation is noncommutative.
func (z *Hyper) Mul(x, y *Hyper) *Hyper {
	p := new(Hyper).Copy(x)
	q := new(Hyper).Copy(y)
	z[0] = (p[0] * q[0])
	z[1] = (p[0] * q[1]) + (p[1] * q[0])
	z[2] = (p[0] * q[2]) + (p[2] * q[0])
	z[3] = (p[0] * q[3]) + (p[1] * q[2]) - (p[2] * q[1]) + (p[3] * q[0])
	return z
}

// DQuad returns the hyper dual quadrance of z, a float64 value.
func (z *Hyper) DQuad() float64 {
	return z[0] * z[0]
}

// Commutator sets z equal to the commutator of x and y, and returns z.
func (z *Hyper) Commutator(x, y *Hyper) *Hyper {
	return z.Sub(new(Hyper).Mul(x, y), new(Hyper).Mul(y, x))
}
