package main

import (
	nlua "github.com/midnightfreddie/nbt-go-lua"
	lua "github.com/yuin/gopher-lua"
)

// luanbt is called to get a Lua environment with nbt manipulation ability// lua vm memory limit; 0 is no limit
const memoryLimitMb = 100

func luanbt() *lua.LState {
	L := lua.NewState()
	err := nlua.Nbt2Lua(L, []byte{1})
	if err != nil {
		panic(err)
	}
	if memoryLimitMb > 0 {
		L.SetMx(memoryLimitMb)
	}
	// lTable, err := nlua.Nbt2LTable([]byte{1})
	return L
}
