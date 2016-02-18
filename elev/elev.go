package elev

import . ".././network"
import (
	"fmt"
	//"net"
	. "time"
)

func master_or_slave() bool {

	var buf Msg_struct
	_ = buf

	var is_master bool

	select {

	case <-After(1 * Second):
		is_master = true

	case buf = <-Send_chan:
		is_master = false

	}

	if is_master {

		fmt.Println("I am master.")

	} else {

		fmt.Println("I am slave.")

	}

	return is_master
}

func Elev_init() {

	is_master := master_or_slave()

}
