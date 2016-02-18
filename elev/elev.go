package elev

import . ".././network"
import (
	"fmt"
	"net"
	"time"
)

func master_or_slave(Broadcast_addr *net.UDPAddr) int {

	conn, _ := net.ListenUDP("udp", Broadcast_addr)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	client := 0

	// i = 0 -> Client = Master
	// i != 0 -> Client = Slave

	buf := make([]byte, 64)

	i, _, _ := conn.ReadFromUDP(buf)

	defer conn.Close()

	if i == 0 {

		client = 1
		fmt.Println("I am master.")

	} else {

		client = -1
		fmt.Println("I am slave.")

	}

	return client
}

func master() {

}

func slave() {
	//Udp_receive()
}

func Elev_init() {

	client := master_or_slave(Broadcast_addr)

	switch client {
	case -1:
		slave()
	case 1:
		master()
	}

}
