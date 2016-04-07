package elev

import . ".././driver"

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

func Calculate_cost(order_matrix [N_FLOORS][N_BUTTONS]int) [N_FLOORS][N_BUTTONS]int {

	dir := Elev_get_motor_direction()

	current_floor := <-Floor_sensor_chan

	cost_matrix := [N_FLOORS][N_BUTTONS]int{}

	var floor_dif int

	for i := 0; i < N_BUTTONS; i++ {

		for n := 0; n < N_FLOORS; n++ {

			if order_matrix[n][i] == 1 {

				floor_dif = n - current_floor

				if dir == DIR_UP {

					if floor_dif > 0 {

						cost_matrix[n][i] = cost_matrix[n][i] + (Abs_val(floor_dif) * COST_PER_FLOOR)

					} else if floor_dif < 0 {

						cost_matrix[n][i] = cost_matrix[n][i] + (Abs_val(floor_dif) * COST_PER_FLOOR) + TURN_COST

					}

				}

				if dir == DIR_DOWN {

					if floor_dif > 0 {

						cost_matrix[n][i] = cost_matrix[n][i] + (Abs_val(floor_dif) * COST_PER_FLOOR) + TURN_COST

					} else if floor_dif < 0 {

						cost_matrix[n][i] = cost_matrix[n][i] + (Abs_val(floor_dif) * COST_PER_FLOOR)

					}

				}

			}

		}

	}

	return cost_matrix

}
