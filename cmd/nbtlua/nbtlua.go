package main

import (
	"encoding/hex"
	"fmt"

	nlua "github.com/midnightfreddie/nbt-go-lua"
	lua "github.com/yuin/gopher-lua"
)

// luanbt is called to get a Lua environment with nbt manipulation ability// lua vm memory limit; 0 is no limit
const memoryLimitMb = 100

// creates and returns a Lua environment
func luanbt() *lua.LState {

	// The nbt I'm using is in Java Edition format (big endian), so set that:
	nlua.UseJavaEncoding()

	// Create gopher-lua environment
	L := lua.NewState()

	// Set memory limit of lua instance (just a safety measure)
	if memoryLimitMb > 0 {
		L.SetMx(memoryLimitMb)
	}

	// This is a simple nbt tag type 1 (byte) with value of 3 and name "Count"
	//   But you would normally probably be loading nbt from a file
	nbtData := []byte{1, 0, 5, 67, 111, 117, 110, 116, 3}

	// Load a small hard-coded nbt
	err := nlua.Nbt2Lua(nbtData, L)
	if err != nil {
		panic(err)
	}

	fmt.Print("\nA tiny sample nbt tag has been loaded into the global nbt variable. ")
	fmt.Print(`Try print(nbt[1].name) or tagType or value. Try changing the name or value. `)
	fmt.Print("Press control-D to exit. (May be control-Z on Windows.)")
	fmt.Print("A hex dump of the modified nbt will print after exit.\n\n")
	return L
}

// This is called after any file scripts are run and after the interactive script prompt
func afterScripts(L *lua.LState) {
	nlua.UseBedrockEncoding()
	nbtOut, err := nlua.Lua2Nbt(L)
	if err != nil {
		panic(err)
	}
	fmt.Print("\nBedrock format nbt hex dump:\n\n")
	fmt.Println(hex.Dump(nbtOut))

	nlua.UseJavaEncoding()
	nbtOut, err = nlua.Lua2Nbt(L)
	if err != nil {
		panic(err)
	}
	fmt.Print("\nJava format nbt hex dump:\n\n")
	fmt.Println(hex.Dump(nbtOut))
}
