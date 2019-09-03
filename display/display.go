package display

const width = 64
const height = 32

var Fonts = [80]uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type Display struct {
	pixels [2048]uint8
}

func New() *Display {
	d := &Display{}
	return d
}

func (d *Display) Reset() {
	for i := range d.pixels {
		d.pixels[i] = 0
	}
}

func (d *Display) setPixel(x int, y int, set uint8) {
	d.pixels[x+y*width] = set
}

func (d *Display) getPixel(x int, y int) uint8 {
	return d.pixels[x+y*width]
}

func (d *Display) GetPixel(index int) uint8 {
	return d.pixels[index]
}

func (d *Display) Draw(x int, y int, sprite *[]uint8) bool {
	rows := len(*sprite)
	collision := false
	for j := 0; j < rows; j++ {
		row := (*sprite)[j]
		for i := 0; i < 8; i++ {
			// Check if value from memory is 1
			newValue := row >> (7 - uint8(i)) & 0x01
			if newValue == 1 {
				// Check if current pixel state is 1
				xi := (x + i) % width
				yj := (y + j) % height
				oldValue := d.getPixel(xi, yj)
				if oldValue == 1 {
					// Pixel changed from set to unset
					collision = true
				}
				// XOR to get new pixel state
				var xor uint8
				if newValue != oldValue {
					xor = 1
				}
				d.setPixel(xi, yj, xor)
			}
		}
	}
	return collision
}
