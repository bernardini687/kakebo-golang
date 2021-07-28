package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/shopspring/decimal"
)

func input(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return string(data)
}

func main() {
	data := input("good.txt")

	consoleLog(data)

	balance, err := Compute(data)
	if err != nil {
		panic(err)
	}

	consoleLog(balance.String())
}

func consoleLog(a interface{}) {
	fmt.Printf("%#v\n", a)
}

// PROTOTYPE IMPLEMENTATION BEGINS:

const (
	NewLine string = "\n"

	Monthly int = 1
	Yearly  int = 12
)

var IntervalDictionary = map[string]int{
	"m": Monthly,
	"M": Monthly,
	"y": Yearly,
	"Y": Yearly,
}

func calcDue(dueFields []string) (decimal.Decimal, error) {
	if len(dueFields) < 2 {
		return decimal.Decimal{}, fmt.Errorf("at least 2 fields required")
	}

	amount, err := decimal.NewFromString(dueFields[0])
	if err != nil {
		return decimal.Decimal{}, err
	}

	interval, ok := IntervalDictionary[dueFields[1]]
	if !ok {
		return decimal.Decimal{}, fmt.Errorf("unknown interval '%s'", dueFields[1])
	}

	return amount.Div(
		decimal.NewFromInt(int64(interval)),
	), nil
}

func Compute(duesData string) (decimal.Decimal, error) {
	trimmed := strings.Trim(duesData, NewLine)

	var balance decimal.Decimal

	for _, dueLine := range strings.Split(trimmed, NewLine) {
		due, err := calcDue(strings.Fields(dueLine))
		if err != nil {
			return decimal.Decimal{}, err
		}

		consoleLog(due.String())

		balance = balance.Add(due)
	}

	return balance, nil
}
