package main

import (
	"fmt"
	"io/ioutil"
	"tmp/kakebo-golang/kakebo"
)

func input(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return string(data)
}

func main() {
	data := input("dues.txt")

	consoleLog(data)

	balance, err := kakebo.CalcBalance(data)
	if err != nil {
		panic(err)
	}

	consoleLog(balance.StringFixed(2))

	data = input("entries.txt")

	consoleLog(data)

	tot, err := kakebo.CalcMonth(data)
	if err != nil {
		panic(err)
	}

	consoleLog(tot.StringFixed(2))
}

func consoleLog(a interface{}) {
	fmt.Printf("%#v\n", a)
}
