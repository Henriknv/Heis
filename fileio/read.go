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

func Read(filename string) ([]int, []int) {

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		Println("File does not exist at location. Making new queue file.")
		os.Create(filename)
	}

	dat, err := ioutil.ReadFile(filename)
	check(err)

	var b [2][]int
	json.Unmarshal(dat, &b)

	internal_orders := b[0]
	external_orders := b[1]

	return internal_orders, external_orders
}
