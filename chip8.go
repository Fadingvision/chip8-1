package chip8

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"syscall/js"

	"github.com/vyun/chip8/vm"
)

const width = 64
const height = 32

var keyMap = map[int]int{
	49: 0x1,
	50: 0x2,
	51: 0x3,
	52: 0xc,
	81: 0x4,
	87: 0x5,
	69: 0x6,
	82: 0xd,
	65: 0x7,
	83: 0x8,
	68: 0x9,
	70: 0xe,
	90: 0xa,
	88: 0x0,
	67: 0xb,
	86: 0xf,
}

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
	c.initializeEvents()

	// c.updateDisplay();
	// Setup callbacks
	// js.Global().Set("c.", js.FuncOf(c.c.))
	js.Global().Set("updateDisplay", js.FuncOf(c.updateDisplay))
	js.Global().Set("executeCycle", js.FuncOf(c.executeCycle))
	js.Global().Set("decrementTimers", js.FuncOf(c.decrementTimers))
	// js.Global().Set("keyUp", js.FuncOf(c.keyUp))
	// js.Global().Set("keyDown", js.FuncOf(c.keyDown))
	c.initMem("INVADERS")
	fmt.Println("main done")
	<-c.done
}

func getElementByID(id string) js.Value {
	return js.Global().Get("document").Call("getElementById", id)
}

func (c *Chip8) initializeEvents() {
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		keyCode := args[0].Get("keyCode").Int()
		if key, ok := keyMap[keyCode]; ok {
			c.keyDown(key)
		}
		return nil
	})
	cb2 := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		keyCode := args[0].Get("keyCode").Int()
		if key, ok := keyMap[keyCode]; ok {
			c.keyUp(key)
		}
		return nil
	})
	cb3 := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		romName := args[0].Get("target").Get("value").String()
		fmt.Println(romName)
		c.initMem(romName)
		return nil
	})
	js.Global().Get("document").Call("addEventListener", "keydown", cb)
	js.Global().Get("document").Call("addEventListener", "keyup", cb2)
	getElementByID("roms").Call("addEventListener", "change", cb3)
}

func initializeCanvas() {
	canvas := getElementByID("canvas")
	context := canvas.Call("getContext", "2d")
	context.Set("fillStyle", "black")
	context.Call("fillRect", 0, 0, width, height)
}

// func (c *Chip8) initMem(this js.Value, inputs []js.Value) interface{} {
// 	len := inputs[0].Get("byteLength")
// 	c.interpreter.Reset()
// 	for i := 0; i < len.Int(); i++ {
// 		data := inputs[0].Call("getUint8", i)
// 		c.interpreter.SetMemory(0x200+i, uint8(data.Int()))

// 	}
// 	return nil
// }
func (c *Chip8) initMem(rom string) interface{} {
	err, buf := loadRom(rom)
	fmt.Println(err)
	c.interpreter.Reset()
	for i := 0; i < len(buf); i++ {
		c.interpreter.SetMemory(0x200+i, buf[i])
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

func (c *Chip8) keyUp(key int) interface{} {
	c.interpreter.GetKeyboard().PressKey(uint8(key), false)
	return nil
}

func (c *Chip8) keyDown(key int) interface{} {
	c.interpreter.GetKeyboard().PressKey(uint8(key), true)
	return nil
}

// MyGoFunc fetches an external resource by making a HTTP request from Go
// The JavaScript method accepts one argument, which is the URL to request
func loadRom(rom string) (error, []byte) {
	resp, err := http.Get("/roms/" + rom)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("response status code: %d", resp.StatusCode)
		return err, nil
	}

	romContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	fmt.Println("fetch done")
	return nil, romContent
}
