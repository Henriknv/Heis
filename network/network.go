package network

//import . ".././constants"

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

// type Msg_struct_MISO struct {
// 	Elev_ID     string
// 	Cost_matrix [N_FLOORS][N_BUTTONS]int
// }

// type Msg_struct_MOSI struct {
// 	Elev_ID                string
// 	Order_array            [N_FLOORS]int
// 	External_lights_matrix [N_FLOORS][N_BUTTONS - 1]int
// }

type Msg_struct struct {
	Text string
}

var Send_chan chan Msg_struct
var Receive_chan chan Msg_struct

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

func Udp_receive(Receieve_chan chan Msg_struct) {

	var message Msg_struct

	for {

		buf := make([]byte, 1024)

		conn, _ := net.ListenUDP("udp", Broadcast_addr)

		time.Sleep(100 * time.Millisecond)

		n, _, _ := conn.ReadFromUDP(buf)

		Unmarshal(buf[:n], &message)

		Receieve_chan <- message

		conn.Close()
	}

}

func Udp_send(Send_chan chan Msg_struct) {

	for {

		conn, _ := net.DialUDP("udp", Local_addr, Broadcast_addr)

		message := <-Send_chan

		buf, _ := Marshal(message)

		conn.Write(buf)

		conn.Close()

	}

}

func Udp_init(local_listen_port int, broadcast_listen_port int) {

	get_broadcast_addr(broadcast_listen_port)
	get_local_addr(local_listen_port)

	go Udp_send(Send_chan)
	go Udp_receive(Receive_chan)

}
