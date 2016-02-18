package main

import "net"
import "fmt"

//import "bufio"
import "time"

import "strconv"

func tcp_receive(conn *net.TCPConn) {

	defer conn.Close()

	for {

		data := make([]byte, 1024)
		conn.Read([]byte(data))
		message := string(data[:])
		fmt.Println(message)
	}
}

func tcp_send(conn *net.TCPConn) {

	defer conn.Close()

	//message := "Hei"
	//data := make([]byte, message)

	//conn.Write(data)

	i := 0

	for {

		time.Sleep(100 * time.Millisecond)
		msg := strconv.Itoa(i) + "\x00"
		i++
		buf := []byte(msg)
		conn.Write(buf)

	}

}

func main() {

	local_addr, _ := net.ResolveTCPAddr("tcp", "129.241.187.155")
	server_addr, _ := net.ResolveTCPAddr("tcp", "129.241.187.23:33546")

	conn, _ := net.DialTCP("tcp", local_addr, server_addr)

	go tcp_receive(conn)
	go tcp_send(conn)

	for {
		time.Sleep(time.Second)
	}

}
