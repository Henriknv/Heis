package driver // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
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
