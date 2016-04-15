package main

import . "./network"
import . "./elev"
import . "fmt"
import . "time"

//import ."./driver"
//import ."./fileio"

func main() {

	broadcast_listen_port := 25001
	local_listen_port := 20020

	Master_output_ch := make(chan MOSI, 100)
	Master_input_ch := make(chan MISO, 100)
	Slave_input_ch := make(chan MOSI, 100)
	Slave_output_ch := make(chan MISO, 100)

	Udp_init(local_listen_port, broadcast_listen_port, Slave_output_ch, Slave_input_ch, Master_input_ch, Master_output_ch)
	Elevator_init()

	go Slave(Slave_input_ch)

	is_master := Master_or_slave(Slave_input_ch)

	Println(is_master)

	if is_master {

		go Master(Master_input_ch, Master_output_ch)
		go Spam(Master_output_ch)

	}

	go Elev_maintenance()
	go Get_orders(Slave_output_ch)
	go Execute_orders()
	go Elev_lights()

	printElevOrders := Tick(100 * Millisecond)

	for {
		select {
		case <-printElevOrders:
			//Println("Elev orders: ", Elev_orders, "\n", "Elev costs: ", Elev_costs)
		default:

		}

	}
}
