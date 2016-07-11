package main

/*
#cgo LDFLAGS: -lluajit-5.1
#include <stdlib.h>
#include <stdio.h>
#include <luajit-2.0/lua.h>
#include <luajit-2.0/lualib.h>
#include <luajit-2.0/lauxlib.h>

int my_func() {
  printf("Hello from the C preamble!\n");
  return 0;
}

*/
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

const src = `
local ffi = require('ffi')
ffi.cdef([[
int my_func();
]])

print("Hello from Lua!")
ffi.C.my_func()
`

func main() {
	C.my_func()
	
	// Initialize state.
	state := C.luaL_newstate()
	if state == nil {
		fmt.Println("Unable to initialize Lua context.")
		os.Exit(1)
	}
	C.luaL_openlibs(state)

	// Compile the script.
	csrc := C.CString(src)
	defer C.free(unsafe.Pointer(csrc))
	if C.luaL_loadstring(state, csrc) != 0 {
		errstring := C.GoString(C.lua_tolstring(state, -1, nil))
		fmt.Printf("Lua error: %v\n", errstring)
		os.Exit(1)
	}

	// Run script.
	if C.lua_pcall(state, 0, 0, 0) != 0 {
		errstring := C.GoString(C.lua_tolstring(state, -1, nil))
		fmt.Printf("Lua execution error: %v\n", errstring)
		os.Exit(1)
	}
}