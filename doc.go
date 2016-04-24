// Package dual implements dual number arithmetic.
package dual

const delta = 0.00000001

// notEquals function returns true if a and b are not equal.
func notEquals(a, b float64) bool {
	return ((a - b) > delta) || ((b - a) > delta)
}
