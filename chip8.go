package main

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Chip8 struct {
	Mem           [4096]byte
	V             [16]byte // V0 to VF 	(TODO: Can we have aliases?)
	I             [1]uint16
	PC            []uint16
	Stack         [16]uint16
	SP            byte
	DelayRegister byte
	SoundTimer    byte
}

// CTor
func NewChip8() *Chip8 {
	newChip := new(Chip8)
	return newChip
}

// Methods
// Reset - what are the start /reset/reboot values for the computer?
func (chip *Chip8) Reset() {
	chip.PC[0] = 0x200 // or 512. Lets assume pc starts here for now.

	// Set all the registers to zero
	// Do it again here outside the constructor/initialization time,
	// since anyone can reset at any time.
	for i := 0; i < 16; i++ {
		chip.V[i] = 0
	}
	chip.I[0] = 0

	// TODO: Add timer initialization

	// Clear memory
	for i := 0; i < 4096; i++ {
		chip.Mem[i] = 0
	}
}

// Can be run in a separate thread
func (chip *Chip8) MainLoop() {
	for {
		time.Sleep(1000 * time.Millisecond)
	}
}

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
	// First seek to the 512th position, from where we will start reading the file.
	seekedOffset, err := f.Seek(512, 0)
	Check(err)

	// TODO: Why does pong.rom not load?
	// Load all the file data into a buffer
	const DAT_SIZE = 4096 - 512
	const PIXEL_SIZE = 16
	buffer := make([]byte, DAT_SIZE)
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

	fmt.Printf("%d bytes @ %d: %s\n", n, seekedOffset, VM.Mem)

	// https://github.com/veandco/go-sdl2
	err = sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO | sdl.INIT_TIMER)
	Check(err)
	defer sdl.Quit()

	window, err := sdl.CreateWindow("CHIP8", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 64*PIXEL_SIZE, 32*PIXEL_SIZE, sdl.WINDOW_SHOWN)
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
