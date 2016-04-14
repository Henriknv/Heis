package elev

//import . ".././driver"
import . ".././constants"

//import . "fmt"

const TURN_COST int = 35
const COST_PER_FLOOR int = 10

var order_matrix = [N_FLOORS][N_BUTTONS]int{}

//{EXT_UP1, EXT_DOWN1, INT_1},
//{EXT_UP2, EXT_DOWN2, INT_2},
//{EXT_UP3, EXT_DOWN3, INT_3},
//{EXT_UP4, EXT_DOWN4, INT_4},

func Abs_val(val int) int {

	if val < 0 {
		return -val
	}
	return val

}

var dir int
var current_floor int
var prev_floor int

var cost_matrix [N_FLOORS][N_BUTTONS]int

func Calculate_cost(order_matrix [N_FLOORS][N_BUTTONS]int) [N_FLOORS][N_BUTTONS]int {

	

	current_floor = <-Floor_sensor_chan

	if current_floor != LIMBO {
		prev_floor = current_floor
	}

	cost_matrix = [N_FLOORS][N_BUTTONS]int{}

	var floor_dif int

	for i := 0; i < N_BUTTONS; i++ {

		for n := 0; n < N_FLOORS; n++ {

			if order_matrix[n][i] == 1 {

				floor_dif = n - prev_floor

				if elev_dir == DIR_UP {

					if floor_dif > 0 {

						cost_matrix[n][i] = (Abs_val(floor_dif) * COST_PER_FLOOR) + 1

					} else if floor_dif < 0 {

						cost_matrix[n][i] = (Abs_val(floor_dif) * COST_PER_FLOOR) + TURN_COST

					}

				} else if elev_dir == DIR_DOWN {

					if floor_dif > 0 {

						cost_matrix[n][i] = (Abs_val(floor_dif) * COST_PER_FLOOR) + TURN_COST

					} else if floor_dif < 0 {

						cost_matrix[n][i] = (Abs_val(floor_dif) * COST_PER_FLOOR) + 1

					}

				} else {

					cost_matrix[n][i] = (Abs_val(floor_dif) * COST_PER_FLOOR) + 1

				}

			}

		}

	}

	//Println(cost_matrix)
	return cost_matrix

}
