package dual

import "testing"

var (
	zero = New(0, 0)
	e0   = New(1, 0)
	e1   = New(0, 1)
)

func TestEquals(t *testing.T) {
	var tests = []struct {
		x    *Dual
		y    *Dual
		want bool
	}{
		{zero, zero, true},
		{e0, e0, true},
		{e1, e1, true},
		{e0, e1, false},
		{e1, e0, false},
		{New(2.03, 3), New(2.0299999999, 3), true},
		{New(1, 2), New(3, 4), false},
	}

	for _, test := range tests {
		if got := test.x.Equals(test.y); got != test.want {
			t.Errorf("Equals(%v, %v) = %v", test.x, test.y, got)
		}
	}
}

func TestString(t *testing.T) {
	var tests = []struct {
		x    *Dual
		want string
	}{
		{zero, "0"},
		{e0, "1"},
		{e1, "1ε"},
		{New(1, 1), "1 + 1ε"},
		{New(1, -1), "1 - 1ε"},
		{New(-1, 1), "-1 + 1ε"},
		{New(-1, -1), "-1 - 1ε"},
	}

	for _, test := range tests {
		if got := test.x.String(); got != test.want {
			t.Errorf("String(%v) = %v, want %v", test.x, got, test.want)
		}
	}
}
