package main

/*
#cgo LDFLAGS: -lluajit-5.1
#include &lt;stdlib.h&gt;
#include &lt;stdio.h&gt;
#include &lt;luajit-2.0/lua.h&gt;
#include &lt;luajit-2.0/lualib.h&gt;
#include &lt;luajit-2.0/lauxlib.h&gt;

int my_func() {
  printf(&#34;Hello from the C preamble!\n&#34;);
  return 0;
}

*/
import &#34;C&#34;

import (
	&#34;fmt&#34;
	&#34;os&#34;
	&#34;unsafe&#34;
)

const src = `
local ffi = require(&#39;ffi&#39;)
ffi.cdef([[
int my_func();
]])

print(&#34;Hello from Lua!&#34;)
ffi.C.my_func()
`

func main() {
	C.my_func()
	
	// Initialize state.
	state := C.luaL_newstate()
	if state == nil {
		fmt.Println(&#34;Unable to initialize Lua context.&#34;)
		os.Exit(1)
	}
	C.luaL_openlibs(state)

	// Compile the script.
	csrc := C.CString(src)
	defer C.free(unsafe.Pointer(csrc))
	if C.luaL_loadstring(state, csrc) != 0 {
		errstring := C.GoString(C.lua_tolstring(state, -1, nil))
		fmt.Printf(&#34;Lua error: %v\n&#34;, errstring)
		os.Exit(1)
	}

	// Run script.
	if C.lua_pcall(state, 0, 0, 0) != 0 {
		errstring := C.GoString(C.lua_tolstring(state, -1, nil))
		fmt.Printf(&#34;Lua execution error: %v\n&#34;, errstring)
		os.Exit(1)
	}
}