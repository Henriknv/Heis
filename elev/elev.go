package elev

import . ".././network"
import . ".././driver"

import (
	. "fmt"
	//"./net"
	. "time"
)

const INTERNAL_BUTTONS int = 2
const N_FLOORS int = 4
const EXT_DOWN_BUTTONS int = 1
const EXT_UP_BUTTONS int = 0
const EXT_BUTTONS = N_FLOORS * 2

var Floor_sensor_chan = make(chan int)

func Master_or_slave() bool {

	var buf Msg_struct
	_ = buf

	var is_master bool

	select {

	case <-After(1 * Second):
		is_master = true

	case buf = <-Send_chan:
		is_master = false

	}

	if is_master {

		Println("I am master.")

	} else {

		Println("I am slave.")

	}

	return is_master
}

func Elevator_init() {

	Elev_init()

	if Elev_get_floor_sensor_signal() == -1 {
		Elev_set_motor_direction(-1)
		for {
			if Elev_get_floor_sensor_signal() > -1 {
				Elev_set_motor_direction(0)
				Elev_set_floor_indicator(Elev_get_floor_sensor_signal())
				break
			}
		}
	}
}

func Elev_maintenance() {

	current_floor := Elev_get_floor_sensor_signal()

	last_floor := current_floor

	Floor_sensor_chan <- current_floor

	for {

		//Floor has changed, sending event to channel.

		current_floor = Elev_get_floor_sensor_signal()

		if !(last_floor == current_floor) {

				last_floor = current_floor

		}

		select {

		case Floor_sensor_chan <- current_floor:
		default:

		}
	}
}

func Execute_order(target_floor int) {

	var current_floor int

	current_floor = <-Floor_sensor_chan
	elev_dir := 0

	for !(current_floor == target_floor) {

		//Println("Current floor:", current_floor)
		//Println("Elevator direction:", elev_dir)

		if current_floor < target_floor && !(current_floor == -1) {

			elev_dir = 1
			Elev_set_motor_direction(elev_dir)

		} else if current_floor > target_floor && !(current_floor == -1) {

			elev_dir = -1
			Elev_set_motor_direction(elev_dir)

		}

		current_floor = <-Floor_sensor_chan

	}

	Elev_set_motor_direction(-elev_dir)
	Sleep(10 * Millisecond)

	elev_dir = 0
	Elev_set_motor_direction(elev_dir)

}

var internal_orders [N_FLOORS]int
var external_orders [EXT_BUTTONS]int

func Get_internal_orders(){

	for i := 0; i < N_FLOORS; i++{
			
		if Elev_get_button_signal(INTERNAL_BUTTONS, i){

			internal_orders[i] = 1

		}
	}

	Sleep(100 * Millisecond)

	Println(internal_orders)

}

func Get_external_orders(){

	for i := 0; i<N_FLOORS; i++{

		if Elev_get_button_signal(EXT_UP_BUTTONS, i){

		external_orders[i] = 1

		}

		if Elev_get_button_signal(EXT_DOWN_BUTTONS, i){

		external_orders[i+N_FLOORS] = 1

		}		
	}
}

func Get_orders(){
	for{
		Get_internal_orders()
		Get_external_orders()
	}
}

func Run_elevator(){

	

	for{

		for i:= 0; i < N_FLOORS; i++{

			if external_orders[i] == 1{

			Execute_order(i)

			external_orders[i] = 0

			}

			if external_orders[i+N_FLOORS] == 1{

			Execute_order(i)

			external_orders[i+N_FLOORS] = 0

			}
		}
	}
}

func Elev_lights(){

	//Henrik fikk på et grønt lys. Dagen har vært bra. Sterke følelser. Mye tårer. Blasfemi. Heisen er ikke fornøyd.

}