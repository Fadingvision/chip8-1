package display

import "testing"

func TestSet(t *testing.T) {
	display := New()
	display.setPixel(1, 1, 1)

	if display.getPixel(1, 1) == 0 {
		t.Error("TestSet failed")
	}
}

func TestDraw(t *testing.T) {
	display := New()

	sprite := []uint8{0x33, 0xca}

	display.Draw(0, 0, &sprite)
	switch {
	case display.getPixel(0, 0) == 1,
		display.getPixel(1, 0) == 1,
		display.getPixel(2, 0) == 0,
		display.getPixel(3, 0) == 0,
		display.getPixel(4, 0) == 1,
		display.getPixel(5, 0) == 1,
		display.getPixel(6, 0) == 0,
		display.getPixel(7, 0) == 0,
		display.getPixel(0, 1) == 0,
		display.getPixel(1, 1) == 0,
		display.getPixel(2, 1) == 1,
		display.getPixel(3, 1) == 1,
		display.getPixel(4, 1) == 0,
		display.getPixel(5, 1) == 1,
		display.getPixel(6, 1) == 0,
		display.getPixel(7, 1) == 1:

		t.Error("TestDraw failed")

	}
}
func TestCollision(t *testing.T) {
	display := New()
	sprite := []uint8{0x30}

	collision := display.Draw(0, 0, &sprite)
	if collision == true {
		t.Error("TestCollision 1 failed")
	}

	sprite[0] = 0x3
	collision = display.Draw(0, 0, &sprite)
	if collision == true {
		t.Error("TestCollision 2 failed")
	}

	sprite[0] = 0x1
	collision = display.Draw(0, 0, &sprite)
	if collision == false {
		t.Error("TestCollision 3 failed")
	}
}

func TestClear(t *testing.T) {
	display := New()
	display.setPixel(1, 1, 1)
	display.Reset()
	if display.getPixel(1, 1) == 1 {
		t.Error("TestClear failed")
	}
}
