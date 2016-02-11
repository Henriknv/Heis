package network

import (
	"fmt"
	"net"
	"strconv"
	"time"
	."encoding/json"
)

var Local_addr *net.UDPAddr
var Broadcast_addr *net.UDPAddr

var broadcast_listen_port int
var local_listen_port int

// NEXT TIME: LAGE NY STRUCT, SOM ER BRUKANDES.

type Order struct {
	
	Text := "Heyo"
}

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

//TODO NESTE GANG: LAGE NYE SEND & RECEIVE FUNKSJONER. CHAN FOR SEND & RECV.

/*func Udp_receive(struct chan receive_chan) {

conn, _ := net.ListenUDP("udp", Broadcast_addr)

defer conn.Close()



for   {

time.Sleep(100 * time.Millisecond)
n, _, _ := conn.ReadFromUDP(buf)
fmt.Println("Received:", string(buf[0:n]))


}
buf := Unmarshal()
}

func Udp_send(struct chan send_chan) {

conn, _ := net.DialUDP("udp", Local_addr, Broadcast_addr)

defer conn.Close()

buf := jsoMarshal(message)

conn.Write(buf)

}*/

func Udp_init(local_listen_port int, broadcast_listen_port int) { //int message_size, send_ch, receive_ch, chan udp_message



	get_broadcast_addr(broadcast_listen_port)	// Setting up broadcast address.
	get_local_addr(local_listen_port)		// Setting up local address.

}
