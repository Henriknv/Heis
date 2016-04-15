package network

import . ".././constants"

import (
	. "encoding/json"
	. "fmt"
	"net"
	"strconv"
	"strings"
)

// Ports and addresses for send/receive:

var Local_addr *net.UDPAddr
var Broadcast_addr *net.UDPAddr

var broadcast_listen_port int
var local_listen_port int

// Struct and channels for send/receive functions.

type MISO struct {
	Elev_id               string
	Local_order_matrix_two    [N_FLOORS][N_BUTTONS]int
	Local_cost_matrix     [N_FLOORS][N_BUTTONS]int
}

type MOSI struct {
	Elev_id               string
	Network_order_matrix [N_FLOORS][N_BUTTONS]int
	Master_cost_matrix   [N_FLOORS][N_BUTTONS]int
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

func Udp_send(mosiCh <-chan MOSI, misoCh <-chan MISO) {

	conn, _ := net.DialUDP("udp", Local_addr, Broadcast_addr)

	for {
		select {
		case miso := <-misoCh:

			buf, _ := Marshal(miso)
			//Println("UDP MISO Send Called")
			conn.Write([]byte("MISO" + string(buf)))
		case mosi := <-mosiCh:
			buf, _ := Marshal(mosi)
			//Println("UDP MOSI Send Called")
			conn.Write([]byte("MOSI" + string(buf)))
		}
		//Println("UDP Send Called")
	}
}

func Udp_receive(mosiCh chan<- MOSI, misoCh chan<- MISO) {

	conn, _ := net.ListenUDP("udp", Broadcast_addr)
	buf := make([]byte, 1024)

	for {
		n, _, _ := conn.ReadFromUDP(buf)
		if strings.HasPrefix(string(buf[:n]), "MISO") {
			var miso MISO
			Unmarshal(buf[4:n], &miso)
			misoCh <- miso
		} else if strings.HasPrefix(string(buf[:n]), "MOSI") {
			var mosi MOSI
			Unmarshal(buf[4:n], &mosi)
			mosiCh <- mosi
		}
	}
}

func Udp_init(local_listen_port int, broadcast_listen_port int, Slave_send_ch chan MISO, Slave_receive_ch chan MOSI, Master_receive chan MISO, Master_send chan MOSI) {

	get_broadcast_addr(broadcast_listen_port)
	get_local_addr(local_listen_port)

	go Udp_send(Master_send, Slave_send_ch)
	go Udp_receive(Slave_receive_ch, Master_receive)

}
