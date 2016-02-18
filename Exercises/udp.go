package udp_send

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func udp_receive() {

	server := ":20020"

	fmt.Println("Connecting to server at", server)

	server_addr, _ := net.ResolveUDPAddr("udp", server)

	fmt.Println("Connected to server:", server)

	conn, _ := net.ListenUDP("udp", server_addr)

	defer conn.Close()

	buf := make([]byte, 1024)

	for {

		time.Sleep(100 * time.Millisecond)
		n, _, _ := conn.ReadFromUDP(buf)
		fmt.Println("Received:", string(buf[0:n]))

	}
}

func udp_send() {

	server_addr, _ := net.ResolveUDPAddr("udp", "129.241.187.23:20020")
	local_addr, _ := net.ResolveUDPAddr("udp", "129.241.187.155")

	conn, _ := net.DialUDP("udp", local_addr, server_addr)

	defer conn.Close()

	i := 0

	for {

		time.Sleep(100 * time.Millisecond)
		msg := strconv.Itoa(i)
		i++
		buf := []byte(msg)
		conn.Write(buf)

	}

}

func main() {

	go udp_send()
	go udp_receive()

	for {
		time.Sleep(time.Second)
	}

}
