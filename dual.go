package dual

import "fmt"

const (
	epsilon = 0.00000001
)

// Dual type represents a dual number over the real numbers.
type Dual struct {
	__ [2]float64
}

// E0 method returns the real part of z.
func (z *Dual) E0() float64 { return z.__[0] }

// E1 method returns the dual part of z.
func (z *Dual) E1() float64 { return z.__[1] }

// equals function.
func equals(x, y float64) bool {
	return ((x - y) < epsilon) && ((y - x) < epsilon)
}

// Equals method returns true if z and x are equal.
func (z *Dual) Equals(x *Dual) bool {
	return (equals(z.__[0], x.__[0]) &&
		equals(z.__[1], x.__[1]))
}

// Clone method clones x onto z.
func (z *Dual) Clone(x *Dual) *Dual {
	for i, v := range x.__ {
		z.__[i] = v
	}

	return z
}

// String method returns the string version of a Dual value.
func (z Dual) String() string {
	if z.__[1] == 0 {
		return fmt.Sprintf("%g", z.__[0])
	}

	if z.__[0] == 0 {
		return fmt.Sprintf("%gε", z.__[1])
	}

	if z.__[1] < 0 {
		return fmt.Sprintf("%g - %gε", z.__[0], -z.__[1])
	}

	return fmt.Sprintf("%g + %gε", z.__[0], z.__[1])
}

// New function returns a pointer to a Dual value made from two given real
// numbers (i.e. float64s).
func New(a, b float64) *Dual {
	z := new(Dual)
	z.__[0] = a
	z.__[1] = b

	return z
}

// Scalar method sets z equal to s*x, and returns z.
func (z *Dual) Scalar(x *Dual, s float64) *Dual {
	for i, v := range x.__ {
		z.__[i] = s * v
	}

	return z
}

// Neg method sets z equal to the negative of x, and returns z.
func (z *Dual) Neg(x *Dual) *Dual {
	return z.Scalar(x, -1)
}

// Conj method sets z equal to the conjugate of x, and returns z.
func (z *Dual) Conj(x *Dual) *Dual {
	z.__[0] = x.__[0]
	z.__[1] = -x.__[1]

	return z
}

// Add method sets z to the sum of x and y, and returns z.
func (z *Dual) Add(x, y *Dual) *Dual {
	for i, v := range x.__ {
		z.__[i] = v + y.__[i]
	}

	return z
}

// Sub method sets z to the sum of x and y, and returns z.
func (z *Dual) Sub(x, y *Dual) *Dual {
	for i, v := range x.__ {
		z.__[i] = v - y.__[i]
	}

	return z
}

// Mul method sets z to the product of x and y, and returns z.
func (z *Dual) Mul(x, y *Dual) *Dual {
	p := new(Dual).Clone(x)
	q := new(Dual).Clone(y)
	z.__[0] = p.__[0] * q.__[0]
	z.__[1] = (p.__[0] * q.__[1]) + (p.__[1] * q.__[0])

	return z
}

// Quad method returns the non-negative quadrance of z.
func (z *Dual) Quad() float64 {
	return (new(Dual).Mul(z, new(Dual).Conj(z))).__[0]
}
