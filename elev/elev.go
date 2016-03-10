package elev

import . ".././network"
import . ".././driver"

import (
	. "fmt"
	//"./net"
	. "time"
)

/*const (
	FLOOR_1 = 0
	FLOOR_2 = 1
	FLOOR_3 = 2
	FLOOR_4 = 3
	LIMBO   = -1
)*/

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

	var current_floor int

	last_floor := Elev_get_floor_sensor_signal()

	for {

		//Floor has changed, sending event to channel.

		current_floor = Elev_get_floor_sensor_signal()

		if !(last_floor == current_floor) {

			Floor_sensor_chan <- current_floor
			last_floor = current_floor

		}
	}
}

func Execute_order(target_floor int) {

	var current_floor int

	current_floor = <-Floor_sensor_chan
	elev_dir := 0

	Println("Target floor:", target_floor)

	for !(current_floor == target_floor) {

		current_floor = <-Floor_sensor_chan
		Sleep(200 * Millisecond)

		Println("Current floor:", current_floor)
		Println("Elevator direction:", elev_dir)

		if current_floor < target_floor && !(current_floor == -1) {

			elev_dir = 1
			Elev_set_motor_direction(elev_dir)

		} else if current_floor > target_floor && !(current_floor == -1) {

			elev_dir = -1
			Elev_set_motor_direction(elev_dir)

		}

	}

	Elev_set_motor_direction(-elev_dir)
	Sleep(5 * Millisecond)

	elev_dir = 0
	Elev_set_motor_direction(elev_dir)

	Println("Order has been executed.")

}

/*
	if !current_floor == -1 {

		Elev_set_floor_indicator(current_floor)

	}
*/
