package kakebo

import (
	"testing"
)

func TestCalcBalance(t *testing.T) {
	data := `-120 y foo
-34.5 m bar
-6 M baz
789 Y xyzzy
`
	want := "15.25"

	got, err := CalcBalance(data)
	balance := got.StringFixed(2)

	if balance != want || err != nil {
		t.Fatalf("\nGOT: %#v, `%v`\nWANT: %#v, `<nil>`", balance, err, want)
	}
}

func TestCalcBalanceEmpty(t *testing.T) {
	data := ""
	want := "at least 2 fields required"

	got, err := CalcBalance(data)
	message := err.Error()
	balance := got.StringFixed(2)

	if balance != "0.00" || message != want {
		t.Fatalf("\nGOT: %#v, `%v`\nWANT: \"0.00\", `%v`", balance, message, want)
	}
}

func TestCalcBalanceInterval(t *testing.T) {
	data := `-120 y foo
-34.5 X bar
-6 M baz
`
	want := "unknown interval 'X'"

	got, err := CalcBalance(data)
	message := err.Error()
	balance := got.StringFixed(2)

	if balance != "0.00" || message != want {
		t.Fatalf("\nGOT: %#v, `%v`\nWANT: \"0.00\", `%v`", balance, message, want)
	}
}
