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

	buf := make([]byte, 64)

	_, _, err := conn.ReadFromUDP(buf)

	if err != nil {

		client = 1
		fmt.Println("I am master.")
		conn.Close();
		Udp_send()

	} else {

		client = -1
		fmt.Println("I am slave.")
		conn.CLose();
		Udp_receive()

	}

	//defer conn.Close()

	return client
}

func Elev_init() {

	master_or_slave(Broadcast_addr)

}