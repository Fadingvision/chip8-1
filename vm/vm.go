package vm

import (
	"fmt"
	"math/rand"

	"github.com/vyun/chip8/display"
	"github.com/vyun/chip8/keyboard"
)

type VM struct {
	memory     [4096]uint8 // 4096 bit address space
	v          [16]uint8   // 16 8-bit general purpose registers, V0 to VF
	i          uint16      // 16-bit I register
	stack      [16]uint16  // 16 16-bit value stack
	sp         uint8       // Stack pointer
	pc         uint16      // Program counter
	delayTimer uint8
	soundTimer uint8
	display    display.Display
	keyboard   keyboard.Keyboard
}

func New() *VM {
	vm := &VM{}
	return vm
}

func (vm *VM) Reset() {
	for i := range vm.memory {
		if i < 80 {
			vm.memory[i] = display.Fonts[i]
		} else {
			vm.memory[i] = 0
		}
	}
	for i := range vm.v {
		vm.v[i] = 0
	}
	for i := range vm.stack {
		vm.stack[i] = 0
	}
	vm.i = 0
	vm.sp = 0
	vm.pc = 0x200 // Program/data space starts at 512 byte offset
	vm.delayTimer = 0
	vm.soundTimer = 0
	vm.display = *display.New()
	vm.keyboard = *keyboard.New()
	vm.display.Reset()
	vm.keyboard.Reset()
}

func (vm *VM) fetch() uint16 {
	opcode := uint16(vm.memory[vm.pc])<<8 | uint16(vm.memory[(vm.pc)+1])
	return opcode
}

