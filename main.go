package main

import . "./network"
import . "./elev"
//import . "fmt"
import ."time"
//import ."./fileio"

func main() {

	broadcast_listen_port := 25001
	local_listen_port := 20020

	send_chan := make(chan Msg_struct, 100)
	receive_chan := make(chan Msg_struct, 100)

	Udp_init(local_listen_port, broadcast_listen_port, send_chan, receive_chan)
	Elevator_init(send_chan, receive_chan)

	is_master := Master_or_slave()

	if is_master {

		go Master()
		go Slave()

	}else{
		go Slave()
	}

	go Elev_maintenance()
	go Get_orders()
	go Execute_orders()
	go Elev_lights()

	for {

		//Println(Read())
		Sleep(1000 * Millisecond)

	}
}
