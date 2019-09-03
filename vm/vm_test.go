package vm

import (
	"testing"
)

// 00EE
func TestRet(t *testing.T) {
	vm := New()
	addr := 0x23
	vm.pc = uint16(addr)

	vm.execute(0x2ABC)
	vm.execute(0x00EE)

	if vm.pc != 0x25 {
		t.Error("TestRet failed, pc")
	} else if vm.sp != 0 {
		t.Error("TestRet failed, sp")
	}
}

// 1nnn
func TestJpAddr(t *testing.T) {
	vm := New()
	vm.execute(0x1A2A)
	if vm.pc != 0x0A2A {
		t.Error("TestJpAddr failed")
	}
}

// 2nnn
func TestCallAddr(t *testing.T) {
	vm := New()
	var addr uint16 = 0x23
	vm.pc = addr
	vm.execute(0x2ABC)
	if vm.pc != 0x0ABC {
		t.Error("TestCallAddr failed - pc")
	} else if vm.sp != 1 {
		t.Error("TestCallAddr failed - sp")
	} else if vm.stack[0] != addr+2 {
		t.Error("TestCallAddr failed - stack")
	}
}

// 3xkk
func TestSeVxByte(t *testing.T) {
	vm := New()
	vm.v[1] = 0xFE

	vm.execute(0x31FE)
	if vm.pc != 4 {
		t.Error("TestSeVxByte failed - pc did not skip")
	}
	vm.execute(0x31FA)
	if vm.pc != 6 {
		t.Error("TestSeVxByte failed - pc skipped")
	}
}

// 4xkk
func TestSneVxByte(t *testing.T) {
	vm := New()
	vm.v[1] = 0xFE

	vm.execute(0x41FE)
	if vm.pc != 2 {
		t.Error("TestSneVxByte failed - pc skipped")
	}
	vm.execute(0x41FA)
	if vm.pc != 6 {
		t.Error("TestSneVxByte failed - pc did not skip")
	}
}

// 5xy0
func TestSeVxVy(t *testing.T) {
	vm := New()
	vm.v[0] = 1
	vm.v[1] = 2
	vm.v[2] = 2

	vm.execute(0x5120)
	if vm.pc != 4 {
		t.Error("TestSeVxVy failed - pc did not skip")
	}
	vm.execute(0x5020)
	if vm.pc != 6 {
		t.Error("TestSeVxVy failed - pc skipped")
	}
}

// 6xy0
func TestLdVxByte(t *testing.T) {
	vm := New()
	vm.v[0] = 1

	vm.execute(0x6008)
	if vm.v[0] != 8 {
		t.Error("TestLdVxByte failed")
	}
}

// 7xkk
func TestAddVxByte(t *testing.T) {
	vm := New()
	vm.v[0] = 1

	vm.execute(0x7008)
	if vm.v[0] != 9 {
		t.Error("TestAddVxByte failed")
	}
}

// 8xy0
func TestLdVxVy(t *testing.T) {
	vm := New()
	vm.v[0] = 1
	vm.v[1] = 0

	vm.execute(0x8010)
	if vm.v[0] != 0 {
		t.Error("TestLdVxVy failed")
	}
}

// 8xy1
func TestOrVxVy(t *testing.T) {
	vm := New()
	vm.v[0] = 0xFF
	vm.v[1] = 0xAA

	vm.execute(0x8011)
	if vm.v[0] != 0xFF {
		t.Error("TestOrVxVy failed")
	}
}

// 8xy2
func TestAndVxVy(t *testing.T) {
	vm := New()
	vm.v[0] = 0xFF
	vm.v[1] = 0xAA

	vm.execute(0x8012)
	if vm.v[0] != 0xAA {
		t.Error("TestAndVxVy failed")
	}
}

// 8xy3
func TestXorVxVy(t *testing.T) {
	vm := New()
	vm.v[0] = 0xFF
	vm.v[1] = 0xAA

	vm.execute(0x8013)
	if vm.v[0] != 0x55 {
		t.Error("TestXorVxVy failed")
	}
}

