package driver // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "elev.h"
*/
import "C"
import . "time"

func Io_init() int {
	return int(C.io_init())
}
func Io_set_bit(channel int) {
	C.io_set_bit(C.int(channel))
}
func Io_clear_bit(channel int) {
	C.io_clear_bit(C.int(channel))
}
func Io_write_analog(channel, value int) {
	C.io_write_analog(C.int(channel), C.int(value))
}
func Io_read_bit(channel int) int {
	return int(C.io_read_bit(C.int(channel)))
}
func Io_read_analog(channel int) int {
	return int(C.io_read_analog(C.int(channel)))
}

func Elev_init() {
	C.elev_init()
}

func Elev_set_motor_direction(dirn int) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(dirn))
}

func Elev_get_motor_direction() int {
	return int(C.elev_get_motor_direction())
}

func Elev_stop_motor() {
	Elev_set_motor_direction(0)
}

func Elev_set_button_lamp(button int, floor int, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}

func Elev_set_floor_indicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func Elev_set_door_open_lamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}

func Elev_set_stop_lamp(value int) {
	C.elev_set_stop_lamp(C.int(value))
}

func Elev_get_button_signal(button int, floor int) bool {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor))) != 0
}

func Elev_get_floor_sensor_signal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func Elev_get_stop_signal() int {
	return int(C.elev_get_stop_signal())
}

func Elev_get_obstruction_signal() int {
	return int(C.elev_get_obstruction_signal())
}

func Elev_open_door() {
	Elev_set_door_open_lamp(1)
	Sleep(3000 * Millisecond)
	Elev_set_door_open_lamp(0)
}
