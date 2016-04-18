package dual

import (
	"fmt"
	"strings"
)

// Dual type represents a dual number over the real numbers.
type Dual [2]float64

// String method returns the string version of a Dual value. It mimics the
// behavior for complex64 and complex128 types.
func (z Dual) String() string {
	a := make([]string, 5)
	a[0] = "("
	a[1] = fmt.Sprintf("%g", z[0])
	if z[1] < 0 {
		a[2] = fmt.Sprintf("%g", z[1])
	} else {
		a[2] = fmt.Sprintf("+%g", z[1])
	}
	a[3] = "ε"
	a[4] = ")"
	return strings.Join(a, "")
}

// Equals method returns true if z and x are equal.
func (z *Dual) Equals(x *Dual) bool {
	for i := range z {
		if notEquals(z[i], x[i]) {
			return false
		}
	}
	return true
}

// Copy method copies x onto z.
func (z *Dual) Copy(x *Dual) *Dual {
	for i, v := range x {
		z[i] = v
	}
	return z
}

// New function returns a pointer to a Dual value made from two given real
// numbers (i.e. float64s): a + bε.
func New(a, b float64) *Dual {
	z := new(Dual)
	z[0] = a
	z[1] = b
	return z
}

// Scal method sets z equal to a*x, and returns z.
func (z *Dual) Scal(x *Dual, a float64) *Dual {
	for i, v := range x {
		z[i] = a * v
	}
	return z
}

// Neg method sets z equal to the negative of x, and returns z.
func (z *Dual) Neg(x *Dual) *Dual {
	return z.Scal(x, -1)
}

// Conj method sets z equal to the conjugate of x, and returns z.
func (z *Dual) Conj(x *Dual) *Dual {
	z[0] = +x[0]
	z[1] = -x[1]
	return z
}

// Add method sets z to the sum of x and y, and returns z.
func (z *Dual) Add(x, y *Dual) *Dual {
	for i, v := range x {
		z[i] = v + y[i]
	}
	return z
}

// Sub method sets z to the difference of x and y, and returns z.
func (z *Dual) Sub(x, y *Dual) *Dual {
	for i, v := range x {
		z[i] = v - y[i]
	}
	return z
}

// Mul method sets z to the product of x and y, and returns z.
func (z *Dual) Mul(x, y *Dual) *Dual {
	p := new(Dual).Copy(x)
	q := new(Dual).Copy(y)
	z[0] = p[0] * q[0]
	z[1] = (p[0] * q[1]) + (p[1] * q[0])
	return z
}

// Quad method returns the non-negative quadrance of z.
func (z *Dual) Quad() float64 {
	return (new(Dual).Mul(z, new(Dual).Conj(z)))[0]
}

// IsZeroDiv method returns true if z is a zero divisor. This is equivalent to
// z being nilpotent: z = bε.
func (z *Dual) IsZeroDiv() bool {
	return !notEquals(z[0], 0)
}

// Inv method sets z equal to the inverse of x, and returns z. If x is a
// zero divisor, then Inv panics.
func (z *Dual) Inv(x *Dual) *Dual {
	if x.IsZeroDiv() {
		panic("zero divisor")
	}
	return z.Scal(new(Dual).Conj(x), 1/x.Quad())
}

// Quo method sets z equal to the quotient x/y, and returns z. If y is a
// zero divisor, then Quo panics.
func (z *Dual) Quo(x, y *Dual) *Dual {
	if y.IsZeroDiv() {
		panic("denominator is a zero divisor")
	}
	return z.Scal(new(Dual).Mul(x, new(Dual).Conj(y)), 1/y.Quad())
}
