# CHIP-8 Interpreter

A CHIP-8 interpreter written in Go and compiled to WebAssembly so that it can be run in the browser. 
Based on Colin Eberhardt's implementation of the project [1], which was originally done in Rust.

![GUI](https://imgur.com/eJbzEt2.png)

## Usage

Build `chip8.wasm` using the Makefile:

    make build
    
Start `server.py` to run locally:

    python3 server.py
    
Finally, head to `localhost:12345`, choose a ROM, and press "Start" to play.


## References

[1] https://blog.scottlogic.com/2017/12/13/chip8-emulator-webassembly-rust.html  
[2] http://devernay.free.fr/hacks/chip8/C8TECH10.HTM  
[3] http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/  
[4] https://blog.gopheracademy.com/advent-2018/go-in-the-browser/
