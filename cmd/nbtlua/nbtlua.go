package main

import (
	"fmt"
	"io/ioutil"

	nlua "github.com/midnightfreddie/nbt-go-lua"
	lua "github.com/yuin/gopher-lua"
)

// luanbt is called to get a Lua environment with nbt manipulation ability// lua vm memory limit; 0 is no limit
const memoryLimitMb = 100

func luanbt() *lua.LState {

	////////// temp hard-coded file
	var inData []byte
	var err error

	inData, err = ioutil.ReadFile(`player.dat`)
	if err != nil {
		panic(err)
	}
	////////////////////

	nlua.UseJavaEncoding()
	L := lua.NewState()
	if memoryLimitMb > 0 {
		L.SetMx(memoryLimitMb)
	}
	// err := nlua.Nbt2Lua([]byte{1, 0, 0, 0x7f}, L)
	err = nlua.Nbt2Lua(inData, L)
	if err != nil {
		panic(err)
	}

	return L
}

func afterScripts(L *lua.LState) {
	nbtOut, err := nlua.Lua2Nbt(L)
	if err != nil {
		panic(err)
	}
	fmt.Println(nbtOut)
}
