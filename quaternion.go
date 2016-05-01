package dual

import (
	"fmt"
	"math"
	"strings"

	"github.com/meirizarrygelpi/qtr"
)

// A Quaternion represents a dual quaternion number as an ordered array of
// eight float64 values.
type Quaternion [8]float64

var (
	// Symbols for the canonical dual quaternion basis.
	symbQ = [8]string{"", "i", "j", "k", "ε", "εi", "εj", "εk"}
)

// String returns the string version of a Quaternion value. If z corresponds to
// the dual quaternion number a + bi + cj + dk + eε + fεi + gεj + hεk, then the
// string is "(a+bi+cj+dk+eε+fεi+gεj+hεk)", similar to complex128 values.
func (z *Quaternion) String() string {
	a := make([]string, 17)
	a[0] = "("
	a[1] = fmt.Sprintf("%g", z[0])
	i := 1
	for j := 2; j < 16; j = j + 2 {
		switch {
		case math.Signbit(z[i]):
			a[j] = fmt.Sprintf("%g", z[i])
		case math.IsInf(z[i], +1):
			a[j] = "+Inf"
		default:
			a[j] = fmt.Sprintf("+%g", z[i])
		}
		a[j+1] = symbQ[i]
		i++
	}
	a[16] = ")"
	return strings.Join(a, "")
}

