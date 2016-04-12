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

	Elev_set_motor_direction(0)

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

// func Execute_orders() {

// 	//current_floor := <-Floor_sensor_chan

// 	current_floor := Elev_get_floor_sensor_signal()

// 	elev_dir := DIR_IDLE

// 	var target_floor int

// 	for {

// 		//Println("Current floor:", current_floor)
// 		//Println("Len orders:", len(Elev_orders))

// 		if len(Elev_orders) > 0 {

// 			target_floor = Elev_orders[0]

// 			if current_floor < target_floor && (current_floor != -1) {

// 				elev_dir = DIR_UP
// 				Elev_set_motor_direction(elev_dir)

// 			} else if current_floor > target_floor && (current_floor != -1) {

// 				elev_dir = DIR_DOWN
// 				Elev_set_motor_direction(elev_dir)

// 			}

// 			for current_floor != target_floor {

// 				current_floor = <-Floor_sensor_chan

// 				for i := 0; i < len(Elev_orders); i++ {
// 					//Println("Current floor1:   ", current_floor)
// 					if current_floor == Elev_orders[i] && current_floor != target_floor {

// 						if Local_order_matrix[Elev_orders[i]][INTERNAL_BUTTONS] == 1 {

// 							Elev_stop_motor()
// 							Elev_open_door()
// 							Println("Delete: 1")
// 							delete_order(i, elev_dir)

// 						} else if Local_order_matrix[Elev_orders[i]][EXT_UP_BUTTONS] == 1 && elev_dir == DIR_UP {

// 							Elev_stop_motor()
// 							Elev_open_door()
// 							delete_order(i, elev_dir)
// 						} else if Local_order_matrix[Elev_orders[i]][EXT_DOWN_BUTTONS] == 1 && elev_dir == DIR_DOWN {

// 							Elev_stop_motor()
// 							Elev_open_door()
// 							delete_order(i, elev_dir)
// 						}

// 					}

// 				}

// 				Elev_set_motor_direction(elev_dir)
// 			}

// 			Elev_stop_motor()
// 			Elev_open_door()
// 			delete_order(0, elev_dir)
// 			Println("Delete: 2")
// 			elev_dir = DIR_IDLE

// 		}
// 	}
// }

// func delete_order(order_index int, last_dir int) {
// 	//Println(Elev_orders, order_index)

// 	if last_dir == DIR_UP {
// 		Local_order_matrix[Elev_orders[order_index]][EXT_UP_BUTTONS] = 0
// 		Local_order_matrix[Elev_orders[order_index]][INTERNAL_BUTTONS] = 0

// 		copy_cost_matrix[Elev_orders[order_index]][EXT_UP_BUTTONS] = 0
// 		copy_cost_matrix[Elev_orders[order_index]][INTERNAL_BUTTONS] = 0
// 	} else if last_dir == DIR_DOWN {
// 		Local_order_matrix[Elev_orders[order_index]][EXT_DOWN_BUTTONS] = 0
// 		Local_order_matrix[Elev_orders[order_index]][INTERNAL_BUTTONS] = 0

// 		copy_cost_matrix[Elev_orders[order_index]][EXT_DOWN_BUTTONS] = 0
// 		copy_cost_matrix[Elev_orders[order_index]][INTERNAL_BUTTONS] = 0
// 	}
// 	if order_index == 0 || order_index == N_FLOORS-1 {
// 		Local_order_matrix[Elev_orders[order_index]][EXT_DOWN_BUTTONS] = 0
// 		copy_cost_matrix[Elev_orders[order_index]][EXT_DOWN_BUTTONS] = 0

// 		Local_order_matrix[Elev_orders[order_index]][EXT_UP_BUTTONS] = 0
// 		copy_cost_matrix[Elev_orders[order_index]][EXT_UP_BUTTONS] = 0
// 	}
// 	if len(Elev_costs) > 1 {
// 		Elev_orders = append(Elev_orders[:order_index], Elev_orders[order_index+1:]...)
// 		Elev_costs = append(Elev_costs[:order_index], Elev_costs[order_index+1:]...)
// 	} else if len(Elev_costs) == 1 {
// 		Elev_orders = Elev_orders[1:len(Elev_orders)]
// 		Elev_costs = Elev_costs[1:len(Elev_costs)]
// 		Println(len(Elev_costs))
// 	}
// }

var elev_dir int