func (vm *VM) execute(opcode uint16) {
	x := uint8((opcode & 0x0F00) >> 8)
	y := uint8((opcode & 0x00F0) >> 4)
	nnn := opcode & 0x0FFF
	kk := uint8(opcode & 0x00FF)
	n := uint8(opcode & 0x000F)

	vm.pc += 2

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x000F {
		case 0x0000:
			// 00E0 - CLS
			vm.display.Reset()
		case 0x000E:
			// 00EE - RET
			vm.sp--
			vm.pc = vm.stack[vm.sp]
		default:
			fmt.Printf("Unknown opcode: %X\n", opcode)
		}
	case 0x1000:
		// 1nnn - JP addr
		vm.pc = nnn
	case 0x2000:
		// 2nnn - CALL addr
		vm.stack[vm.sp] = vm.pc
		vm.pc = nnn
		vm.sp++
	case 0x3000:
		// 3xkk - SE Vx, byte
		if vm.v[x] == kk {
			vm.pc += 2
		}
	case 0x4000:
		// 4xkk - SNE Vx, byte
		if vm.v[x] != kk {
			vm.pc += 2
		}
	case 0x5000:
		// 5xy0 - SE Vx, Vy
		if vm.v[x] == vm.v[y] {
			vm.pc += 2
		}
	case 0x6000:
		// 6xkk - LD Vx, byte
		vm.v[x] = kk
	case 0x7000:
		// 7xkk - ADD Vx, byte
		vm.v[x] += kk
	case 0x8000:
		// 8xy0 - LD Vx, Vy
		switch opcode & 0x000F {
		case 0x0000:
			// 8xy0 - LD Vx, Vy
			vm.v[x] = vm.v[y]
		case 0x0001:
			// 8xy1 - OR Vx, Vy
			vm.v[x] |= vm.v[y]
		case 0x0002:
			// 8xy2 - AND Vx, Vy
			vm.v[x] &= vm.v[y]
		case 0x0003:
			// 8xy3 - XOR Vx, Vy
			vm.v[x] ^= vm.v[y]
		case 0x0004:
			// 8xy4 - ADD Vx, Vy
			vm.v[0xF] = 0
			if vm.v[x] > (0xFF - vm.v[y]) {
				vm.v[0xF] = 1
			}
			vm.v[x] += vm.v[y]

		case 0x0005:
			// 8xy5 - SUB Vx, Vy
			vm.v[0xF] = 0
			if vm.v[x] > vm.v[y] {
				vm.v[0xF] = 1
			}
			vm.v[x] -= vm.v[y]
		case 0x0006:
			// 8xy6 - SHR Vx {, Vy}
			vm.v[0xF] = vm.v[x] & 0x01
			vm.v[x] >>= 1
		case 0x0007:
			// 8xy7 - SUBN Vx, Vy
			vm.v[0xF] = 0
			if vm.v[y] > vm.v[x] {
				vm.v[0xF] = 1
			}
			vm.v[x] = vm.v[y] - vm.v[x]
		case 0x000E:
			// 8xyE - SHL Vx {, Vy}
			vm.v[0xF] = vm.v[x] & 0x80
			vm.v[x] <<= 1
		}
	case 0x9000:
		// 9xy0 - SNE Vx, Vy
		if vm.v[x] != vm.v[y] {
			vm.pc += 2
		}
	case 0xA000:
		// Annn - LD I, addr
		vm.i = nnn
	case 0xB000:
		// Bnnn - JP V0, addr
		vm.pc = nnn + uint16(vm.v[0])
	case 0xC000:
		// Cxkk - RND Vx, byte
		vm.v[x] = uint8(rand.Intn(256)) & kk
	case 0xD000:
		// Dxyn - DRW Vx, Vy, nibble
		sprite := vm.memory[vm.i : vm.i+uint16(n)]
		vm.v[0xF] = 0
		collision := vm.display.Draw(int(vm.v[x]), int(vm.v[y]), &sprite)
		if collision {
			vm.v[0xF] = 1
		}

	case 0xE000:
		switch opcode & 0x00FF {
		case 0x009E:
			// Ex9E - SKP Vx
			if vm.keyboard.IsPressed(vm.v[x]) == true {
				vm.pc += 2
			}
		case 0x00A1:
			// ExA1 - SKNP Vx
			if vm.keyboard.IsPressed(vm.v[x]) == false {
				vm.pc += 2
			}
		default:
			fmt.Printf("Unknown opcode: %X\n", opcode)
		}
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x0007:
			// Fx07 - LD Vx, DT
			vm.v[x] = vm.delayTimer
		case 0x000A:
			// Fx0A - LD Vx, K
			vm.pc -= 2
			for i := uint8(0); i < 16; i++ {
				if vm.keyboard.IsPressed(i) == true {
					vm.v[x] = i
					vm.pc += 2
					break
				}
			}
		case 0x0015:
			// Fx15 - LD DT, Vx
			vm.delayTimer = vm.v[x]
		case 0x0018:
			// Fx18 - LD ST, Vx
			vm.soundTimer = vm.v[x]
		case 0x001E:
			// Fx1E - ADD I, Vx
			vm.i += uint16(vm.v[x])
		case 0x0029:
			// Fx29 - LD F, Vx
			vm.i = uint16(vm.v[x]) * 5
		case 0x0033:
			// Fx33 - LD B, Vx
			vm.memory[vm.i] = vm.v[x] / 100
			vm.memory[vm.i+1] = (vm.v[x] / 10) % 10
			vm.memory[vm.i+2] = (vm.v[x] % 100) % 10
		case 0x0055:
			// Fx55 - LD [I], Vx
			for i := uint8(0); i <= x; i++ {
				vm.memory[vm.i+uint16(i)] = vm.v[i]
			}
		case 0x0065:
			// Fx65 - LD Vx, [I]
			for i := uint8(0); i <= x; i++ {
				vm.v[i] = vm.memory[vm.i+uint16(i)]
			}
		default:
			fmt.Printf("Unknown opcode: %X\n", opcode)
		}
	default:
		fmt.Printf("Unknown opcode: %X\n", opcode)

	}
}

func (vm *VM) Cycle() {
	opcode := vm.fetch()
	vm.execute(opcode)
}

func (vm *VM) DecrementTimers() {
	if vm.delayTimer > 0 {
		vm.delayTimer--
	}
	if vm.soundTimer > 0 {
		vm.soundTimer--
	}
}

func (vm *VM) GetDisplay() *display.Display {
	return &vm.display
}

func (vm *VM) GetKeyboard() *keyboard.Keyboard {
	return &vm.keyboard
}

func (vm *VM) SetMemory(index int, data uint8) {
	vm.memory[index] = data
}