// Equals returns true if z and y are equal.
func (z *Quaternion) Equals(y *Quaternion) bool {
	for i := range z {
		if notEquals(z[i], y[i]) {
			return false
		}
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Quaternion) Copy(y *Quaternion) *Quaternion {
	for i, v := range y {
		z[i] = v
	}
	return z
}

// NewQuaternion returns a pointer to a Quaternion value made from eight given
// float64 values.
func NewQuaternion(a, b, c, d, e, f, g, h float64) *Quaternion {
	z := new(Quaternion)
	z[0] = a
	z[1] = b
	z[2] = c
	z[3] = d
	z[4] = e
	z[5] = f
	z[6] = g
	z[7] = h
	return z
}

// IsQuaternionInf returns true if any of the components of z are infinite.
func (z *Quaternion) IsQuaternionInf() bool {
	for _, v := range z {
		if math.IsInf(v, 0) {
			return true
		}
	}
	return false
}

// QuaternionInf returns a pointer to a dual quaternion infinity value.
func QuaternionInf(a, b, c, d, e, f, g, h int) *Quaternion {
	z := new(Quaternion)
	z[0] = math.Inf(a)
	z[1] = math.Inf(b)
	z[2] = math.Inf(c)
	z[3] = math.Inf(d)
	z[4] = math.Inf(e)
	z[5] = math.Inf(f)
	z[6] = math.Inf(g)
	z[7] = math.Inf(h)
	return z
}

// IsQuaternionNaN returns true if any component of z is NaN and neither is an
// infinity.
func (z *Quaternion) IsQuaternionNaN() bool {
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

// QuaternionNaN returns a pointer to a dual quaternion NaN value.
func QuaternionNaN() *Quaternion {
	nan := math.NaN()
	return &Quaternion{nan, nan, nan, nan, nan, nan, nan, nan}
}

// ScalR sets z equal to y scaled by a (with a being float64), and returns z.
//
// This is a special case of Mul:
// 		Mul(z, Quaternion{a, 0, 0, 0, 0, 0, 0, 0})
func (z *Quaternion) ScalR(y *Quaternion, a float64) *Quaternion {
	for i, v := range y {
		z[i] = a * v
	}
	return z
}

// ScalH sets z equal to y scaled by h (with h being qtr.Hamilton), and returns
// z.
//
// This is a special case of Mul:
// 		Mul(z, Quaternion{h[0], h[1], h[2], h[3], 0, 0, 0, 0})
func (z *Quaternion) ScalH(y *Quaternion, h *qtr.Hamilton) *Quaternion {
	z[0] = (y[0] * h[0])
	z[1] = (y[1] * h[0])
	z[2] = (y[2] * h[0])
	z[3] = (y[3] * h[0])
	z[4] = (y[4] * h[0])
	z[5] = (y[5] * h[0])
	z[6] = (y[6] * h[0])
	z[7] = (y[7] * h[0])
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Quaternion) Neg(y *Quaternion) *Quaternion {
	return z.ScalR(y, -1)
}

// DConj sets z equal to the dual conjugate of y, and returns z.
func (z *Quaternion) DConj(y *Quaternion) *Quaternion {
	z[0] = +y[0]
	z[1] = +y[1]
	z[2] = +y[2]
	z[3] = +y[3]
	z[4] = -y[4]
	z[5] = -y[5]
	z[6] = -y[6]
	z[7] = -y[7]
	return z
}

// Conj sets z equal to the complex conjugate of y, and returns z.
func (z *Quaternion) Conj(y *Quaternion) *Quaternion {
	z[0] = +y[0]
	z[1] = -y[1]
	z[2] = -y[2]
	z[3] = -y[3]
	z[4] = +y[4]
	z[5] = -y[5]
	z[6] = -y[6]
	z[7] = -y[7]
	return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Quaternion) Add(x, y *Quaternion) *Quaternion {
	for i, v := range x {
		z[i] = v + y[i]
	}
	return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Quaternion) Sub(x, y *Quaternion) *Quaternion {
	for i, v := range x {
		z[i] = v - y[i]
	}
	return z
}

// Mul sets z equal to the product of x and y, and returns z.
//
// The basic rules are:
// 		i * i = j * j = k * k = -1
// 		i * j = -j * i = k
// 		j * k = -k * j = i
// 		k * i = -i * k = j
// 		ε * ε = 0
// 		ε * i = i * ε = εi
// 		ε * j = j * ε = εj
// 		ε * k = k * ε = εk
// 		εi * i = i * εi = -ε
// 		εj * j = j * εj = -ε
// 		εk * k = k * εk = -ε
// 		εi * j = -j * εi = εk
// 		εj * k = -k * εj = εi
// 		εk * i = -i * εk = εj
// 		ε * εi = εi * ε = 0
// 		ε * εj = εj * ε = 0
// 		ε * εk = εk * ε = 0
// 		εi * εi = εj * εj = εk * εk = 0
// 		εi * εj = εj * εi = 0
// 		εi * εk = εk * εi = 0
// 		εj * εk = εk * εj = 0
// 		εj * εk = εk * εj = 0
// This multiplication rule is noncommutative but associative.
func (z *Quaternion) Mul(x, y *Quaternion) *Quaternion {
	p := new(Quaternion).Copy(x)
	q := new(Quaternion).Copy(y)
	z[0] = (p[0] * q[0]) - (p[1] * q[1]) - (p[2] * q[2]) - (p[3] * q[3])
	z[1] = (p[0] * q[1]) + (p[1] * q[0]) + (p[2] * q[3]) - (p[3] * q[2])
	z[2] = (p[0] * q[2]) - (p[1] * q[3]) + (p[2] * q[0]) + (p[3] * q[1])
	z[3] = (p[0] * q[3]) + (p[1] * q[2]) - (p[2] * q[1]) + (p[3] * q[0])
	z[4] = (p[0] * q[4]) - (p[1] * q[5]) - (p[2] * q[6]) - (p[3] * q[7]) +
		(p[4] * q[0]) - (p[5] * q[1]) - (p[6] * q[2]) - (p[7] * q[3])
	z[5] = (p[0] * q[5]) + (p[1] * q[4]) + (p[2] * q[7]) - (p[3] * q[6]) +
		(p[4] * q[1]) + (p[5] * q[0]) + (p[6] * q[3]) - (p[7] * q[2])
	z[6] = (p[0] * q[6]) - (p[1] * q[7]) + (p[2] * q[4]) + (p[3] * q[5]) +
		(p[4] * q[2]) - (p[5] * q[3]) + (p[6] * q[0]) + (p[7] * q[1])
	z[7] = (p[0] * q[7]) + (p[1] * q[6]) - (p[2] * q[5]) + (p[3] * q[4]) +
		(p[4] * q[3]) + (p[5] * q[2]) - (p[6] * q[1]) + (p[7] * q[0])
	return z
}

// Commutator sets z equal to the commutator of x and y, and returns z.
func (z *Quaternion) Commutator(x, y *Quaternion) *Quaternion {
	return z.Sub(new(Quaternion).Mul(x, y), new(Quaternion).Mul(y, x))
}

// Quad returns the quadrance of z, a dual real value.
func (z *Quaternion) Quad() *Real {
	p := new(Quaternion).Mul(z, new(Quaternion).Conj(z))
	d := new(Real)
	d[0] = p[0]
	d[1] = p[4]
	return d
}

// DQuad returns the dual quadrance of z, a qtr.Hamilton value.
func (z *Quaternion) DQuad() *qtr.Hamilton {
	p := new(Quaternion).Mul(z, new(Quaternion).DConj(z))
	q := new(qtr.Hamilton)
	q[0] = p[0]
	q[1] = p[1]
	q[2] = p[2]
	q[3] = p[3]
	return q
}

// IsZeroDiv returns true if z is a zero divisor. This is equivalent to
// z being nilpotent (i.e. z² = 0).
func (z *Quaternion) IsZeroDiv() bool {
	return !notEquals(z[0], 0) && !notEquals(z[1], 0)
}

// Inv sets z equal to the inverse of y, and returns z. If y is a zero divisor,
// then Inv panics.
func (z *Quaternion) Inv(y *Quaternion) *Quaternion {
	if y.IsZeroDiv() {
		panic("zero divisor")
	}
	return z.ScalH(
		new(Quaternion).DConj(y),
		new(qtr.Hamilton).Inv(y.DQuad()),
	)
}

// Quo sets z equal to the quotient of x and y, and returns z. If y is a zero
// divisor, then Quo panics.
func (z *Quaternion) Quo(x, y *Quaternion) *Quaternion {
	if y.IsZeroDiv() {
		panic("zero divisor denominator")
	}
	return z.ScalH(
		new(Quaternion).Mul(x, new(Quaternion).DConj(y)),
		new(qtr.Hamilton).Inv(y.DQuad()),
	)
}
