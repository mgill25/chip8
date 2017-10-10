package main

import (
	"math"
	"math/rand"
)

// Opcode Definitions

func (chip8 *Chip8) Opcode0NNN(opcode uint16) {
}

func (chip8 *Chip8) Opcode00E0(opcode uint16) {
}

func (chip8 *Chip8) Opcode00EE(opcode uint16) {
}

// Opcode1NNN Jumps to address NNN
func (chip8 *Chip8) Opcode1NNN(opcode uint16) {
	chip8.PC = chip8.ArgNNN(opcode)
}

func (chip8 *Chip8) Opcode2NNN(opcode uint16) {
}

// Opcode3XNN Skips the next instruction if Vx = NN
func (chip8 *Chip8) Opcode3XNN(opcode uint16) {
	// Assumption is working the same as IP (pointing towards the next instruction)
	// So if we want to skip the next instruction, we just add 2 bytes (since each opcode is 2 bytes)
	// https://en.wikipedia.org/wiki/Program_counter
	if uint16(chip8.V[chip8.ArgX(opcode)]) == chip8.ArgNN(opcode) {
		chip8.PC += 2
	}
}

func (chip8 *Chip8) Opcode4XNN(opcode uint16) {
	if uint16(chip8.V[chip8.ArgX(opcode)]) != chip8.ArgNN(opcode) {
		chip8.PC += 2
	}
}

func (chip8 *Chip8) Opcode5XY0(opcode uint16) {
	if chip8.V[chip8.ArgX(opcode)] == chip8.V[chip8.ArgY(opcode)] {
		chip8.PC += 2
	}
}

// Opcode6XNN Vx = NN
func (chip8 *Chip8) Opcode6XNN(opcode uint16) {
	chip8.V[chip8.ArgX(opcode)] = byte(chip8.ArgNN(opcode))
}

// Opcode7XNN Adds Vx to NN (carry flag is not changed).
// XXX: How does a computer check for integer overflow?
// https://en.wikipedia.org/wiki/Integer_overflow
func (chip8 *Chip8) Opcode7XNN(opcode uint16) {
	chip8.V[chip8.ArgX(opcode)] += chip8.V[chip8.ArgNN(opcode)]
}

// Opcode8XY0 Vx = Vy
func (chip8 *Chip8) Opcode8XY0(opcode uint16) {
	chip8.V[chip8.ArgX(opcode)] = chip8.V[chip8.ArgY(opcode)]
}

// Opcode8XY1 Vx = Vx | Vy
func (chip8 *Chip8) Opcode8XY1(opcode uint16) {
	Vx := chip8.V[chip8.ArgX(opcode)]
	Vy := chip8.V[chip8.ArgY(opcode)]
	chip8.V[chip8.ArgX(opcode)] = Vx | Vy
}

// Opcode8XY2 Vx = Vx & Vy
func (chip8 *Chip8) Opcode8XY2(opcode uint16) {
	Vx := chip8.V[chip8.ArgX(opcode)]
	Vy := chip8.V[chip8.ArgY(opcode)]
	chip8.V[chip8.ArgX(opcode)] = Vx & Vy
}

// Opcode8XY3 Vx = Vx ^ Vy
func (chip8 *Chip8) Opcode8XY3(opcode uint16) {
	Vx := chip8.V[chip8.ArgX(opcode)]
	Vy := chip8.V[chip8.ArgY(opcode)]
	chip8.V[chip8.ArgX(opcode)] = Vx ^ Vy
}

// Opcode8XY4 Vx += Vy. VF = 1 when there is a carry, otherwise 0.
func (chip8 *Chip8) Opcode8XY4(opcode uint16) {
	Vx := chip8.V[chip8.ArgX(opcode)]
	Vy := chip8.V[chip8.ArgY(opcode)]
	chip8.V[chip8.ArgX(opcode)] = Vx + Vy
	chip8.V[0xf] = byte(B2i(Vx > (math.MaxUint8 - Vy)))
}

// Opcode8XY5 Vx -= Vy. VF = 0 when there is a borrow, otherwise 1.
func (chip8 *Chip8) Opcode8XY5(opcode uint16) {
	Vx := chip8.V[chip8.ArgX(opcode)]
	Vy := chip8.V[chip8.ArgY(opcode)]
	chip8.V[chip8.ArgX(opcode)] = Vx - Vy
	chip8.V[0xf] = byte(B2i(Vx >= Vy))
}

// Opcode8XY6 Vx = Vy = Vy >> 1. VF = Least significant bit of Vy before the shift.
func (chip8 *Chip8) Opcode8XY6(opcode uint16) {
	Vy := chip8.V[chip8.ArgY(opcode)]
	chip8.V[chip8.ArgX(opcode)] = chip8.V[chip8.ArgY(opcode)] >> 1
	chip8.V[chip8.ArgY(opcode)] = chip8.V[chip8.ArgY(opcode)] >> 1
	chip8.V[0xf] = Vy & 0x1 // LSB
}

