# First attempt at writing a Virtual Machine

## What is a Virtual Machine

A virtual machine is an *emulation* of a computer system.

Virtual Machines are based on computer architectures, and provide functionality of a physical computer.

### 2 different kinds of VMs

1. System Virtual Machines
	- Emulates the entire computer system, including the Operating System, the boot loader etc.
	- Example would be Virtualbox or VMWare.
2. Process Virtual Machines
	- Designed to execute programs in a platform-independent environment.
	- Abstract platform for an *intermediate language* used as the intermediate representation
	  by a compiler.

	- Examples: Java Virtual Machine (JVM)

## How to create a virtual machine

Good introductory resources:

	[1] https://en.wikibooks.org/wiki/Creating_a_Virtual_Machine/Introduction
	[2] https://www.reddit.com/r/EmuDev/
	[3] http://devernay.free.fr/hacks/chip8/C8TECH10.HTM
	[4] http://mattmik.com/files/chip8/mastering/chip8.html
	[5] https://stackoverflow.com/questions/448673/how-do-emulators-work-and-how-are-they-written

## Beginner VM - [Chip-8](https://en.wikipedia.org/wiki/CHIP-8)

Chip-8 is an intermediate programming langauge which runs on the Chip-8 Virtual Machine.

- Made to allow video games to be easily programmable.
- A large number of classic video games have been ported to Chip-8, such as Pong, Space Invaders, Tetris, Pac-Man.

## Credits
I took the structure of the code from the work done by Gynvael in C++.

Go [here](https://gaming.youtube.com/channel/UCCkVMojdBWS-JtH7TliWkVg) for Gynvael's screencasts.
