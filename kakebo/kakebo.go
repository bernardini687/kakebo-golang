package kakebo

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// FormatEntries
//
// Input example:
//
//     1,2 foo
//     3,45 bar
//     6 baz
//     78.09 xyzzy
//
// Output example:
//
//     Foo	1.20
//     Bar	3.45
//     Baz	6.00
//     Xyzzy	78.09
//
func FormatEntries(entryData string) (string, error) {
	return "", nil
}

// CalcBalance
//
// Input example:
//
//     -120 y foo
//     -34.5 m bar
//     -6 M baz
//     789 Y xyzzy
//
// Output example:
//
//     15.25
//
func CalcBalance(dueData string) (decimal.Decimal, error) {
	return sumValues(dueData, extractDueValue)
}

// CalcMonth
//
// Input example:
//
//     Foo	1.20
//     Bar	3.45
//     Baz	6.00
//     Xyzzy	78.09
//
// Output example:
//
//     88.74
//
func CalcMonth(monthData string) (decimal.Decimal, error) {
	return sumValues(monthData, extractEntryValue)
}

// DisplayMonth
//
// Input example:
//
//     Foo	1.20
//     Bar	3.45
//     Baz	6.00
//
// Output example:
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
	tot, err := CalcMonth(monthData)
	if err != nil {
		return "", err
	}

	var lines []string

	lines = append(lines, fmt.Sprintln(period.Month(), period.Year()))
	lines = append(lines, monthData)
	lines = append(lines, fmt.Sprintf("Tot\t%s\n", tot))

	display := strings.Join(lines, "\n")

	return strings.ReplaceAll(display, ".", ","), nil
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

func extractDueValue(fields []string) (decimal.Decimal, error) {
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

func extractEntryValue(fields []string) (decimal.Decimal, error) {
	return decimal.NewFromString(fields[1])
}

// func extractEntryValue(fields []string) (decimal.Decimal, error) {
// 	if len(fields) < 1 {
// 		return decimal.Decimal{}, fmt.Errorf("at least 1 field required")
// 	}

// 	amount, err := decimal.NewFromString(fields[0])
// 	if err != nil {
// 		return decimal.Decimal{}, err
// 	}

// 	return amount, nil
// }

func sumValues(data string, extractor func([]string) (decimal.Decimal, error)) (decimal.Decimal, error) {
	var tot decimal.Decimal

	for _, line := range lines(data) {
		val, err := extractor(strings.Fields(line))
		if err != nil {
			return decimal.Decimal{}, err
		}

		tot = tot.Add(val)
	}

	return tot, nil
}

func lines(data string) []string {
	return strings.Split(strings.Trim(data, "\n"), "\n")
}
