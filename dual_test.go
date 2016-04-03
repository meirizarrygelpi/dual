package dual

import "testing"

var (
	zero = New(0, 0)
	e0   = New(1, 0)
	e1   = New(0, 1)
)

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
