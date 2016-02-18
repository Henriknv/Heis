package network

import (
	. "encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

// Ports and addresses for send/receive:

var Local_addr *net.UDPAddr
var Broadcast_addr *net.UDPAddr

var broadcast_listen_port int
var local_listen_port int

// Struct and channels for send/receive functions.

type msg_struct struct {
	Text string
}

// Functions for aquiring broadcast and local addresses:

func get_broadcast_addr(broadcast_listen_port int) (err error) {

	Broadcast_addr, err = net.ResolveUDPAddr("udp", "129.241.187.255:"+strconv.Itoa(broadcast_listen_port))
	if err != nil {
		return err
	}
	fmt.Println("Printing broadcast address:" + Broadcast_addr.String())

	return

}

func get_local_addr(local_listen_port int) (err error) {

	temp_conn, err := net.DialUDP("udp", nil, Broadcast_addr)
	if err != nil {
		return err
	}
	defer temp_conn.Close()
	temp_addr := temp_conn.LocalAddr()
	Local_addr, err = net.ResolveUDPAddr("udp", temp_addr.String())

	Local_addr.Port = local_listen_port
	fmt.Println(Local_addr.String())
	fmt.Println("Printing local port:" + strconv.Itoa(Local_addr.Port))

	return

}

// Functions to send/receive a message via UDP:

func Udp_receive(Receieve_chan chan msg_struct) {

	var message msg_struct

	for {

		buf := make([]byte, 1024)

		conn, _ := net.ListenUDP("udp", Broadcast_addr)

		time.Sleep(100 * time.Millisecond)

		n, _, _ := conn.ReadFromUDP(buf)

		Unmarshal(buf[:n], &message)
		fmt.Println(message)

		Receieve_chan <- message

		conn.Close()
	}
}

func Udp_send(Send_chan chan msg_struct) {

	//var message msg_struct

	for {

		conn, _ := net.DialUDP("udp", Local_addr, Broadcast_addr)

		message := <-Send_chan

		buf, _ := Marshal(message)
		fmt.Println(buf)

		conn.Write(buf)

		conn.Close()
	}
}

func Udp_init(local_listen_port int, broadcast_listen_port int) { //int message_size, send_ch, receive_ch, chan udp_message

	get_broadcast_addr(broadcast_listen_port) // Setting up broadcast address.
	get_local_addr(local_listen_port)         // Setting up local address.

	Send_chan := make(chan msg_struct)
	Receive_chan := make(chan msg_struct)

	go Udp_send(Send_chan)
	go Udp_receive(Receive_chan)

	One_test := msg_struct{Text: "Halla"}

	Send_chan <- One_test

	//var two_test msg_struct

	Two_test := <-Receive_chan

	fmt.Println(Two_test.Text)

}
