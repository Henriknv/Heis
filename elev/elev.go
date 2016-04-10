package elev

import . ".././network"
import . ".././driver"
import . ".././constants"

import (
	. "fmt"
	. "time"
)

var Local_order_matrix [N_FLOORS][N_BUTTONS]int
var external_orders [][]int
var Elev_orders []int
var Elev_costs []int

var INTERNAL_COSTS = 1
var INTERNAL_FLOORS = 0

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

// func Master(){

// 	var elev_Id_best string
// 	var elev_Id_start string
// 	var best_value[N_FLOORS][N_BUTTONS] int
// 	var best_matrix [N_FLOORS][N_BUTTONS]int

// 	var received Msg_struct_MISO

// 	for{

// 		select {
// 			case received <- Receive_chan:
// 			default:
// 		}

// 		elev_Id_received = received.Elev_ID

// 		for i := 0; i < N_FLOORS; i++{
// 			for j := 0; j<N_BUTTONS-1; j++{
// 				if  received.Cost_matrix[i][j] < best_value[i][j] && received.Cost_matrix[i][j] != 0{
// 					best_matrix[i][j] = received.Elev_ID
// 					best_value[i][j] = received.Cost_matrix
// 				}
// 			}
// 		}
// 	}

// 	for i := 0; i < N_FLOORS; i++{
// 			for j := 0; j<N_BUTTONS-1; j++{

// 				}
// 			}
// 		}

// }

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

func Execute_orders() {

	current_floor := <-Floor_sensor_chan

	elev_dir := DIR_IDLE

	Println("TISS")

	for {

		if len(Elev_orders) > 0 {

			target_floor := Elev_orders[0]

			for !(current_floor == target_floor) {

				if current_floor < target_floor && (current_floor != -1) {

					elev_dir = DIR_UP
					Elev_set_motor_direction(elev_dir)

				} else if current_floor > target_floor && (current_floor != -1) {

					elev_dir = DIR_DOWN
					Elev_set_motor_direction(elev_dir)

				}

				current_floor = <-Floor_sensor_chan

				// if elev_dir == DIR_UP && current_floor != LIMBO {
				// 	if external_orders[current_floor] == 1 {

				// 		Execute_order(current_floor)
				// 		external_orders[current_floor] = 0
				// 		internal_orders[current_floor] = 0
				// 	}
				// }

				// if elev_dir == DIR_DOWN && current_floor != LIMBO {

				// 	if external_orders[current_floor+N_FLOORS] == 1 {

				// 		Execute_order(current_floor)
				// 		external_orders[current_floor+N_FLOORS] = 0
				// 		internal_orders[current_floor] = 0
				// 	}
				// }
			}

			Elev_set_motor_direction(-elev_dir)
			Sleep(10 * Millisecond)
			elev_dir = DIR_IDLE
			Elev_set_motor_direction(elev_dir)

			Elev_set_door_open_lamp(1)
			Sleep(3000 * Millisecond)
			Elev_set_door_open_lamp(0)
			Sleep(500 * Millisecond)

		}
	}
}

func Get_internal_orders() {

	for i := 0; i < N_FLOORS; i++ {

		if Elev_get_button_signal(INTERNAL_BUTTONS, i) {

			Local_order_matrix[i][INTERNAL_BUTTONS] = 1

		}
	}
}

func Get_external_orders() {

	for i := 0; i < N_FLOORS; i++ {

		if Elev_get_button_signal(EXT_UP_BUTTONS, i) {

			Local_order_matrix[i][EXT_UP_BUTTONS] = 1

		}

		if Elev_get_button_signal(EXT_DOWN_BUTTONS, i) {

			Local_order_matrix[i][EXT_DOWN_BUTTONS] = 1

		}
	}
}

func Get_orders() [N_FLOORS][N_BUTTONS]int {

	for {

		Get_internal_orders()
		Get_external_orders()
		Sort_orders()
	}
}

var sorting_bool bool

