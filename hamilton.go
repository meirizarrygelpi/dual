package dual

import (
	"fmt"
	"math"
	"strings"

	"github.com/meirizarrygelpi/qtr"
)

// A Hamilton represents a dual Hamilton quaternion as an ordered array of two
// pointers to qtr.Hamilton values.
type Hamilton [2]*qtr.Hamilton

var (
	// Symbols for the canonical dual Hamilton quaternion basis.
	symbHamilton = [8]string{"", "i", "j", "k", "ε", "εi", "εj", "εk"}
)

// String returns the string version of a Hamilton value. If z corresponds to
// the dual Hamilton quaternion a + bi + cj + dk + eε + fεi + gεj + hεk, then
// the string is "(a+bi+cj+dk+eε+fεi+gεj+hεk)", similar to complex128 values.
func (z *Hamilton) String() string {
	v := make([]float64, 8)
	v[0], v[1], v[2], v[3] = (z[0])[0], (z[0])[1], (z[0])[2], (z[0])[3]
	v[4], v[5], v[6], v[7] = (z[1])[0], (z[1])[1], (z[1])[2], (z[1])[3]
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
		a[j+1] = symbQ[i]
		i++
	}
	a[16] = ")"
	return strings.Join(a, "")
}

// Equals returns true if z and y are equal.
func (z *Hamilton) Equals(y *Hamilton) bool {
	for i := range z {
		if !z[i].Equals(y[i]) {
			return false
		}
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Hamilton) Copy(y *Hamilton) *Hamilton {
	for i, v := range y {
		z[i] = new(qtr.Hamilton).Copy(v)
	}
	return z
}

// NewHamilton returns a pointer to a Hamilton value made from eight given
// float64 values.
func NewHamilton(a, b, c, d, e, f, g, h float64) *Hamilton {
	z := new(Hamilton)
	z[0] = qtr.NewHamilton(a, b, c, d)
	z[1] = qtr.NewHamilton(e, f, g, h)
	return z
}

// IsHamiltonInf returns true if any of the components of z are infinite.
func (z *Hamilton) IsHamiltonInf() bool {
	for _, v := range z {
		if v.IsHamiltonInf() {
			return true
		}
	}
	return false
}

// HamiltonInf returns a pointer to a dual Hamilton quaternion infinity value.
func HamiltonInf(a, b, c, d, e, f, g, h int) *Hamilton {
	z := new(Hamilton)
	z[0] = qtr.HamiltonInf(a, b, c, d)
	z[1] = qtr.HamiltonInf(e, f, g, h)
	return z
}

// IsHamiltonNaN returns true if any component of z is NaN and neither is an
// infinity.
func (z *Hamilton) IsHamiltonNaN() bool {
	for _, v := range z {
		if v.IsHamiltonInf() {
			return false
		}
	}
	for _, v := range z {
		if v.IsHamiltonNaN() {
			return true
		}
	}
	return false
}

// HamiltonNaN returns a pointer to a dual Hamilton quaternion NaN value.
func HamiltonNaN() *Hamilton {
	z := new(Hamilton)
	z[0] = qtr.HamiltonNaN()
	z[1] = qtr.HamiltonNaN()
	return z
}

// ScalR sets z equal to y scaled by a (with a being a qtr.Hamilton pointer) on
// the right, and returns z.
//
// This is a special case of Mul:
// 		ScalR(y, a) = Mul(y, Hamilton{a, 0})
func (z *Hamilton) ScalR(y *Hamilton, a *qtr.Hamilton) *Hamilton {
	for i, v := range y {
		z[i] = new(qtr.Hamilton).Mul(v, a)
	}
	return z
}

// ScalL sets z equal to y scaled by a (with a being a qtr.Hamilton pointer) on
// the left, and returns z.
//
// This is a special case of Mul:
// 		ScalL(y, a) = Mul(Hamilton{a, 0}, y)
func (z *Hamilton) ScalL(a *qtr.Hamilton, y *Hamilton) *Hamilton {
	for i, v := range y {
		z[i] = new(qtr.Hamilton).Mul(a, v)
	}
	return z
}

// Dil sets z equal to the dilation of y by a, and returns z.
//
// This is a special case of Mul:
// 		Dil(y, a) = Mul(y, Hamilton{qtr.Hamilton{a, 0, 0, 0}, 0})
func (z *Hamilton) Dil(y *Hamilton, a float64) *Hamilton {
	for i, v := range y {
		z[i] = new(qtr.Hamilton).Scal(v, a)
	}
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Hamilton) Neg(y *Hamilton) *Hamilton {
	return z.Dil(y, -1)
}

// Conj sets z equal to the quaternion conjugate of y, and returns z.
func (z *Hamilton) Conj(y *Hamilton) *Hamilton {
	for i, v := range y {
		z[i] = new(qtr.Hamilton).Conj(v)
	}
	return z
}

// DualConj sets z equal to the dual conjugate of y, and returns z.
func (z *Hamilton) DualConj(y *Hamilton) *Hamilton {
	z[0] = new(qtr.Hamilton).Copy(y[0])
	z[1] = new(qtr.Hamilton).Neg(y[1])
	return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Hamilton) Add(x, y *Hamilton) *Hamilton {
	for i, v := range x {
		z[i] = new(qtr.Hamilton).Add(v, y[i])
	}
	return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Hamilton) Sub(x, y *Hamilton) *Hamilton {
	for i, v := range x {
		z[i] = new(qtr.Hamilton).Sub(v, y[i])
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
func (z *Hamilton) Mul(x, y *Hamilton) *Hamilton {
	p := new(Hamilton).Copy(x)
	q := new(Hamilton).Copy(y)
	z[0] = new(qtr.Hamilton).Mul(p[0], q[0])
	z[1] = new(qtr.Hamilton).Add(
		new(qtr.Hamilton).Mul(p[0], q[1]),
		new(qtr.Hamilton).Mul(p[1], q[0]),
	)
	return z
}

// Commutator sets z equal to the commutator of x and y, and returns z.
func (z *Hamilton) Commutator(x, y *Hamilton) *Hamilton {
	return z.Sub(new(Hamilton).Mul(x, y), new(Hamilton).Mul(y, x))
}

// Quad returns the quadrance of z, a pointer to a Real value.
func (z *Hamilton) Quad() *Real {
	p := new(Hamilton).Mul(z, new(Hamilton).Conj(z))
	r := new(Real)
	r[0] = (p[0])[0]
	r[1] = (p[1])[0]
	return r
}

// IsZeroDiv returns true if z is a zero divisor. This is equivalent to
// z being nilpotent (i.e. z² = 0).
func (z *Hamilton) IsZeroDiv() bool {
	return !z[0].Equals(&qtr.Hamilton{0, 0, 0, 0})
}
