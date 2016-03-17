package main

import . "./network"
import . "./elev"

func main() {

	broadcast_listen_port := 25001
	local_listen_port := 20020

	Udp_init(local_listen_port, broadcast_listen_port)
	Elevator_init()

	//is_master := Master_or_slave()

	go Elev_maintenance()
	go Get_orders()
	go Run_elevator()

	for {

	}
}
