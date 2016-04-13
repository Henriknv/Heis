package fileio

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

import .".././constants"

func Write(Local_order_matrix [N_FLOORS][N_BUTTONS]int) {

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		os.Create(filename)
	}

	buf, _ := json.Marshal(Local_order_matrix)

	err := ioutil.WriteFile(filename, buf, 0644)
	check(err)
}
