# Notes from the Gynvael stream on Chip-8 VM

- Usually in VMs or even older machines, you have a frame buffer, and you address a specific point like a pixel and turn it on/off, but for Sprites, you have a register which you have to set, and another set of registers which basically you say you want to place the sprite and say etc etc (bit complicated, look into this later!)

- Sprites are 8 pixels wide and 1 to 15 pixels in height.

- XOR operation to flip the color of the sprite.


35 Opcodes (waaay easier than x86). Also 4.5x more opcodes than Brainfuck!

- Each opcode is 2 bytes long. 
- Big-Endian (x86 is little-endian)


Instructions are stored in memory, and are 2 bytes long. 1st byte of each instruction should be located
at an even address.

Opcodes have the following symbols:
	NNN: Address
	NN: 8-bit Constant
	N: 4-bit Constant
	X and Y: 4-bit register identifier
	PC: Program Counter
	I: 16-bit register