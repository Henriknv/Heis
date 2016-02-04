package main

import . "./network"
import . "./elev"

func main() {

	broadcast_listen_port := 25000
	local_listen_port := 20020

	Udp_init(local_listen_port, broadcast_listen_port)
	Elev_init()

}
