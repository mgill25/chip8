package main

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// Chip8 VM memory and register layout.
type Chip8 struct {
	Mem           [4096]byte
	V             [16]byte // V0 to VF 	(TODO: Can we have aliases?)
	I             uint16
	PC            uint16
	Stack         [16]uint16
	SP            byte
	DelayRegister byte
	SoundTimer    byte

	OpcodeTable []OpcodeTableEntry
}

// OpcodeTableEntry Defines a single opcode, its mask and corresponding handler
type OpcodeTableEntry struct {
	opcode  uint16
	mask    uint16
	Handler func(uint16)
}

// NewChip8 CTor
func NewChip8() *Chip8 {
	newChip := new(Chip8)
	newChip.InitializeOpcodeTable()
	newChip.Reset()
	return newChip
}

// InitializeOpcodeTable Initialize each opcode with its AND mask, and
// the corresponding opcode handler function.
func (chip *Chip8) InitializeOpcodeTable() {
	chip.OpcodeTable = []OpcodeTableEntry{
		{ /* 0x0NNN */ 0x0000, 0xF000, chip.Opcode0NNN},
		{ /* 0x00E0 */ 0x00E0, 0xFFFF, chip.Opcode00E0},
		{ /* 0x00EE */ 0x00EE, 0xFFFF, chip.Opcode00EE},
		{ /* 0x1NNN */ 0x1000, 0xF000, chip.Opcode1NNN},
		{ /* 0x2NNN */ 0x2000, 0xF000, chip.Opcode2NNN},
		{ /* 0x3XNN */ 0x3000, 0xF000, chip.Opcode3XNN},
		{ /* 0x4XNN */ 0x4000, 0xF000, chip.Opcode4XNN},
		{ /* 0x5XY0 */ 0x5000, 0xF00F, chip.Opcode5XY0},
		{ /* 0x6XNN */ 0x6000, 0xF000, chip.Opcode6XNN},
		{ /* 0x7XNN */ 0x7000, 0xF000, chip.Opcode7XNN},
		{ /* 0x8XY0 */ 0x8000, 0xF00F, chip.Opcode8XY0},
		{ /* 0x8XY1 */ 0x8001, 0xF00F, chip.Opcode8XY1},
		{ /* 0x8XY2 */ 0x8002, 0xF00F, chip.Opcode8XY2},
		{ /* 0x8XY3 */ 0x8003, 0xF00F, chip.Opcode8XY3},
		{ /* 0x8XY4 */ 0x8004, 0xF00F, chip.Opcode8XY4},
		{ /* 0x8XY5 */ 0x8005, 0xF00F, chip.Opcode8XY5},
		{ /* 0x8XY6 */ 0x8006, 0xF00F, chip.Opcode8XY6},
		{ /* 0x8XY7 */ 0x8007, 0xF00F, chip.Opcode8XY7},
		{ /* 0x8XYE */ 0x800E, 0xF00F, chip.Opcode8XYE},
		{ /* 0x9XY0 */ 0x9000, 0xF00F, chip.Opcode9XY0},
		{ /* 0xANNN */ 0xA000, 0xF000, chip.OpcodeANNN},
		{ /* 0xBNNN */ 0xB000, 0xF000, chip.OpcodeBNNN},
		{ /* 0xCXNN */ 0xC000, 0xF000, chip.OpcodeCXNN},
		{ /* 0xDXYN */ 0xD000, 0xF000, chip.OpcodeDXYN},
		{ /* 0xEX9E */ 0xE09E, 0xF0FF, chip.OpcodeEX9E},
		{ /* 0xEXA1 */ 0xE001, 0xF0FF, chip.OpcodeEXA1},
		{ /* 0xFX07 */ 0xF007, 0xF0FF, chip.OpcodeFX07},
		{ /* 0xFX0A */ 0xF00A, 0xF0FF, chip.OpcodeFX0A},
		{ /* 0xFX15 */ 0xF015, 0xF0FF, chip.OpcodeFX15},
		{ /* 0xFX18 */ 0xF018, 0xF0FF, chip.OpcodeFX18},
		{ /* 0xFX1E */ 0xF01E, 0xF0FF, chip.OpcodeFX1E},
		{ /* 0xFX29 */ 0xF029, 0xF0FF, chip.OpcodeFX29},
		{ /* 0xFX33 */ 0xF033, 0xF0FF, chip.OpcodeFX33},
		{ /* 0xFX55 */ 0xF055, 0xF0FF, chip.OpcodeFX55},
		{ /* 0xFX65 */ 0xF065, 0xF0FF, chip.OpcodeFX65}}
}

// Reset - what are the start /reset/reboot values for the computer?
func (chip *Chip8) Reset() {
	chip.PC = 0x200 // or 512. Lets assume pc starts here for now.
	// http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#0.1
	//  +---------------+= 0x200 (512) Start of most Chip-8 programs

	// Set all the registers to zero
	// Do it again here outside the constructor/initialization time,
	// since anyone can reset at any time.
	for i := 0; i < 16; i++ {
		chip.V[i] = 0
	}
	chip.I = 0

	// TODO: Add timer initialization

	// Zero the memory
	for i := 0; i < 4096; i++ {
		chip.Mem[i] = 0
	}
}

// MainLoop Can be run in a separate goroutine
func (chip *Chip8) MainLoop() {
	for {
		time.Sleep(1000 * time.Millisecond)
	}
}

// VMThreadFunc launches the main loop
func VMThreadFunc(vm *Chip8) {
	vm.MainLoop()
}

func main() {
	// Chip8 Entry point
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "[Error] No filename passed!\nUsage: chip8 <filename>\n")
		os.Exit(2)
	}
	fname := os.Args[1]
	fmt.Println("ROM file name is:", fname)
	fmt.Println("Initializing the Virtual Machine.")
	VM := NewChip8()
	fmt.Println("Reading the input file.")

	f, err := os.Open(fname)
	Check(err)
	defer f.Close()

	// Load all the file data into a buffer
	const datSize = 4096 - 512
	const pixelSize = 16
	buffer := make([]byte, datSize)
	n, err := f.Read(buffer)
	Check(err)

	if n == 0 {
		fmt.Println("Failed to read the RAM image")
		os.Exit(2)
	}

	// Copy from this temporary buffer into the chip memory
	// Rather than iterating over the entire length of the buffer we have allocated,
	// only iterate over only the number of bytes read.
	for i := 0; i < n; i++ {
		VM.Mem[512+i] = buffer[i]
	}

	// fmt.Println(n, "bytes read", VM.Mem)

	// https://github.com/veandco/go-sdl2
	err = sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO | sdl.INIT_TIMER)
	Check(err)
	defer sdl.Quit()

	window, err := sdl.CreateWindow("CHIP8", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 64*pixelSize, 32*pixelSize, sdl.WINDOW_SHOWN)
	Check(err)
	defer window.Destroy()

	if window == nil {
		sdl.Quit()
		os.Exit(4)
	}

	// TODO: Read up on VSync and accelerated renderer.
	sdlRenderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	Check(err)
	defer sdlRenderer.Destroy()

	surface, err := window.GetSurface()
	Check(err)
	// fmt.Println(surface.W, surface.Pitch)

	// Launch the VM in a separate goroutine (a separate thread)
	go VMThreadFunc(VM)

	pixels := surface.Pixels()
	pixels[256*surface.Pitch+512*4] = 0xff // a blue pixel. Each pixel == 4 bytes, blue, green, red, alpha

	// Start the event handler
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			fmt.Println("t: ", t)
			break
		}
	}
	window.UpdateSurface()

	time.Sleep(3000 * time.Millisecond)
}
