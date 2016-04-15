package elev

import . ".././network"
import . ".././driver"
import . ".././constants"
//import . ".././fileio"

import (
	. "fmt"
	. "time"
)

var Local_order_matrix [N_FLOORS][N_BUTTONS]int
var Network_order_matrix [N_FLOORS][N_BUTTONS]int
var Master_cost_matrix [N_FLOORS][N_BUTTONS]int
var Master_order_matrix [N_FLOORS][N_BUTTONS]int

var external_orders [][]int
var Elev_orders []int
var Elev_costs []int

var INTERNAL_COSTS = 1
var INTERNAL_FLOORS = 0

var Floor_sensor_chan = make(chan int)
var Elev_id string

func Master_or_slave(Slave_input_ch <-chan MOSI) bool {

	var buf MOSI
	_ = buf

	var is_master bool

	select {

	case <-After(1 * Second):
		is_master = true

	case buf = <-Slave_input_ch:
		is_master = false

	}

	if is_master {

		Println("I am master.")
		Elev_id = Local_addr.String()[12:15]
		Println(Elev_id)

	} else {

		Println("I am slave.")
		Elev_id = Local_addr.String()[12:15]
		Println(Elev_id)

	}

	//Println(Network_order_matrix, "\n", Master_cost_matrix)

	return is_master
}

func Elevator_init() {

	Elev_init()

	Elev_set_motor_direction(0)

	//Local_order_matrix = Read()

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

		current_floor = Elev_get_floor_sensor_signal()

		if !(last_floor == current_floor) {

			last_floor = current_floor
		}

		select {

		case Floor_sensor_chan <- current_floor:
		default:

		}

		Sleep(10*Millisecond)
	}
}

var elev_dir int
var prev_dir int

