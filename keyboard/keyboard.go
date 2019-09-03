package keyboard

type Keyboard struct {
	keys [16]bool
}

func New() *Keyboard {
	k := &Keyboard{}
	return k
}

func (k *Keyboard) Reset() {
	for i := range k.keys {
		k.keys[i] = false
	}
}

func (k *Keyboard) PressKey(index uint8, press bool) {
	k.keys[index] = press
}

func (k *Keyboard) IsPressed(index uint8) bool {
	return k.keys[index]
}
