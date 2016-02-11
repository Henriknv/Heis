package main

import . "encoding/json"
import "fmt"

type test struct {
	Text string
}

func main() {

	one_test := &test{Text: "HOR2Gq"}

	fmt.Println(one_test.Text)

	random, _ := Marshal(one_test)

	fmt.Println(random)

	var two_test test

	Unmarshal(random, &two_test)

	fmt.Println(two_test.Text)

}
