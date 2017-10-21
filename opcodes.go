package main

import (
	"fmt"
	"math"
	"math/rand"
)

// Opcode Definitions

func (chip8 *Chip8) Opcode0NNN(opcode uint16) {
	// TODO: This is probably not implemented correctly.
	chip8.Opcode2NNN(opcode)
}

// Opcode00E0 Clears the screen
func (chip8 *Chip8) Opcode00E0(opcode uint16) {
	for i := 0; i < 64*32; i++ {
		chip8.FrameBuffer[i] = 0
	}
	chip8.RedrawScreen()
}

// Opcode00EE Returns from a subroutine
func (chip8 *Chip8) Opcode00EE(opcode uint16) {
	if chip8.CallStack.Len() == 0 {
		fmt.Println("pc: %.4x: stack empty, but return used", chip8.PC)
	}
	chip8.PC = chip8.CallStack.Pop().(uint16) // type assertion used, because interface{}
}

// Opcode1NNN Jumps to address NNN
func (chip8 *Chip8) Opcode1NNN(opcode uint16) {
	chip8.PC = chip8.ArgNNN(opcode)
}

// Opcode2NN Calls subroutine at NNN
func (chip8 *Chip8) Opcode2NNN(opcode uint16) {
	chip8.CallStack.Push(chip8.PC)
	chip8.PC = chip8.ArgNNN(opcode)
	// Should perhaps limit the number of items you can push on stack.
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
	chip8.V[chip8.ArgX(opcode)] = byte(Vx + Vy)
	chip8.V[0xf] = byte(B2i(Vx > (math.MaxUint8 - Vy)))
}

// Opcode8XY5 Vx -= Vy. VF = 0 when there is a borrow, otherwise 1.
func (chip8 *Chip8) Opcode8XY5(opcode uint16) {
	Vx := chip8.V[chip8.ArgX(opcode)]
	Vy := chip8.V[chip8.ArgY(opcode)]
	chip8.V[chip8.ArgX(opcode)] = byte(Vx - Vy)
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
	chip8.PC = (chip8.ArgNNN(opcode) + uint16(chip8.V[0])) & 0xfff // TODO: Why mask?
}

// OpcodeCXNN Sets VX to the result of a bitwise and operation on a
// random number (Typically: 0 to 255) and NN.
func (chip8 *Chip8) OpcodeCXNN(opcode uint16) {
	randomNum := uint8(rand.Intn(255)) // rand.Intn(max - min) + min
	chip8.V[chip8.ArgX(opcode)] = randomNum & uint8(chip8.ArgNN(opcode))
}

/*
OpcodeDXYN Draws a sprite at coordinate (VX, VY) that has
a width of 8 pixels and a height of N pixels. Each row of
8 pixels is read as bit-coded starting from memory location I;
I value doesn’t change after the execution of this instruction.
As described above, VF is set to 1 if any screen pixels are flipped
from set to unset when the sprite is drawn, and to 0 if that doesn’t happen
*/
func (chip8 *Chip8) OpcodeDXYN(opcode uint16) {
	var px uint8
	var j, k, n uint16
	flipped := false

	// Coordinates and height. TODO: Not wrap around?
	x := uint16(chip8.V[chip8.ArgX(opcode)] % 64)
	y := uint16(chip8.V[chip8.ArgY(opcode)] % 32)
	n = chip8.ArgN(opcode)

	// Read rows from memory starting from I.
	for j = 0; j < n; j++ {
		y++                       // Increment y, as the coordinate where we are drawing.s
		px = chip8.Mem[chip8.I+j] // get the memory
		// draw it
		for k = 0; k < 8; k++ {
			// Typical formula for framebuffer: (x + k + y * 64)
			// Assume LSB first
			bit := (px >> k) & 1
			fb := chip8.FrameBuffer[x+k+y*64]
			if (bit != 0) && (fb != 0) {
				flipped = flipped || true
			}
			chip8.FrameBuffer[x+k+y*64] ^= bit
		}
	}

	if flipped {
		chip8.V[0xf] = 1
	} else {
		chip8.V[0xf] = 0
	}

	chip8.RedrawScreen()
}

// OpcodeEx9E Skips the next instruction if the key stored in VX is pressed.
// (Usually the next instruction is a jump to skip a code block)
func (chip8 *Chip8) OpcodeEX9E(opcode uint16) {
	idx := chip8.V[chip8.ArgX(opcode)]
	pressed := false
	if idx < 0x10 {
		pressed = chip8.KeyPressed[idx]
	}
	if pressed {
		chip8.PC += 2
	}
}

// OpcodeEXA1 Skips the next instruction if the key stored in VX isn't pressed.
// (Usually the next instruction is a jump to skip a code block)
func (chip8 *Chip8) OpcodeEXA1(opcode uint16) {
	idx := chip8.V[chip8.ArgX(opcode)]
	pressed := false
	if idx < 0x10 {
		pressed = chip8.KeyPressed[idx]
	}
	if !pressed {
		chip8.PC += 2
	}

}

// FX07 Sets VX to the value of the delay timer.
func (chip8 *Chip8) OpcodeFX07(opcode uint16) {
	chip8.V[chip8.ArgX(opcode)] = chip8.DelayTimer
}

func (chip8 *Chip8) OpcodeFX0A(opcode uint16) {
}

// FX15 Sets the delay timer to VX.
func (chip8 *Chip8) OpcodeFX15(opcode uint16) {
	chip8.DelayTimer = chip8.V[chip8.ArgX(opcode)]
}

// FX18 Sets the sound timer to VX.
func (chip8 *Chip8) OpcodeFX18(opcode uint16) {
	chip8.SoundTimer = chip8.V[chip8.ArgX(opcode)]
}

// OpcodeFX1E Adds Vx to I
func (chip8 *Chip8) OpcodeFX1E(opcode uint16) {
	chip8.I = (chip8.I + uint16(chip8.V[chip8.ArgX(opcode)])) & 0xfff // TODO: Why mask?
}

// FX29 Sets I to the location of the sprite for the character in VX.
// Characters 0-F (in hexadecimal) are represented by a 4x5 font.
func (chip8 *Chip8) OpcodeFX29(opcode uint16) {
	chip8.I = 0
	// TODO: Add Font Support
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
