package nlua

import (
	"encoding/binary"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// NewState is called to get a Lua environment with nbt manipulation ability// lua vm memory limit; 0 is no limit
const memoryLimitMb = 100

// Used by all converters; change with UseJavaEncoding() or UseBedrockEncoding()
var byteOrder = binary.ByteOrder(binary.LittleEndian)

// UseJavaEncoding sets the module to decode/encode from/to big endian NBT for Minecraft Java Edition
func UseJavaEncoding() {
	byteOrder = binary.BigEndian
}

// UseBedrockEncoding sets the module to decode/encode from/to little endian NBT for Minecraft Bedrock Edition
func UseBedrockEncoding() {
	byteOrder = binary.LittleEndian
}

// Turns an int64 (nbt long) into a least-/most- significant 32 bits pair
func longToIntPair(i int64) (least uint32, most uint32) {
	least = uint32(i & 0xffffffff)
	most = uint32(i >> 32)
	return
}

func intPairToLong(least uint32, most uint32) int64 {
	var i int64
	i = int64(least) | (int64(most) << 32)
	return i
}

// NbtParseError is when the nbt data does not match an expected pattern. Pass it message string and downstream error
type NbtParseError struct {
	s string
	e error
}

func (e NbtParseError) Error() string {
	var s string
	if e.e != nil {
		s = fmt.Sprintf(": %s", e.e.Error())
	}
	return fmt.Sprintf("Error parsing NBT: %s%s", e.s, s)
}

// LuaNbtError is when the lua nbt table data does not match an expected pattern. Pass it message string and downstream error
type LuaNbtError struct {
	s string
	e error
}

func (e LuaNbtError) Error() string {
	var s string
	if e.e != nil {
		s = fmt.Sprintf(": %s", e.e.Error())
	}
	return fmt.Sprintf("Error lua nbt to native nbt: %s%s", e.s, s)
}

func NewState() *lua.LState {
	L := lua.NewState()
	// Set memory limit of lua instance (just a safety measure)
	if memoryLimitMb > 0 {
		L.SetMx(memoryLimitMb)
	}
	Nlua(L)
	return L
}

// Nlua injects load and save functions into a lua environment
func Nlua(L *lua.LState) {

}
