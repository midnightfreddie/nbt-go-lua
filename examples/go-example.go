package main

import (
	"fmt"

	nlua "github.com/midnightfreddie/nbt-go-lua"
	lua "github.com/yuin/gopher-lua"
)

func main() {

	// create gopher-lua environment / state / vm
	L := lua.NewState()
	defer L.Close()

	// This is a simple nbt tag type 1 (byte) with value of 3 and name "Count"
	//   But you would normally probably be loading nbt from a file
	nbtData := []byte{1, 0, 5, 67, 111, 117, 110, 116, 3}

	// The nbt is in big endian (Java Edition) format, so let's specify that
	nlua.UseJavaEncoding()

	// Read the raw nbt dat into the Lua environment
	err := nlua.Nbt2Lua(nbtData, L)
	if err != nil {
		panic(err)
	}

	// Now look at it from Lua!
	//    Global `nbt` is a table/array, so the first–in this case only–top-level tag is `nbt[1]`
	err = L.DoString(`print("From Lua:", nbt[1].name)`)
	if err != nil {
		panic(err)
	}

	// We'll stop handling errors for the rest of the example, but you should handle them in real code

	// You can also execute a Lua file script
	_ = L.DoFile("tree-walk.lua")

	// Lets change the value to 64
	_ = L.DoString(`nbt[1].value = 64`)

	// Now let's convert our modified data back to nbt
	modifiedNbtData, _ := nlua.Lua2Nbt(L)
	fmt.Println(nbtData, "- original")
	fmt.Println(modifiedNbtData, "- modified")

	// You'd probably write that to a file, perhaps also gzipping it
}
