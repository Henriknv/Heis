package main

import . "./network"
import . "./elev"
import . "fmt"
import (
	"time"
)

func main() {

	broadcast_listen_port := 25001
	local_listen_port := 20020

	Udp_init(local_listen_port, broadcast_listen_port)
	Elevator_init()

	is_master := Master_or_slave()
	Println(is_master)

	go Elev_maintenance()
	go Get_orders()
	go Execute_orders()
	//go Run_elevator()
	go Elev_lights()

	for {
		//Println("Elev_orders:", Elev_orders)
		//Println(Local_order_matrix)
		time.Sleep(1000 * time.Millisecond)
	}
}
