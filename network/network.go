package network

import . ".././constants"

import (
	. "encoding/json"
	. "fmt"
	"net"
	"strconv"
//	"time"
)

// Ports and addresses for send/receive:

var Local_addr *net.UDPAddr
var Broadcast_addr *net.UDPAddr

var broadcast_listen_port int
var local_listen_port int

// Struct and channels for send/receive functions.


type Msg_struct struct {
	Elev_id string
	External_cost_matrix [N_FLOORS][N_BUTTONS]int
	Order_status bool
	Cost int
	Order int
}



// Functions for aquiring broadcast and local addresses:

func get_broadcast_addr(broadcast_listen_port int) (err error) {

	Broadcast_addr, err = net.ResolveUDPAddr("udp", "129.241.187.255:"+strconv.Itoa(broadcast_listen_port))

	if err != nil {
		return err
	}

	Println("Printing broadcast address:" + Broadcast_addr.String())

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

	Println(Local_addr.String())
	Println("Printing local port:" + strconv.Itoa(Local_addr.Port))

	return

}

// Functions to send/receive a message via UDP:

func Udp_receive(Receieve_chan chan Msg_struct) {

	var message Msg_struct

	for {

		buf := make([]byte, 1024)

		conn, _ := net.ListenUDP("udp", Broadcast_addr)

		//time.Sleep(100 * time.Millisecond)

		n, _, _ := conn.ReadFromUDP(buf)

		Unmarshal(buf[:n], &message)
		Println("received:", message)
		Receieve_chan <- message


		conn.Close()

	}

}
	
func Udp_send(Send_chan chan Msg_struct) {

	var message Msg_struct

	for {

		message = <- Send_chan

		conn, _ := net.DialUDP("udp", Local_addr, Broadcast_addr)
		
		buf, _ := Marshal(message)

		conn.Write(buf)

		conn.Close()
	}
}

func Udp_init(local_listen_port int, broadcast_listen_port int, send_chan chan Msg_struct, receive_chan chan Msg_struct) {

	get_broadcast_addr(broadcast_listen_port)
	get_local_addr(local_listen_port)

	go Udp_send(send_chan)
	go Udp_receive(receive_chan)

}
