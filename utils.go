package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// Check Generic utility function to check for errors
func Check(e error) {
	if e != nil {
		fmt.Println("Got error:", e)
	}
}

// B2i Convert a boolean value to an int
func B2i(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

// Opcode parsing helpers

// ArgNNN Return the NNN address value
func (chip *Chip8) ArgNNN(opcode uint16) uint16 {
	return opcode & 0x0fff
}

// ArgNN Return the NN 8-bit constant
func (chip *Chip8) ArgNN(opcode uint16) uint16 {
	return opcode & 0x00ff
}

// ArgN Return the N 4-bit constant
func (chip *Chip8) ArgN(opcode uint16) uint16 {
	return opcode & 0x000f
}

// ArgX Return the 4-bit register identifier X
// X is at the 2nd byte position. Example: 8XY1
func (chip *Chip8) ArgX(opcode uint16) uint16 {
	return (opcode & 0x0f00) >> 8
}

// ArgY Return the 4-bit register identifier
// X is at the 3rd byte position. Example: 8XY1
func (chip *Chip8) ArgY(opcode uint16) uint16 {
	return (opcode & 0x00f0) >> 4
}

func TranslateCodeToIndex(key sdl.Keycode) int {
	switch key {
	case sdl.K_0:
		return 0
	case sdl.K_1:
		return 1
	case sdl.K_2:
		return 2
	case sdl.K_3:
		return 3
	case sdl.K_4:
		return 4
	case sdl.K_5:
		return 5
	case sdl.K_6:
		return 6
	case sdl.K_7:
		return 7
	case sdl.K_8:
		return 8
	case sdl.K_9:
		return 9
	case sdl.K_a:
		return 10
	case sdl.K_b:
		return 11
	case sdl.K_c:
		return 12
	case sdl.K_d:
		return 13
	case sdl.K_e:
		return 14
	case sdl.K_f:
		return 15

	case sdl.K_KP_0:
		return 0
	case sdl.K_KP_1:
		return 1
	case sdl.K_KP_2:
		return 2
	case sdl.K_KP_3:
		return 3
	case sdl.K_KP_4:
		return 4
	case sdl.K_KP_5:
		return 5
	case sdl.K_KP_6:
		return 6
	case sdl.K_KP_7:
		return 7
	case sdl.K_KP_8:
		return 8
	case sdl.K_KP_9:
		return 9

	case sdl.K_DOWN:
		return 8
	case sdl.K_UP:
		return 2
	case sdl.K_LEFT:
		return 4
	case sdl.K_RIGHT:
		return 6
	}

	return -1
}

// Stupid Golang
// https://mrekucci.blogspot.in/2015/07/dont-abuse-mathmax-mathmin.html
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
