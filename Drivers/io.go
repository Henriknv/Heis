package Drivers // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "elev.h"
*/
import "C"

func io_init() int {
	return int(C.io_init())
}
func io_set_bit(channel int) {
	C.io_set_bit(C.int(channel))
}
func io_clear_bit(channel int) {
	C.io_clear_bit(C.int(channel))
}
func io_write_analog(channel, value int) {
	C.io_write_analog(C.int(channel), C.int(value))
}
func io_read_bit(channel int) int {
	return int(C.io_read_bit(C.int(channel)))
}
func io_read_analog(channel int) int {
	return int(C.io_read_analog(C.int(channel)))
}

func elev_init() {
	C.elev_init()
}

func elev_set_motor_direction(elev_motor_direction_t dirn){
	C.elev_set_motor_direction(elev_motor_direction_t dirn)
}

func elev_set_button_lamp(elev_button_type_t button, int floor, int value){
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}

func elev_set_floor_indicator(int floor){
	C.elev_set_floor_indicator(C.int(floor))
}

func elev_set_door_open_lamp(int value){
	C.elev_set_door_open_lamp(C.int(value))
}

func elev_set_stop_lamp(int value){
	C.elev_set_stop_lamp(C.int(value))
}

func elev_get_button_signal(elev_button_type_t button, int floor) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor)))
}

func elev_get_floor_sensor_signal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func elev_get_stop_signal() int {
	return int(C.elev_get_stop_signal())
}

func elev_get_obstruction_signal() int {
	return int(C.elev_get_obstruction_signal())
}

