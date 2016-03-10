package main

import . "./network"
import . "./elev"

func main() {

	broadcast_listen_port := 25001
	local_listen_port := 20020

	Udp_init(local_listen_port, broadcast_listen_port)
	Elevator_init()
	//is_master := Master_or_slave()
	/*
		Execute_order(1)
		Execute_order(3)
		Execute_order(1)
		Execute_order(2)*/

	go Elev_maintenance()

	Execute_order(0)
	Execute_order(3)
	Execute_order(0)

	for {

	}
}