// Opcode8XY7 Vx = Vy - Vx. VF = 0 when there is a borrow, otherwise 1.
func (chip8 *Chip8) Opcode8XY7(opcode uint16) {
	Vx := chip8.V[chip8.ArgX(opcode)]
	Vy := chip8.V[chip8.ArgY(opcode)]
	chip8.V[chip8.ArgX(opcode)] = Vy - Vx
	chip8.V[0xf] = byte(B2i(Vx > Vy))
}

// Opcode8XYE Vx = Vy = Vy << 1. VF = Most significant bit of Vy before the shift.
func (chip8 *Chip8) Opcode8XYE(opcode uint16) {
	Vy := chip8.V[chip8.ArgY(opcode)]
	chip8.V[chip8.ArgX(opcode)] = chip8.V[chip8.ArgY(opcode)] << 1
	chip8.V[chip8.ArgY(opcode)] = chip8.V[chip8.ArgY(opcode)] << 1
	chip8.V[0xf] = Vy >> 7 // MSB
}

// Opcode9XY0 Skip the next instruction if Vx != Vy
func (chip8 *Chip8) Opcode9XY0(opcode uint16) {
	if chip8.V[chip8.ArgX(opcode)] != chip8.V[chip8.ArgY(opcode)] {
		chip8.PC += 2
	}
}

// OpcodeANNN Set I to the address NNN
func (chip8 *Chip8) OpcodeANNN(opcode uint16) {
	chip8.I = chip8.ArgNNN(opcode)
}

// OpcodeBNNN Jumps to the address NNN plus V0
func (chip8 *Chip8) OpcodeBNNN(opcode uint16) {
	chip8.PC = chip8.ArgNNN(opcode) + uint16(chip8.V[0])
}

// OpcodeCXNN Sets VX to the result of a bitwise and operation on a
// random number (Typically: 0 to 255) and NN.
func (chip8 *Chip8) OpcodeCXNN(opcode uint16) {
	randomNum := uint8(rand.Intn(255)) // rand.Intn(max - min) + min
	chip8.V[chip8.ArgX(opcode)] = randomNum & uint8(chip8.ArgNN(opcode))
}

func (chip8 *Chip8) OpcodeDXYN(opcode uint16) {
}

func (chip8 *Chip8) OpcodeEX9E(opcode uint16) {
}

func (chip8 *Chip8) OpcodeEXA1(opcode uint16) {
}

func (chip8 *Chip8) OpcodeFX07(opcode uint16) {
}

func (chip8 *Chip8) OpcodeFX0A(opcode uint16) {
}

func (chip8 *Chip8) OpcodeFX15(opcode uint16) {
}

func (chip8 *Chip8) OpcodeFX18(opcode uint16) {
}

// OpcodeFX1E Adds Vx to I
func (chip8 *Chip8) OpcodeFX1E(opcode uint16) {
	chip8.I = chip8.I + uint16(chip8.V[chip8.ArgX(opcode)])
}

func (chip8 *Chip8) OpcodeFX29(opcode uint16) {
}

// OpcodeFX33 Stores the Binary-Coded Decimal representation of Vx
// with the most significant of the three digits at the address i,
// the middle digit at i + 1, and the least significant at i + 2.
func (chip8 *Chip8) OpcodeFX33(opcode uint16) {
	bcd := chip8.V[chip8.ArgX(opcode)]
	chip8.Mem[chip8.I+0] = bcd / 100
	chip8.Mem[chip8.I+1] = (bcd / 10) % 10
	chip8.Mem[chip8.I+2] = bcd % 10
}

// OpcodeFX55 Stores V0 to Vx (including Vx) in memory starting at
// address I. I is increased by 1 for each value written.
func (chip8 *Chip8) OpcodeFX55(opcode uint16) {
	lastRegister := chip8.ArgX(opcode)
	var j uint16
	for j = 0; j < lastRegister; chip8.I, j = chip8.I+1, j+1 {
		chip8.Mem[chip8.I] = chip8.V[j]
	}
}

// OpcodeFX65 Fills V0 to Vx (including Vx) with values from memory
// starting at address I. I is increased by 1 for each value written.
func (chip8 *Chip8) OpcodeFX65(opcode uint16) {
	lastRegister := chip8.ArgX(opcode)
	var j uint16
	for j = 0; j < lastRegister; chip8.I, j = chip8.I+1, j+1 {
		chip8.V[j] = chip8.Mem[chip8.I]
	}
}