func Execute_orders() {
	current_floor = <-Floor_sensor_chan

	const FIRST_ELEMENT = 0
	checked_floor = current_floor
	var target_floor int

	for {

		Sort_orders()

		if len(Elev_orders) > 0 {

			target_floor = Elev_orders[FIRST_ELEMENT]
			Println(target_floor)

			if target_floor > current_floor {
				elev_dir = DIR_UP

			} else if target_floor < current_floor {
				elev_dir = DIR_DOWN

			} else {
				elev_dir = DIR_IDLE
			}

			Elev_set_motor_direction(elev_dir)

			for current_floor != target_floor {

				current_floor = <-Floor_sensor_chan
				Sort_orders()

				if current_floor != LIMBO {

					for i := 0; i < len(Elev_orders); i++ {

						if Elev_orders[i] == current_floor && current_floor != target_floor {

							Elev_stop_motor()
							Elev_open_door()
							Println("Delete 2")
							delete_order(i, elev_dir)
							Elev_set_motor_direction(elev_dir)

						}
					}
				}
			}

			Elev_stop_motor()
			Elev_open_door()
			Println("Delete 1")
			delete_order(FIRST_ELEMENT, elev_dir)

		}
	}
}

func delete_order(order_index int, dir int) {

	if dir == DIR_UP {
		Local_order_matrix[Elev_orders[order_index]][INTERNAL_BUTTONS] = 0
		Local_order_matrix[Elev_orders[order_index]][EXT_UP_BUTTONS] = 0

	} else if dir == DIR_DOWN {

		Local_order_matrix[Elev_orders[order_index]][INTERNAL_BUTTONS] = 0
		Local_order_matrix[Elev_orders[order_index]][EXT_DOWN_BUTTONS] = 0

	} else if dir == DIR_IDLE {

		Local_order_matrix[Elev_orders[order_index]][INTERNAL_BUTTONS] = 0
		Local_order_matrix[Elev_orders[order_index]][EXT_DOWN_BUTTONS] = 0
		Local_order_matrix[Elev_orders[order_index]][EXT_UP_BUTTONS] = 0

	}

	if Elev_orders[order_index] == 0 || Elev_orders[order_index] == N_FLOORS-1 {

		Local_order_matrix[Elev_orders[order_index]][EXT_UP_BUTTONS] = 0
		Local_order_matrix[Elev_orders[order_index]][EXT_DOWN_BUTTONS] = 0

	}

	if len(Elev_orders) > 0 {

		Elev_orders = append(Elev_orders[:order_index], Elev_orders[order_index+1:]...)
		Elev_costs = append(Elev_costs[:order_index], Elev_costs[order_index+1:]...)

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

	}
}

var sorting_bool bool

var copy_cost_matrix [N_FLOORS][N_BUTTONS]int

var old_cost_matrix [N_FLOORS][N_BUTTONS]int

var checked_floor int

func Sort_orders() {

	copy_cost_matrix = Calculate_cost(Local_order_matrix)

	current_floor = <-Floor_sensor_chan

	if old_cost_matrix != copy_cost_matrix {

		for i := 0; i < N_FLOORS; i++ {

			if copy_cost_matrix[i][INTERNAL_BUTTONS] != 0 {

				if len(Elev_orders) > 0 {

					sorting_bool = true
					//Println("Ute: ", Elev_orders)

					for j := 0; j < len(Elev_orders); j++ {
						//Println("Inne: ", Elev_orders)
						//Println("I: ", i, " J: ", j, "  Cost: ", copy_cost_matrix[i][INTERNAL_BUTTONS])
						if Elev_orders[j] == i {
							sorting_bool = false
							Elev_costs[j] = copy_cost_matrix[i][INTERNAL_BUTTONS]

						}

					}

					if sorting_bool == true {

						Elev_orders = append(Elev_orders, i)
						Elev_costs = append(Elev_costs, copy_cost_matrix[i][INTERNAL_BUTTONS])

					}

				} else if len(Elev_orders) == 0 {

					Elev_orders = append(Elev_orders, i)
					Elev_costs = append(Elev_costs, copy_cost_matrix[i][INTERNAL_BUTTONS])

				}
			}
		}

		old_cost_matrix = copy_cost_matrix
	}

	// if len(external_orders) > 0 {

	// 	for i := 0; i < len(external_orders); i++ {

	// 		temp_orders = append(temp_orders, external_orders[i][INTERNAL_FLOORS])
	// 		temp_costs = append(temp_costs, external_orders[i][INTERNAL_COSTS])

	// 	}

	// }

	//Elev_orders = merge(temp_orders, temp_costs)
	//Println("Inn: ", Elev_costs, "   ordrs  ", Elev_orders)
	sort()
	//Println("UT:  ", Elev_costs, "   orders: ", Elev_orders)
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

	//Println("Elev_orders: ", Elev_orders, "Elev_costs: ", Elev_costs)

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
			current_order = Local_order_matrix
		}
	}
}
