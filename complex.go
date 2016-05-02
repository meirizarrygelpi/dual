package dual

import (
	"fmt"
	"math"
	"math/cmplx"
	"strings"
)

// A Complex represents a dual complex number as an ordered array of two
// complex128 values.
type Complex [2]complex128

var (
	// Symbols for the canonical dual complex basis.
	symbComplex = [4]string{"", "i", "ε", "εi"}
)

// String returns the string representation of a Complex value. If z
// corresponds to the dual complex number a + bi + cε + dεi, then the string is
// "(a+bi+cε+dεi)", similar to complex128 values.
func (z *Complex) String() string {
	v := make([]float64, 4)
	v[0], v[1] = real(z[0]), imag(z[0])
	v[2], v[3] = real(z[1]), imag(z[1])
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
		a[j+1] = symbComplex[i]
		i++
	}
	a[8] = ")"
	return strings.Join(a, "")
}

// Equals returns true if z and y are equal.
func (z *Complex) Equals(y *Complex) bool {
	for i := range z {
		if notEquals(real(z[i]), real(y[i])) {
			return false
		}
		if notEquals(imag(z[i]), imag(y[i])) {
			return false
		}
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Complex) Copy(y *Complex) *Complex {
	z[0] = y[0]
	z[1] = y[1]
	return z
}

// NewComplex returns a pointer to a Complex value made from four given float64
// values.
func NewComplex(a, b, c, d float64) *Complex {
	z := new(Complex)
	z[0] = complex(a, b)
	z[1] = complex(c, d)
	return z
}

// IsComplexInf returns true if any of the components of z are infinite.
func (z *Complex) IsComplexInf() bool {
	if cmplx.IsInf(z[0]) || cmplx.IsInf(z[1]) {
		return true
	}
	return false
}

// ComplexInf returns a pointer to a dual complex infinity value.
func ComplexInf(a, b, c, d int) *Complex {
	z := new(Complex)
	z[0] = complex(math.Inf(a), math.Inf(b))
	z[1] = complex(math.Inf(c), math.Inf(d))
	return z
}

// IsComplexNaN returns true if any component of z is NaN and neither is an
// infinity.
func (z *Complex) IsComplexNaN() bool {
	if cmplx.IsInf(z[0]) || cmplx.IsInf(z[1]) {
		return false
	}
	if cmplx.IsNaN(z[0]) || cmplx.IsNaN(z[1]) {
		return true
	}
	return false
}

// ComplexNaN returns a pointer to a dual complex NaN value.
func ComplexNaN() *Complex {
	nan := cmplx.NaN()
	z := new(Complex)
	z[0] = nan
	z[1] = nan
	return z
}

// Scal sets z equal to y scaled by a (with a being a complex128), and returns
// z.
//
// This is a special case of Mul:
// 		Scal(y, a) = Mul(y, Complex{a, 0})
func (z *Complex) Scal(y *Complex, a complex128) *Complex {
	z[0] = y[0] * a
	z[1] = y[1] * a
	return z
}

// Dil sets z equal to the dilation of y by a, and returns z.
//
// This is a special case of Mul:
// 		Dil(y, a) = Mul(y, Complex{complex(a, 0), 0})
func (z *Complex) Dil(y *Complex, a float64) *Complex {
	z[0] = y[0] * complex(a, 0)
	z[1] = y[1] * complex(a, 0)
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Complex) Neg(y *Complex) *Complex {
	z[0] = -y[0]
	z[1] = -y[1]
	return z
}

// DualConj sets z equal to the dual conjugate of y, and returns z.
func (z *Complex) DualConj(y *Complex) *Complex {
	z[0] = +y[0]
	z[1] = -y[1]
	return z
}

// Conj sets z equal to the complex conjugate of y, and returns z.
func (z *Complex) Conj(y *Complex) *Complex {
	z[0] = cmplx.Conj(y[0])
	z[1] = cmplx.Conj(y[1])
	return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Complex) Add(x, y *Complex) *Complex {
	z[0] = x[0] + y[0]
	z[1] = x[1] + y[1]
	return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Complex) Sub(x, y *Complex) *Complex {
	z[0] = x[0] - y[0]
	z[1] = x[1] - y[1]
	return z
}

// Mul sets z equal to the product of x and y, and returns z.
//
// The basic rules are:
//      ε * ε = 0
//      i * i = -1
//      εi * εi = 0
//      ε * i = i * ε = εi
//      εi * i = i * εi = -ε
//      ε * εi = εi * ε = 0
// This multiplication rule is commutative and associative.
func (z *Complex) Mul(x, y *Complex) *Complex {
	p := new(Complex).Copy(x)
	q := new(Complex).Copy(y)
	z[0] = p[0] * q[0]
	z[1] = (p[0] * q[1]) + (p[1] * q[0])
	return z
}

// Quad returns the quadrance of z, a dual real value.
func (z *Complex) Quad() *Real {
	p := new(Complex).Mul(z, new(Complex).Conj(z))
	d := new(Real)
	d[0] = real(p[0])
	d[1] = real(p[1])
	return d
}

// DualQuad returns the dual quadrance of z, a complex128 value.
func (z *Complex) DualQuad() complex128 {
	return (new(Complex).Mul(z, new(Complex).DualConj(z)))[0]
}

// IsZeroDiv returns true if z is a zero divisor. This is equivalent to
// z being nilpotent (i.e. z² = 0).
func (z *Complex) IsZeroDiv() bool {
	return z[0] != complex(0, 0)
}

// Inv sets z equal to the inverse of y, and returns z. If y is a zero divisor,
// then Inv panics.
func (z *Complex) Inv(y *Complex) *Complex {
	if y.IsZeroDiv() {
		panic("zero divisor")
	}
	return z.Scal(new(Complex).DualConj(y), 1/y.DualQuad())
}

// Quo sets z equal to the quotient of x and y, and returns z. If y is a zero
// divisor, then Quo panics.
func (z *Complex) Quo(x, y *Complex) *Complex {
	if y.IsZeroDiv() {
		panic("zero divisor denominator")
	}
	return z.Scal(new(Complex).Mul(x, new(Complex).DualConj(y)), 1/y.DualQuad())
}
