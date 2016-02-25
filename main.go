package main

import . "./ex6"

func main() {

	broadcast_listen_port := 25001
	local_listen_port := 20020

	Udp_init(local_listen_port, broadcast_listen_port)
	Elev_init()

	for {

	}
}
