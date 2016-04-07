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
const DIR_UP = 1
const DIR_DOWN = -1
const DIR_IDLE = 0
const LIMBO = -1
const N_BUTTONS = 4

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
		Elev_Id := Local_addr.String()[12:15]
		Println(Elev_Id)

	} else {

		Println("I am slave.")
		Elev_Id := Local_addr.String()[12:15]
		Println(Elev_Id)

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

	current_floor := <-Floor_sensor_chan
	elev_dir := DIR_IDLE

	for !(current_floor == target_floor) {

		if current_floor < target_floor && !(current_floor == -1) {

			elev_dir = DIR_UP
			Elev_set_motor_direction(elev_dir)

		} else if current_floor > target_floor && !(current_floor == -1) {

			elev_dir = DIR_DOWN
			Elev_set_motor_direction(elev_dir)

		}

		current_floor = <-Floor_sensor_chan

		if elev_dir == DIR_UP && current_floor != LIMBO {
			if external_orders[current_floor] == 1 {

				Execute_order(current_floor)
				external_orders[current_floor] = 0
				internal_orders[current_floor] = 0
			}
		}

		if elev_dir == DIR_DOWN && current_floor != LIMBO {

			if external_orders[current_floor+N_FLOORS] == 1 {

				Execute_order(current_floor)
				external_orders[current_floor+N_FLOORS] = 0
				internal_orders[current_floor] = 0
			}
		}

	}

	Elev_set_motor_direction(-elev_dir)
	Sleep(10 * Millisecond)

	elev_dir = DIR_IDLE
	Elev_set_motor_direction(elev_dir)
	Elev_set_door_open_lamp(1)
	Sleep(3000 * Millisecond)
	Elev_set_door_open_lamp(0)
	Sleep(430 * Millisecond)

}

var internal_orders [N_FLOORS]int
var external_orders [EXT_BUTTONS]int

func Get_internal_orders() {

	for i := 0; i < N_FLOORS; i++ {

		if Elev_get_button_signal(INTERNAL_BUTTONS, i) {

			internal_orders[i] = 1

		}
	}

}

func Get_external_orders() {

	for i := 0; i < N_FLOORS; i++ {

		if Elev_get_button_signal(EXT_UP_BUTTONS, i) {

			external_orders[i] = 1

		}

		if Elev_get_button_signal(EXT_DOWN_BUTTONS, i) {

			external_orders[i+N_FLOORS] = 1

		}
	}
}

func Get_orders() {
	for {
		Get_internal_orders()
		Get_external_orders()
	}
}

func Run_elevator() {

	internal := false

	for {

		for i := 0; i < N_FLOORS; i++ {

			if internal_orders[i] == 1 {

				Execute_order(i)

				internal_orders[i] = 0

				internal = true

			}
		}

		if !internal {
			for i := 0; i < N_FLOORS; i++ {

				if external_orders[i] == 1 {

					Execute_order(i)

					external_orders[i] = 0

				}

				if external_orders[i+N_FLOORS] == 1 {

					Execute_order(i)

					external_orders[i+N_FLOORS] = 0

				}
			}
		}

		internal = false
	}
}

func Elev_lights() {

	current_floor := <-Floor_sensor_chan
	previous_floor := current_floor
	current_internal := internal_orders
	current_external := external_orders

	for {

		current_floor = <-Floor_sensor_chan
		if !(current_floor == previous_floor) {
			if current_floor >= 0 {
				Elev_set_floor_indicator(current_floor)
				previous_floor = current_floor
			}
		}

		if current_internal != internal_orders {
			for i := 0; i < N_FLOORS; i++ {

				if internal_orders[i] == 1 {
					Elev_set_button_lamp(INTERNAL_BUTTONS, i, 1)
				} else {
					Elev_set_button_lamp(INTERNAL_BUTTONS, i, 0)
				}
			}

			current_internal = internal_orders
		}

		if current_external != external_orders {
			for i := 0; i < N_FLOORS; i++ {

				if external_orders[i] == 1 {
					Elev_set_button_lamp(EXT_UP_BUTTONS, i, 1)
				} else {
					Elev_set_button_lamp(EXT_UP_BUTTONS, i, 0)
				}

				if external_orders[i+N_FLOORS] == 1 {
					Elev_set_button_lamp(EXT_DOWN_BUTTONS, i, 1)
				} else {
					Elev_set_button_lamp(EXT_DOWN_BUTTONS, i, 0)
				}
			}
			current_external = external_orders
		}

	}

}
