package kakebo

import (
	"testing"
	"time"
)

// Test CalcBalance
//
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
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: %#v, `<nil>`", balance, err, want)
	}
}

func TestCalcBalanceFieldsErr(t *testing.T) {
	data := `-120 y foo
-34.5
-6 M baz
`
	want := "at least 2 fields required"

	got, err := CalcBalance(data)
	message := err.Error()
	balance := got.StringFixed(2)

	if balance != "0.00" || message != want {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: \"0.00\", `%v`", balance, message, want)
	}
}

func TestCalcBalanceEmptyErr(t *testing.T) {
	data := ""
	want := "at least 2 fields required"

	got, err := CalcBalance(data)
	message := err.Error()
	balance := got.StringFixed(2)

	if balance != "0.00" || message != want {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: \"0.00\", `%v`", balance, message, want)
	}
}

func TestCalcBalanceAmountErr(t *testing.T) {
	data := `-120 y foo
-34,5 m bar
-6 M baz
`
	want := "can't convert -34,5 to decimal"

	got, err := CalcBalance(data)
	message := err.Error()
	balance := got.StringFixed(2)

	if balance != "0.00" || message != want {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: \"0.00\", `%v`", balance, message, want)
	}
}

func TestCalcBalanceIntervalErr(t *testing.T) {
	data := `-120 y foo
-34.5 X bar
-6 M baz
`
	want := "unknown interval 'X'"

	got, err := CalcBalance(data)
	message := err.Error()
	balance := got.StringFixed(2)

	if balance != "0.00" || message != want {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: \"0.00\", `%v`", balance, message, want)
	}
}

// Test CalcBalance
//
func TestCalcMonth(t *testing.T) {
	data := `1.2 foo
3.45 bar
6 baz
78.09 xyzzy
`
	want := "88.74"

	got, err := CalcMonth(data)
	balance := got.StringFixed(2)

	if balance != want || err != nil {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: %#v, `<nil>`", balance, err, want)
	}
}

func TestCalcMonthFieldsErr(t *testing.T) {
	data := `1.2 foo

6 baz
78.09 xyzzy
`
	want := "at least 1 field required"

	got, err := CalcMonth(data)
	message := err.Error()
	balance := got.StringFixed(2)

	if balance != "0.00" || message != want {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: \"0.00\", `%v`", balance, message, want)
	}
}

func TestCalcMonthEmptyErr(t *testing.T) {
	data := ""
	want := "at least 1 field required"

	got, err := CalcMonth(data)
	message := err.Error()
	balance := got.StringFixed(2)

	if balance != "0.00" || message != want {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: \"0.00\", `%v`", balance, message, want)
	}
}

func TestCalcMonthAmountErr(t *testing.T) {
	data := `1.2 foo
bar
6 baz
78.09 xyzzy
`
	want := "can't convert bar to decimal"

	got, err := CalcMonth(data)
	message := err.Error()
	balance := got.StringFixed(2)

	if balance != "0.00" || message != want {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: \"0.00\", `%v`", balance, message, want)
	}
}

// Test DisplayMonth
//
func TestDisplayMonth(t *testing.T) {
	data := `1.2 foo
3.45 bar
6 baz
78.09 xyzzy
`
	want := `November 2009

Foo	1.20
Bar	3.45
Baz	6.00
Xyzzy	78.09

Tot	88.74
`

	date := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	got, err := DisplayMonth(data, date)

	if got != want || err != nil {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: %#v, `<nil>`", got, err, want)
	}
}
