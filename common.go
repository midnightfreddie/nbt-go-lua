package nlua

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io/ioutil"

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
	L.SetGlobal("loadnbt", L.NewFunction(loadNbt))
	L.SetGlobal("savenbt", L.NewFunction(saveNbt))
	L.SetGlobal("use_bedrock_encoding", L.NewFunction(useBedrockEncoding))
	L.SetGlobal("use_java_encoding", L.NewFunction(useJavaEncoding))
}

func loadNbt(L *lua.LState) int {
	var inData []byte
	var err error
	path := L.ToString(1)
	inData, err = ioutil.ReadFile(path)
	if err != nil {
		// TODO: proper error handling inside lua?
		fmt.Println("Error reading file:", err)
		return 0
	}
	// is it gzipped?
	if (inData[0] == 0x1f) && (inData[1] == 0x8b) {
		var uncompressed []byte
		buf := bytes.NewReader(inData)
		zr, err := gzip.NewReader(buf)
		if err != nil {
			// TODO: proper error handling inside lua?
			fmt.Println("Error creating gzip reader on buf:", err)
			return 0
		}
		uncompressed, err = ioutil.ReadAll(zr)
		if err != nil {
			// TODO: proper error handling inside lua?
			fmt.Println("Error un-gzipping file:", err)
			return 0
		}
		inData = uncompressed
	}

	err = Nbt2Lua(inData, L)
	if err != nil {
		// TODO: proper error handling inside lua?
		fmt.Println("Error converting file:", err)
		return 0
	}
	return 0
}

// stub
func saveNbt(L *lua.LState) int {
	path := L.ToString(1)
	compress := L.ToBool(2)
	outData, err := Lua2Nbt(L)
	if err != nil {
		// TODO: proper error handling inside lua?
		fmt.Println("Error converting lua to nbt:", err)
		return 0
	}
	if compress {
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		_, err := zw.Write(outData)
		if err != nil {
			// TODO: proper error handling inside lua?
			fmt.Println("Error creating gzip writer on buf:", err)
			return 0
		}
		err = zw.Close()
		if err != nil {
			// TODO: proper error handling inside lua?
			fmt.Println("Error gzipping file:", err)
			return 0
		}
		outData = buf.Bytes()
	}
	err = ioutil.WriteFile(path, outData, 0644)
	if err != nil {
		// TODO: proper error handling inside lua?
		fmt.Println("Error writing file:", err)
		return 0
	}
	return 0
}

// lua wrapper for UseBedrockEncoding()
func useBedrockEncoding(L *lua.LState) int {
	UseBedrockEncoding()
	return 0
}

// lua wrapper for UseJavaEncoding()
func useJavaEncoding(L *lua.LState) int {
	UseJavaEncoding()
	return 0
}
