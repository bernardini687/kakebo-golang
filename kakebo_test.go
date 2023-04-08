package kakebo

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// Test FormatEntries
func TestFormatEntries(t *testing.T) {
	data := `1.2 foo
3.45 bar
6 baz
78.09 xyzzy
`
	want := `Foo	1.20
Bar	3.45
Baz	6.00
Xyzzy	78.09
`

	got, err := FormatEntries(data)

	if got != want || err != nil {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: %#v, `<nil>`", got, err, want)
	}
}

func TestFormatEntriesComma(t *testing.T) {
	data := `1,2 foo
3,45 bar
6 baz
78,09 xyzzy
`
	want := `Foo	1.20
Bar	3.45
Baz	6.00
Xyzzy	78.09
`

	got, err := FormatEntries(data)

	if got != want || err != nil {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: %#v, `<nil>`", got, err, want)
	}
}

func TestFormatEntriesFieldsErr(t *testing.T) {
	data := `1.2 foo
3.45
6 baz
78.09 xyzzy
`
	want := "at least 2 fields required"

	got, err := FormatEntries(data)
	message := err.Error()

	if message != want || got != "" {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: \"0.00\", `%v`", got, message, want)
	}
}

func TestFormatEntriesEmptyErr(t *testing.T) {
	data := ""
	want := "at least 2 fields required"

	got, err := FormatEntries(data)
	message := err.Error()

	if message != want || got != "" {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: \"0.00\", `%v`", got, message, want)
	}
}

func TestFormatEntriesAmountErr(t *testing.T) {
	data := `1.2 foo
bar bar
6 baz
78.09 xyzzy
`
	want := "can't convert bar to decimal"

	got, err := FormatEntries(data)
	message := err.Error()

	if message != want || got != "" {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: \"0.00\", `%v`", got, message, want)
	}
}

// Test CalcBalance
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

	if message != want || balance != "0.00" {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: \"0.00\", `%v`", balance, message, want)
	}
}

func TestCalcBalanceEmptyErr(t *testing.T) {
	data := ""
	want := "at least 2 fields required"

	got, err := CalcBalance(data)
	message := err.Error()
	balance := got.StringFixed(2)

	if message != want || balance != "0.00" {
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

	if message != want || balance != "0.00" {
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

	if message != want || balance != "0.00" {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: \"0.00\", `%v`", balance, message, want)
	}
}

// Test CalcMonth
func TestCalcMonth(t *testing.T) {
	monthData := `Foo	1.20
Bar	3.45
Baz	6.00
Xyzzy	78.09
`
	want := "88.74"

	got, err := CalcMonth(monthData)
	balance := got.StringFixed(2)

	if balance != want || err != nil {
		t.Fatalf("\n GOT: %#v, `%v`\nWANT: %#v, `<nil>`", balance, err, want)
	}
}

// Test DisplayMonth
func TestDisplayMonth(t *testing.T) {
	monthData := `Foo	1.20
Bar	3.45
Baz	6.00
Xyzzy	78.09
`
	want := `November 2009

Foo	1,20
Bar	3,45
Baz	6,00
Xyzzy	78,09

Tot	100,00
`

	date := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	tot := decimal.NewFromInt(100)
	got := DisplayMonth(date, monthData, tot)

	if got != want {
		t.Fatalf("\n GOT: %#v\nWANT: %#v", got, want)
	}
}

// Test DisplayStats
func TestDisplayStats(t *testing.T) {
	want := `10 November 2009

Save goal	100,00
Monthly budget	900,00
Daily budget	30,00

End of month	33%
Amount spent	11%
`

	date := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	bal := decimal.NewFromInt(1000)
	tot := decimal.NewFromInt(100)
	got := DisplayStats(date, bal, tot, 10)

	if got != want {
		t.Fatalf("\n GOT: %#v\nWANT: %#v", got, want)
	}
}

// Test DisplayDues
func TestDisplayDues(t *testing.T) {
	dues := `-120 y foo
-34.5 m bar
-6 M baz
789 Y xyzzy
1200 M incoming
`

	want := `Incoming	1200,00
Xyzzy	65,75
Baz	-6,00
Foo	-10,00
Bar	-34,50
`

	got := DisplayDues(dues)

	if got != want {
		t.Fatalf("\n GOT: %#v\nWANT: %#v", got, want)
	}
}

func TestDisplayDuesInvalidDuesErr(t *testing.T) {
	dues := `-120 y foo
-34.5
-6 M baz
`

	want := "invalid dues"

	got := DisplayDues(dues)

	if got != want {
		t.Fatalf("\n GOT: %#v\nWANT: %#v", got, want)
	}
}