// 8xy4
func TestAddVxVy(t *testing.T) {
	vm := New()
	vm.v[0] = 10
	vm.v[1] = 100
	vm.v[2] = 250

	vm.execute(0x8014)
	if vm.v[0] != 110 {
		t.Error("TestAddVxVy failed")
	} else if vm.v[0xF] != 0 {
		t.Error("TestAddVxVy failed - incorrect overflow")
	}
	vm.execute(0x8024)
	if vm.v[0] != 0x68 {
		t.Error("TestAddVxVy failed")
	} else if vm.v[0xF] != 1 {
		t.Error("TestAddVxVy failed - did not detect overflow")
	}
}

// 9xy0
func TestSneVxVy(t *testing.T) {
	vm := New()
	vm.v[0] = 1
	vm.v[1] = 2
	vm.v[2] = 2

	vm.execute(0x9120)
	if vm.pc != 2 {
		t.Error("TestSneVxVy failed - pc skipped")
	}
	vm.execute(0x9020)
	if vm.pc != 6 {
		t.Error("TestSneVxVy failed - pc did not skip")
	}
}

// Annn
func TestLdIAddr(t *testing.T) {
	vm := New()

	vm.execute(0xA111)
	if vm.i != 0x111 {
		t.Error("TestLdIAddr failed")
	}
}

// Bnnn
func TestJpV0Addr(t *testing.T) {
	vm := New()
	vm.v[0] = 2
	vm.execute(0xBA2A)
	if vm.pc != 0x0A2C {
		t.Error("TestJpV0Addr failed")
	}
}

// Fx1E
func TestAddIVx(t *testing.T) {
	vm := New()
	vm.i = 2
	vm.v[0] = 6
	vm.execute(0xF01E)
	if vm.i != 8 {
		t.Error("TestAddIVx failed")
	}
}

// Fx33
func TestLdBVx(t *testing.T) {
	vm := New()
	vm.i = 0x300
	vm.v[2] = 234

	// load v2's BCD representation in memory starting at i
	vm.execute(0xF233)
	if vm.memory[vm.i] != 2 {
		t.Error("TestLdMemIVx failed, i")
	} else if vm.memory[vm.i+1] != 3 {
		t.Error("TestLdMemIVx failed, i+1")
	} else if vm.memory[vm.i+2] != 4 {
		t.Error("TestLdMemIVx failed, i+2")
	}
}

// Fx55
func TestLdMemIVx(t *testing.T) {
	vm := New()
	vm.v[0] = 1
	vm.v[1] = 2
	vm.v[2] = 3
	vm.v[3] = 4
	vm.i = 0x300

	// only load v0-v2 into memory starting at i
	vm.execute(0xF255)
	if vm.memory[vm.i] != 1 {
		t.Error("TestLdMemIVx failed, i")
	} else if vm.memory[vm.i+1] != 2 {
		t.Error("TestLdMemIVx failed, i+1")
	} else if vm.memory[vm.i+2] != 3 {
		t.Error("TestLdMemIVx failed, i+2")
	} else if vm.memory[vm.i+3] != 0 {
		t.Error("TestLdMemIVx failed, i+3")
	}
}

// Fx65
func TestLdVxMemI(t *testing.T) {
	vm := New()
	vm.memory[vm.i] = 1
	vm.memory[vm.i+1] = 2
	vm.memory[vm.i+2] = 3
	vm.memory[vm.i+3] = 4

	vm.execute(0xF265)

	// only load into v0-v2
	if vm.v[0] != 1 {
		t.Error("TestLdMemIVx failed, v[0]")
	} else if vm.v[1] != 2 {
		t.Error("TestLdMemIVx failed, v[1]")
	} else if vm.v[2] != 3 {
		t.Error("TestLdMemIVx failed, v[2]")
	} else if vm.v[3] != 0 {
		t.Error("TestLdMemIVx failed, v[3]")
	}

}
