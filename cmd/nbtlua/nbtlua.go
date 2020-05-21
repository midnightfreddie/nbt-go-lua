package main

import (
	"fmt"

	nlua "github.com/midnightfreddie/nbt-go-lua"
	lua "github.com/yuin/gopher-lua"
)

// creates and returns a Lua environment
func luanbt() *lua.LState {

	// We'll default to Java encoding for this executable
	nlua.UseJavaEncoding()

	// Create gopher-lua environment
	L := nlua.NewState()

	fmt.Print("\nLoad an NBT file with loadfile(path-to-nbt). ")
	fmt.Print(`Try print(nbt[1].name) or tagType or value. Try changing the name or value. `)
	fmt.Print("Press control-D to exit. (May be control-Z on Windows.)")
	fmt.Print("A hex dump of the modified nbt will print after exit.\n\n")
	return L
}