func Sort_orders() {

	costs := Calculate_cost(Local_order_matrix)
	//Println("Costs:", costs)

	for i := 0; i < N_FLOORS; i++ {

		if costs[i][INTERNAL_BUTTONS] != 0 {

			Println("Hei")

			if len(Elev_orders) > 0 {

				Println("BÆSJ")
				sorting_bool = true

				for j := 0; j < len(Elev_orders); j++ {

					if Elev_orders[j] == i && Elev_costs[j] == costs[i][INTERNAL_COSTS] {
						sorting_bool = false
					}

				}

				if sorting_bool {

					Println("This is as far as I GO.")
					Elev_orders = append(Elev_orders, i)
					Elev_costs = append(Elev_costs, costs[i][INTERNAL_BUTTONS])

				}

			} else if len(Elev_orders) == 0 {

				Elev_orders = append(Elev_orders, i)
				Elev_costs = append(Elev_costs, costs[i][INTERNAL_BUTTONS])

			}
		}
	}

	// if len(external_orders) > 0 {

	// 	for i := 0; i < len(external_orders); i++ {

	// 		temp_orders = append(temp_orders, external_orders[i][INTERNAL_FLOORS])
	// 		temp_costs = append(temp_costs, external_orders[i][INTERNAL_COSTS])

	// 	}

	// }

	//Elev_orders = merge(temp_orders, temp_costs)

	Println("PROMP")

	sort()
}

// func merge(array_one []int, array_two []int) [12][2]int {

// 	temp_len := len(array_one)
// 	Println(temp_len)

// 	var merged_array [12][2]int

// 	for i := 0; i < temp_len; i++ {

// 		merged_array[i][INTERNAL_FLOORS] = array_one[i]
// 		merged_array[i][INTERNAL_COSTS] = array_two[i]

// 	}

// 	return merged_array
// }

func sort() {

	counter := 1

	var temp_floor int
	var temp_cost int

	for counter > 0 {

		counter = 0

		for i := 0; i < len(Elev_costs)-1; i++ {

			if Elev_costs[i+1] < Elev_costs[i] {

				temp_cost = Elev_costs[i]
				Elev_costs[i] = Elev_costs[i+1]
				Elev_costs[i+1] = temp_cost

				temp_floor = Elev_orders[i]
				Elev_orders[i] = Elev_orders[i+1]
				Elev_orders[i+1] = temp_floor

				counter = counter + 1

			}
		}
	}

	Println("Elev_orders: ", Elev_orders, "Elev_costs: ", Elev_costs)

}

// func Run_elevator() {

// 	internal := false

// 	for {

// 		for i := 0; i < N_FLOORS; i++ {

// 			if internal_orders[i] == 1 {

// 				Execute_order(i)

// 				internal_orders[i] = 0

// 				internal = true

// 			}
// 		}

// 		if !internal {
// 			for i := 0; i < N_FLOORS; i++ {

// 				if external_orders[i] == 1 {

// 					Execute_order(i)

// 					external_orders[i] = 0

// 				}

// 				if external_orders[i+N_FLOORS] == 1 {

// 					Execute_order(i)

// 					external_orders[i+N_FLOORS] = 0

// 				}
// 			}
// 		}

// 		internal = false
// 	}
// }

func Elev_lights() {

	current_floor := <-Floor_sensor_chan
	previous_floor := current_floor
	current_order := Local_order_matrix

	for {

		current_floor = <-Floor_sensor_chan

		if !(current_floor == previous_floor) {

			if current_floor >= 0 {

				Elev_set_floor_indicator(current_floor)
				previous_floor = current_floor

			}
		}

		if current_order != Local_order_matrix {
			for i := 0; i < N_BUTTONS; i++ {

				for j := 0; j < N_FLOORS; j++ {

					if Local_order_matrix[j][i] == 1 {

						if i == INTERNAL_BUTTONS {
							Elev_set_button_lamp(INTERNAL_BUTTONS, j, 1)
						}
						if i == EXT_UP_BUTTONS {
							Elev_set_button_lamp(EXT_UP_BUTTONS, j, 1)
						}

						if i == EXT_DOWN_BUTTONS {
							Elev_set_button_lamp(EXT_DOWN_BUTTONS, j, 1)
						}

					} else {

						if i == INTERNAL_BUTTONS {
							Elev_set_button_lamp(INTERNAL_BUTTONS, j, 0)
						}

						if i == EXT_UP_BUTTONS {
							Elev_set_button_lamp(EXT_UP_BUTTONS, j, 0)
						}

						if i == EXT_DOWN_BUTTONS {
							Elev_set_button_lamp(EXT_DOWN_BUTTONS, j, 0)
						}
					}
				}
			}
		}
	}
}
