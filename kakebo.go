package kakebo

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// FormatEntries
//
// Input example:
//
//	1,2 foo
//	3,45 bar
//	6 baz
//	78.09 xyzzy
//
// Output example:
//
//	Foo	1.20
//	Bar	3.45
//	Baz	6.00
//	Xyzzy	78.09
func FormatEntries(entryData string) (string, error) {
	var formattedEntries []string

	for _, line := range lines(entryData) {
		entry, err := formatEntry(strings.Fields(line))
		if err != nil {
			return "", err
		}

		formattedEntries = append(formattedEntries, entry)
	}

	return strings.Join(formattedEntries, ""), nil
}

func formatEntry(fields []string) (string, error) {
	if len(fields) < 2 {
		return "", fmt.Errorf("at least 2 fields required")
	}

	s := strings.Replace(fields[0], ",", ".", 1)
	amount, err := decimal.NewFromString(s)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\t%s\n", strings.Title(fields[1]), money(amount)), nil
}

// CalcMonth
//
// Input example:
//
//	Foo	1.20
//	Bar	3.45
//	Baz	6.00
//	Xyzzy	78.09
//
// Output example:
//
//	88.74
func CalcMonth(monthData string) (decimal.Decimal, error) {
	return sumValues(monthData, extractFormattedEntryValue)
}

func extractFormattedEntryValue(fields []string) (decimal.Decimal, error) {
	return decimal.NewFromString(fields[1])
}

// CalcBalance
//
// Input example:
//
//	-120 y foo
//	-34.5 m bar
//	-6 M baz
//	789 Y xyzzy
//
// Output example:
//
//	15.25
func CalcBalance(dueData string) (decimal.Decimal, error) {
	return sumValues(dueData, extractDueValue)
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

// DisplayMonth
//
// Input example:
//
//	Foo	1.20
//	Bar	3.45
//	Baz	6.00
//
// Output example:
//
//	January 2020
//
//	Foo	1,20
//	Bar	3,45
//	Baz	6,00
//
//	Tot	10,65
func DisplayMonth(date time.Time, monthData string, monthTot decimal.Decimal) string {
	var lines []string

	lines = append(lines, fmt.Sprintln(date.Month(), date.Year()))   // header
	lines = append(lines, monthData)                                 // body
	lines = append(lines, fmt.Sprintf("Tot\t%s\n", money(monthTot))) // footer

	text := strings.Join(lines, "\n")

	return replaceDotsWithCommas(text)
}

// DisplayStats
//
// Input example:
//
//	time.Time{2009/11/10}, decimal.Decimal{1000}, decimal.Decimal{100}, 10
//
// Output example:
//
//	10 November 2009
//
//	Save goal	100,00
//	Monthly budget	900,00
//	Daily budget	30,00
//
//	End of month	33%
//	Amount spent	11%
func DisplayStats(date time.Time, balance, monthTot decimal.Decimal, savePercentage int) string {
	var lines []string

	hundred := decimal.NewFromInt(100)
	savePercent := decimal.NewFromInt(int64(savePercentage))
	today := decimal.NewFromInt(int64(date.Day()))

	saveGoal := balance.Div(hundred).Mul(savePercent)
	monthlyBudget := balance.Sub(saveGoal)
	monthDays := daysOfMonth(date)
	dailyBudget := monthlyBudget.DivRound(monthDays, 2)
	endOfMonthPercentage := hundred.Mul(today).DivRound(monthDays, 0)
	spentAmountPercentage := hundred.Mul(monthTot).DivRound(monthlyBudget, 0)

	y, m, d := date.Date()
	body := [][]string{
		{"Save goal", money(saveGoal)},
		{"Monthly budget", money(monthlyBudget)},
		{"Daily budget", money(dailyBudget)},
	}
	footer := [][]string{
		{"End of month", percentage(endOfMonthPercentage)},
		{"Amount spent", percentage(spentAmountPercentage)},
	}

	lines = append(lines, fmt.Sprintln(d, m, y))
	lines = append(lines, formatStats(body))
	lines = append(lines, formatStats(footer))

	text := strings.Join(lines, "\n")

	return replaceDotsWithCommas(text)
}

func daysOfMonth(date time.Time) decimal.Decimal {
	y, m, _ := date.Date()
	endOfMonth := time.Date(y, m+1, 0, 0, 0, 0, 0, date.Location())

	return decimal.NewFromInt(int64(endOfMonth.Day()))
}

func formatStats(stats [][]string) string {
	var text string

	for _, cols := range stats {
		text += fmt.Sprintf("%s\t%s\n", cols[0], cols[1])
	}

	return text
}

type Due struct {
	Amount      decimal.Decimal
	Description string
}

// DisplayDues
//
// Input example:
//
//	-120 y foo
//	-34.5 m bar
//	-6 M baz
//	789 Y xyzzy
//	1200 M incoming
//
// Output example:
//
//	Bar	-34,50
//	Foo	-10,00
//	Baz	-6,00
func DisplayDues(dueData string) string {
	var dues []Due

	for _, line := range lines(dueData) {
		fields := strings.Fields(line)
		val, err := extractDueValue(fields)

		if err != nil {
			return "invalid dues"
		}

		if val.LessThan(decimal.Zero) {
			dues = append(dues, Due{decimal.Decimal.Abs(val), fields[2]})
		}
	}

	sort.SliceStable(dues, func(a, b int) bool {
		return dues[a].Amount.GreaterThan(dues[b].Amount)
	})

	var lines []string

	for _, due := range dues {
		lines = append(lines, fmt.Sprintf("%s\t%s\n", strings.Title(due.Description), money(due.Amount)))
	}

	text := strings.Join(lines, "")

	return replaceDotsWithCommas(text)
}

//
// Common stuff
//

func lines(data string) []string {
	return strings.Split(strings.Trim(data, "\n"), "\n")
}

func money(amount decimal.Decimal) string {
	return amount.StringFixed(2)
}

func percentage(amount decimal.Decimal) string {
	return amount.String() + "%"
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

func replaceDotsWithCommas(text string) string {
	return strings.ReplaceAll(text, ".", ",")
}
