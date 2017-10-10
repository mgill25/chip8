package main

import "fmt"

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
