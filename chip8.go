package chip8

import (
	"syscall/js"

	"github.com/vyun/chip8/vm"
)

const width = 64
const height = 32

type Chip8 struct {
	interpreter vm.VM
	done        chan bool
}

func New() *Chip8 {
	c := Chip8{interpreter: *vm.New(), done: make(chan bool)}
	return &c
}

func (c *Chip8) Start() {
	initializeCanvas()
	// Setup callbacks
	js.Global().Set("initMem", js.FuncOf(c.initMem))
	js.Global().Set("updateDisplay", js.FuncOf(c.updateDisplay))
	js.Global().Set("executeCycle", js.FuncOf(c.executeCycle))
	js.Global().Set("decrementTimers", js.FuncOf(c.decrementTimers))
	js.Global().Set("keyUp", js.FuncOf(c.keyUp))
	js.Global().Set("keyDown", js.FuncOf(c.keyDown))
	<-c.done
}

func getElementByID(id string) js.Value {
	return js.Global().Get("document").Call("getElementById", id)
}

func initializeCanvas() {
	canvas := getElementByID("canvas")
	context := canvas.Call("getContext", "2d")
	context.Set("fillStyle", "black")
	context.Call("fillRect", 0, 0, width, height)
}

func (c *Chip8) initMem(this js.Value, inputs []js.Value) interface{} {
	len := inputs[0].Get("byteLength")
	c.interpreter.Reset()
	for i := 0; i < len.Int(); i++ {
		data := inputs[0].Call("getUint8", i)
		c.interpreter.SetMemory(0x200+i, uint8(data.Int()))

	}
	return nil
}

func (c *Chip8) updateDisplay(this js.Value, inputs []js.Value) interface{} {
	canvas := getElementByID("canvas")
	context := canvas.Call("getContext", "2d")
	image := context.Call("createImageData", width, height)
	imageData := image.Get("data")
	for i := 0; i < width*height; i++ {
		if c.interpreter.GetDisplay().GetPixel(i) == 1 {
			imageData.SetIndex(i*4, 0xff)
			imageData.SetIndex(i*4+1, 0xff)
			imageData.SetIndex(i*4+2, 0xff)
			imageData.SetIndex(i*4+3, 255)
		} else {
			imageData.SetIndex(i*4, 0)
			imageData.SetIndex(i*4+1, 0)
			imageData.SetIndex(i*4+2, 0)
			imageData.SetIndex(i*4+3, 255)
		}
	}

	context.Call("putImageData", image, 0, 0)
	return nil
}

func (c *Chip8) executeCycle(this js.Value, inputs []js.Value) interface{} {
	c.interpreter.Cycle()
	return nil
}

func (c *Chip8) decrementTimers(this js.Value, inputs []js.Value) interface{} {
	c.interpreter.DecrementTimers()
	return nil
}

func (c *Chip8) keyUp(this js.Value, inputs []js.Value) interface{} {
	c.interpreter.GetKeyboard().PressKey(uint8(inputs[0].Int()), false)
	return nil
}

func (c *Chip8) keyDown(this js.Value, inputs []js.Value) interface{} {
	c.interpreter.GetKeyboard().PressKey(uint8(inputs[0].Int()), true)
	return nil
}
