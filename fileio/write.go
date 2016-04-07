package fileio

import (
	"encoding/json"
	. "fmt"
	"io/ioutil"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Write(filename string, internal_orders []int, external_orders []int) {

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		os.Create(filename)
	}

	var orders [2][]int

	orders[0] = internal_orders
	orders[1] = external_orders

	b, _ := json.Marshal(orders)

	err := ioutil.WriteFile(filename, b, 0644)
	check(err)
}
