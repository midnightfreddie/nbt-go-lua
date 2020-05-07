package nlua

import (
	"encoding/binary"
	"fmt"
)

// Version is the json document's nbt2JsonVersion:
const Version = "0.0.0"

// Nbt2JsonUrl is inserted in the json document as nbt2JsonUrl
const Nbt2JsonUrl = "https://github.com/midnightfreddie/nbt-go-lua"

// Name is the json document's name:
var Name = "Named Binary Tag to Lua"

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

// JsonParseError is when the json data does not match an expected pattern. Pass it message string and downstream error
type JsonParseError struct {
	s string
	e error
}

func (e JsonParseError) Error() string {
	var s string
	if e.e != nil {
		s = fmt.Sprintf(": %s", e.e.Error())
	}
	return fmt.Sprintf("Error parsing json2nbt: %s%s", e.s, s)
}