func Execute_orders() {

	current_floor = <-Floor_sensor_chan
	var target_floor int

	for {

		Sort_orders(Master_cost_matrix)

		//Println(Master_cost_matrix)

		elev_dir = DIR_IDLE
		if len(Elev_orders) > 0 {

			target_floor = Elev_orders[0]

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

				Sort_orders(Master_cost_matrix)

				//Println(Master_cost_matrix)

				if current_floor != LIMBO {

					for i := 0; i < len(Elev_orders); i++ {

						if Elev_orders[i] == current_floor && current_floor != target_floor {

							Elev_stop_motor()
							Elev_open_door()
							delete_order(i, elev_dir)
							Elev_set_motor_direction(elev_dir)

						}
					}
				}
			}

			Elev_stop_motor()
			Elev_open_door()
			delete_order(0, elev_dir)
			//Println("TRACKING: ", "Orders: ", Elev_orders, "		Costs: ", Elev_costs)

			prev_dir = elev_dir

		}

		Sleep(10*Millisecond)
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

		Local_order_matrix[Elev_orders[order_index]][EXT_UP_BUTTONS] = 0
		Local_order_matrix[Elev_orders[order_index]][EXT_DOWN_BUTTONS] = 0
		Local_order_matrix[Elev_orders[order_index]][INTERNAL_BUTTONS] = 0
		
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

func Get_orders(Slave_output_ch chan MISO) [N_FLOORS][N_BUTTONS]int {

	var temp [N_FLOORS][N_BUTTONS]int
	var temp_slave_out MISO

	for {

		//temp = Local_order_matrix

		Get_internal_orders()
		Get_external_orders()

		if Local_order_matrix != temp {
			temp_slave_out.Elev_id = Elev_id
			temp_slave_out.Local_order_matrix = Local_order_matrix

			temp_slave_out.Local_cost_matrix = Calculate_cost()

			Slave_output_ch <- temp_slave_out
			//Println("id:  ", temp_slave_out.Elev_id)
			//Println("Local_order_matrix: ", temp_slave_out.Local_order_matrix)
			//Println("Local_cost_matrix: ", temp_slave_out.Local_cost_matrix)

			//Write(Local_order_matrix)

			temp = Local_order_matrix
			Println("Get_orders: ",Local_order_matrix)

		}
		Sleep(10*Millisecond)
	}
}

var sorting_bool bool

var copy_cost_matrix [N_FLOORS][N_BUTTONS]int

var old_cost_matrix [N_FLOORS][N_BUTTONS]int

func Sort_orders(copy_cost_matrix [N_FLOORS][N_BUTTONS]int) {

	

	if old_cost_matrix != copy_cost_matrix {

		//Println(copy_cost_matrix)

		for i := 0; i < N_FLOORS; i++ {

			for j := 0; j < N_BUTTONS; j++ {

				if copy_cost_matrix[i][j] != 0 {

					if len(Elev_orders) > 0 {

						sorting_bool = true

						for k := 0; k < len(Elev_orders); k++ {

							if Elev_orders[k] == i {
								if copy_cost_matrix[i][j] < Elev_costs[k] {

									Elev_costs[k] = copy_cost_matrix[i][j]

								}

								sorting_bool = false
							}
						}

						if sorting_bool == true {
						
							Elev_orders = append(Elev_orders, i)
							Elev_costs = append(Elev_costs, copy_cost_matrix[i][j])

						}

					} else if len(Elev_orders) == 0 {

						Elev_orders = append(Elev_orders, i)
						Elev_costs = append(Elev_costs, copy_cost_matrix[i][j])

					}
				}
			}
		}

		old_cost_matrix = copy_cost_matrix
		sequential_sort()

		Println("Orders: ", Elev_orders, "		Costs: ", Elev_costs)
	}
}

func sequential_sort() {

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
}

func Elev_lights() {

	current_floor := <-Floor_sensor_chan
	previous_floor := current_floor
	current_order := Master_order_matrix

	for {

		current_floor = <-Floor_sensor_chan

		if !(current_floor == previous_floor) {

			if current_floor >= 0 {

				Elev_set_floor_indicator(current_floor)
				previous_floor = current_floor

			}
		}

		if current_order != Master_order_matrix {
			for i := 0; i < N_BUTTONS; i++ {

				for j := 0; j < N_FLOORS; j++ {

					if Master_order_matrix[j][i] == 1 {

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

			current_order = Master_order_matrix

		}

		Sleep(10*Millisecond)
	}
}

func Master(Master_input_ch chan MISO, Master_output_ch chan MOSI) {

	//var temp_out MOSI

	var Network_order_matrix [N_FLOORS][N_BUTTONS]int

	var first_order bool
	var lowest_cost int
	var lowest_elev_id string

	//var temp_cost_matrix [N_FLOORS][N_BUTTONS]int

	Online_elevators := make(map[string]MISO)
	Online_elevators_prev := make(map[string]MISO)
	Slave_order_matrices := make(map[string]*[N_FLOORS][N_BUTTONS]int)

	for {

		select{

		case temp_in := <-Master_input_ch:

			Network_order_matrix = [N_FLOORS][N_BUTTONS]int{}
			Slave_order_matrices[temp_in.Elev_id] = &[N_FLOORS][N_BUTTONS]int{}
			Online_elevators[temp_in.Elev_id] = temp_in

			for i:=0; i < N_FLOORS; i++{

				for j:=0; j < N_BUTTONS; j++{

					for Elevator := range Online_elevators{

						if Online_elevators[Elevator].Local_order_matrix[i][j] == 1{

							Network_order_matrix[i][j] = 1

						}
					}

					if Online_elevators[temp_in.Elev_id].Local_order_matrix[i][j] < Online_elevators_prev[temp_in.Elev_id].Local_order_matrix[i][j]{

						Network_order_matrix[i][j] = 0
					
					}

					first_order = true
					lowest_cost = len(Online_elevators) * N_FLOORS * N_BUTTONS * 10

					for Elevator := range Online_elevators{

						if Online_elevators[Elevator].Local_cost_matrix[i][j] < lowest_cost && Network_order_matrix[i][j] == 1 {

							lowest_cost = Online_elevators[Elevator].Local_cost_matrix[i][j]
							Slave_order_matrices[Elevator][i][j] = lowest_cost
							
							if !first_order{

								Slave_order_matrices[lowest_elev_id][i][j] = 0
							
							}

							lowest_elev_id = Elevator
							first_order = false
							
						} 
					}
				}
			}

			for Elevator := range Online_elevators{

				//Println(MOSI{Elev_id: Elevator, Network_order_matrix: Network_order_matrix, Master_cost_matrix: *Slave_order_matrices[Elevator]})
				Master_output_ch <- MOSI{Elev_id: Elevator, Network_order_matrix: Network_order_matrix, Master_cost_matrix: *Slave_order_matrices[Elevator]}

			}			
		default:

		}

		Sleep(10*Millisecond)

	}
}

func Slave(Slave_input_ch chan MOSI) {

	for {

		select{
		
		case temp_in := <-Slave_input_ch:

			if temp_in.Elev_id == Elev_id {

				Master_cost_matrix = temp_in.Master_cost_matrix
				Master_order_matrix = temp_in.Network_order_matrix
				//Println("First print: ", Master_cost_matrix)

				for i := 0; i < N_FLOORS; i++ {
					
					if Local_order_matrix[i][INTERNAL_BUTTONS] == 1{
					
						Master_cost_matrix[i][INTERNAL_BUTTONS] = Calculate_cost()[i][INTERNAL_BUTTONS]
					
					}
					
					Master_order_matrix[i][INTERNAL_BUTTONS] = Local_order_matrix[i][INTERNAL_BUTTONS]

					for j := 0; j < N_BUTTONS-1; j++ {

						if Master_cost_matrix[i][j] > 0 {

							Local_order_matrix[i][j] = 1

						}
					}
				}

				//Println("Second print: ", Master_cost_matrix)

			}
		
		default:

		}
	}
}

func Spam(Master_output_ch chan MOSI){
	var temp MOSI
	temp.Elev_id = "000"
	for{

		Master_output_ch <- temp

		Sleep(500 * Millisecond)		
	}
}
