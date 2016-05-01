package dual

import (
    "fmt"
    "math"
    "strings"
)

// A Complex represents a dual complex number as an ordered array of four
// float64 values.
type Complex [4]float64

var (
    symbC = [4]string{"", "i", "ε", "εi"}
)

// String returns the string representation of a Complex value. If z
// corresponds to the dual complex number a + bi + cε + diε, then the string is
// "(a+bi+cε+diε)", similar to complex128 values.
func (z *Complex) String() string {
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
        a[j+1] = symbC[i]
        i++
    }
    a[8] = ")"
    return strings.Join(a, "")
}

// Equals returns true if z and y are equal.
func (z *Complex) Equals(y *Complex) bool {
    for i := range z {
        if notEquals(z[i], y[i]) {
            return false
        }
    }
    return true
}

// Copy copies y onto z, and returns z.
func (z *Complex) Copy(y *Complex) *Complex {
    for i, v := range y {
        z[i] = v
    }
    return z
}

// NewComplex returns a pointer to a Complex value made from four given float64
// values.
func NewComplex(a, b, c, d float64) *Complex {
    z := new(Complex)
    z[0] = a
    z[1] = b
    z[2] = c
    z[3] = d
    return z
}

// IsComplexInf returns true if any of the components of z are infinite.
func (z *Complex) IsComplexInf() bool {
    for _, v := range z {
        if math.IsInf(v, 0) {
            return true
        }
    }
    return false
}

// ComplexInf returns a pointer to a dual complex infinity value.
func ComplexInf(a, b, c, d int) *Complex {
    z := new(Complex)
    z[0] = math.Inf(a)
    z[1] = math.Inf(b)
    z[2] = math.Inf(c)
    z[3] = math.Inf(d)
    return z
}

// IsComplexNaN returns true if any component of z is NaN and neither is an
// infinity.
func (z *Complex) IsComplexNaN() bool {
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

// ComplexNaN returns a pointer to a dual complex NaN value.
func ComplexNaN() *Complex {
    nan := math.NaN()
    return &Complex{nan, nan, nan, nan}
}

// Scal sets z equal to y scaled by a, and returns z.
func (z *Complex) Scal(y *Complex, a float64) *Complex {
    for i, v := range y {
        z[i] = a * v
    }
    return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Complex) Neg(y *Complex) *Complex {
    return z.Scal(y, -1)
}

// DConj sets z equal to the dual conjugate of y, and returns z.
func (z *Complex) DConj(y *Complex) *Complex {
    z[0] = +y[0]
    z[1] = +y[1]
    z[2] = -y[2]
    z[3] = -y[3]
    return z
}

// Conj sets z equal to the complex conjugate of y, and returns z.
func (z *Complex) Conj(y *Complex) *Complex {
    z[0] = +y[0]
    z[1] = -y[1]
    z[2] = +y[2]
    z[3] = -y[3]
    return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Complex) Add(x, y *Complex) *Complex {
    for i, v := range x {
        z[i] = v + y[i]
    }
    return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Complex) Sub(x, y *Complex) *Complex {
    for i, v := range x {
        z[i] = v - y[i]
    }
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
func (z *Complex) Mul(x, y *Complex) *Complex {
    p := new(Complex).Copy(x)
    q := new(Complex).Copy(y)
    z[0] = (p[0] * q[0]) - (p[1] * q[1])
    z[1] = (p[0] * q[1]) + (p[1] * q[0])
    z[2] = (p[0] * q[2]) - (p[1] * q[3]) + (p[2] * q[0]) - (p[3] * q[1])
    z[3] = (p[0] * q[3]) + (p[1] * q[2]) + (p[2] * q[1]) + (p[3] * q[0])
    return z
}

// Quad returns the quadrance of z, a dual real value.
func (z *Complex) Quad() *Real {
    p := new(Complex).Mul(z, new(Complex).Conj(z))
    d := new(Real)
    d[0] = p[0]
    d[1] = p[2]
    return d
}

// DQuad returns the dual quadrance of z, a complex128 value.
func (z *Complex) DQuad() complex128 {
    p := new(Complex).Mul(z, new(Complex).DConj(z))
    return complex(p[0], p[1])
}
