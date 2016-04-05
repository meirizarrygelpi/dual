package dual

import "fmt"

const epsilon = 0.00000001

// Dual type represents a dual number over the real numbers.
type Dual [2]float64

// notEquals function returns true if a and b are not equal.
func notEquals(a, b float64) bool {
	return ((a - b) > epsilon) || ((b - a) > epsilon)
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

// Set method sets z equal to x.
func (z *Dual) Set(x *Dual) *Dual {
	for i, v := range x {
		z[i] = v
	}

	return z
}

// String method returns the string version of a Dual value.
func (z Dual) String() string {
	if z[1] == 0 {
		return fmt.Sprintf("%g", z[0])
	}

	if z[0] == 0 {
		return fmt.Sprintf("%gε", z[1])
	}

	if z[1] < 0 {
		return fmt.Sprintf("%g - %gε", z[0], -z[1])
	}

	return fmt.Sprintf("%g + %gε", z[0], z[1])
}

// New function returns a pointer to a Dual value made from two given real
// numbers (i.e. float64s): a + bε.
func New(a, b float64) *Dual {
	z := new(Dual)
	z[0] = a
	z[1] = b

	return z
}

// Scalar method sets z equal to a*x, and returns z.
func (z *Dual) Scalar(x *Dual, a float64) *Dual {
	for i, v := range x {
		z[i] = a * v
	}

	return z
}

// Neg method sets z equal to the negative of x, and returns z.
func (z *Dual) Neg(x *Dual) *Dual {
	return z.Scalar(x, -1)
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
	p := new(Dual).Set(x)
	q := new(Dual).Set(y)
	z[0] = p[0] * q[0]
	z[1] = (p[0] * q[1]) + (p[1] * q[0])

	return z
}

// Quad method returns the non-negative quadrance of z.
func (z *Dual) Quad() float64 {
	return (new(Dual).Mul(z, new(Dual).Conj(z)))[0]
}
