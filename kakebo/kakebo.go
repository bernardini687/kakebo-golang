package kakebo

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// CalcBalance
//
// Example:
//
//     -120 y foo
//     -34.5 m bar
//     -6 M baz
//     789 Y xyzzy
//
func CalcBalance(duesData string) (decimal.Decimal, error) {
	dues := strings.Trim(duesData, "\n")

	var balance decimal.Decimal

	for _, due := range strings.Split(dues, "\n") {
		net, err := calcDue(strings.Fields(due))
		if err != nil {
			return decimal.Decimal{}, err
		}

		balance = balance.Add(net)
	}

	return balance, nil
}

// CalcMonth
//
// Example:
//
//     1.2 foo
//     3.45 bar
//     6 baz
//     78.09 xyzzy
//
func CalcMonth(monthData string) (decimal.Decimal, error) {
	entries := strings.Trim(monthData, "\n") // TODO: much similar to CalcBalance, introduce tests then refactor!

	var tot decimal.Decimal

	for _, entry := range strings.Split(entries, "\n") {
		val, err := calcEntry(strings.Fields(entry))
		if err != nil {
			return decimal.Decimal{}, err
		}

		tot = tot.Add(val)
	}

	return tot, nil
}

// DisplayMonth
//
// Example:
//
//     January 2020
//
//     Foo	1,2
//     Bar	3,45
//     Baz	6,00
//
//     Tot	10,65
//
func DisplayMonth(monthData string, period time.Time) (string, error) {
	// TODO: potential sub-optimal solution:
	// when building the month display, we would iterate once for the total
	// and a second time for the formatting of the display rows.
	//
	// modularity vs. efficiency issue?
	//
	tot, err := CalcMonth(monthData)
	if err != nil {
		return "", err
	}

	var display []string

	display = append(display, fmt.Sprintln(period.Month(), period.Year()))

	// TODO: build each line from `monthData`
	display = append(display, fmt.Sprintf("%s\t%s", "Foo", "1,2"))
	display = append(display, fmt.Sprintf("%s\t%s", "Bar", "2,45"))
	display = append(display, fmt.Sprintf("%s\t%s", "Xyzzy", "6,00"))

	display = append(display, fmt.Sprintf("\nTot\t%s\n", tot))

	return strings.Join(display, "\n"), nil
}

const (
	monthly int64 = 1
	yearly  int64 = 12
)

var intervalDictionary = map[string]int64{
	"m": monthly,
	"M": monthly,
	"y": yearly,
	"Y": yearly,
}

func calcDue(fields []string) (decimal.Decimal, error) {
	if len(fields) < 2 {
		return decimal.Decimal{}, fmt.Errorf("at least 2 fields required")
	}

	amount, err := decimal.NewFromString(fields[0])
	if err != nil {
		return decimal.Decimal{}, err
	}

	interval, ok := intervalDictionary[fields[1]]
	if !ok {
		return decimal.Decimal{}, fmt.Errorf("unknown interval '%s'", fields[1])
	}

	return amount.Div(decimal.NewFromInt(interval)), nil
}

func calcEntry(fields []string) (decimal.Decimal, error) {
	if len(fields) < 1 {
		return decimal.Decimal{}, fmt.Errorf("at least 1 field required")
	}

	amount, err := decimal.NewFromString(fields[0])
	if err != nil {
		return decimal.Decimal{}, err
	}

	return amount, nil
}
