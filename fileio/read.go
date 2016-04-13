package fileio

import (
	"encoding/json"
	. "fmt"
	"io/ioutil"
	"os"
)

import .".././constants"

const filename = "order_backup.txt"

var read_buf [N_FLOORS][N_BUTTONS]int

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Read() [N_FLOORS][N_BUTTONS]int{

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		Println("File does not exist at location. Making new queue file.")
		os.Create(filename)
	}

	dat, err := ioutil.ReadFile(filename)
	check(err)

	json.Unmarshal(dat, &read_buf)

	return read_buf
}
