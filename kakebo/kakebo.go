package kakebo

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// TODO:
// monthData should be coming from MM.txt files.
// i should make shure those files contains only valid amounts.
// so one missing function should be one that takes an input and checks for errors, giving back only valid amounts when no error
// happened.
// thus, any monthData should be treated as a non-empty MM.txt content that only contains valid amounts!
// CalcMonth() and entriesForDisplay() should be affected!

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
	return sumValues(duesData, extractDueValue)
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
	return sumValues(monthData, extractEntryValue)
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
	// here we iterate once over `monthData` to calculate the total.
	// we'll iterate a second time when `entriesForDisplay()` runs.
	// not ideal, but:
	//   1) `displayMonth` will hardly contain a lot of lines
	//   2) `CalcMonth` checks for errors, so if `entriesForDisplay()` runs, we can assume the data is good
	//
	tot, err := CalcMonth(monthData)
	if err != nil {
		return "", err
	}

	var display []string

	display = append(display, fmt.Sprintln(period.Month(), period.Year()))
	display = append(display, entriesForDisplay(monthData)...)
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
	if len(fields) < 1 {
		return decimal.Decimal{}, fmt.Errorf("at least 1 field required")
	}

	amount, err := decimal.NewFromString(fields[0])
	if err != nil {
		return decimal.Decimal{}, err
	}

	return amount, nil
}

func entriesForDisplay(monthData string) []string {
	entries := lines(monthData)

	for i, entry := range entries {
		fields := strings.Fields(entry)
		// this function is meant to run after a call to `CalcMonth()` thus,
		// the input data has already been checked for errors
		amount, _ := decimal.NewFromString(fields[0])
		entries[i] = fmt.Sprintf("%s\t%s", strings.Title(fields[1]), amount.StringFixed(2))
	}

	return entries
}

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
