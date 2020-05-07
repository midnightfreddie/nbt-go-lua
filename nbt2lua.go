package nlua

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"

	lua "github.com/yuin/gopher-lua"
)

// NbtJson is the top-level JSON document; it is exported for reflect, and client code shouldn't use it
type NbtJson struct {
	Name           string             `json:"name"`
	Version        string             `json:"version"`
	Nbt2JsonUrl    string             `json:"nbt2JsonUrl"`
	ConversionTime string             `json:"conversionTime,omitempty"`
	Comment        string             `json:"comment,omitempty"`
	Nbt            []*json.RawMessage `json:"nbt"`
}

// NbtTag represents one NBT tag for each struct; it is exported for reflect, and client code shouldn't use it
type NbtTag struct {
	TagType byte        `json:"tagType"`
	Name    string      `json:"name"`
	Value   interface{} `json:"value,omitempty"`
}

// NbtTagList represents an NBT tag list; it is exported for reflect, and client code shouldn't use it
type NbtTagList struct {
	TagListType byte          `json:"tagListType"`
	List        []interface{} `json:"list"`
}

// NbtLong stores a 64-bit int into two 32-bit values for json portability. ValueMost are the high 32 bits and ValueLeast are the low 32 bits.
//   using uint32s to avoid Go trying to outsmart us on "negative" int32s
type NbtLong struct {
	ValueLeast uint32 `json:"valueLeast"`
	ValueMost  uint32 `json:"valueMost"`
}

// Turns an int64 (nbt long) into a valueLeast/valueMost json pair
func longToIntPair(i int64) NbtLong {
	var nbtLong NbtLong
	nbtLong.ValueLeast = uint32(i & 0xffffffff)
	nbtLong.ValueMost = uint32(i >> 32)
	return nbtLong
}

func intPairToLong(nbtLong NbtLong) int64 {
	var i int64
	i = int64(nbtLong.ValueLeast) | (int64(nbtLong.ValueMost) << 32)
	return i
}

// Nbt2Lua converts uncompressed NBT byte array to gopher-lua LTable
func Nbt2Lua(b []byte, L *lua.LState) error {
	lTable := L.NewTable()
	buf := bytes.NewReader(b)
	for buf.Len() > 0 {
		element, err := getTag(buf, L)
		if err != nil {
			return err
		}
		lTable.Append(element)
	}
	L.SetGlobal("nbt", lTable)
	return nil
}

// getTag broken out form Nbt2Json to allow recursion with reader but public input is []byte
func getTag(r *bytes.Reader, L *lua.LState) (lua.LValue, error) {
	lTable := L.NewTable()
	var tagType byte
	err := binary.Read(r, byteOrder, &tagType)
	if err != nil {
		return lua.LNil, NbtParseError{"Reading TagType", err}
	}
	L.RawSet(lTable, lua.LString("tagType"), lua.LNumber(tagType))
	// do not try to fetch name for TagType 0 which is compound end tag
	if tagType != 0 {
		var err error
		var nameLen int16
		err = binary.Read(r, byteOrder, &nameLen)
		if err != nil {
			return lTable, NbtParseError{"Reading Name length", err}
		}
		name := make([]byte, nameLen)
		err = binary.Read(r, byteOrder, &name)
		if err != nil {
			return lTable, NbtParseError{fmt.Sprintf("Reading Name - is UseJavaEncoding or UseBedrockEncoding set correctly? Name length decoded is %d", nameLen), err}
		}
		L.RawSet(lTable, lua.LString("name"), lua.LString(string(name[:])))
	}
	value, err := getPayload(r, tagType, L)
	if err != nil {
		return lTable, err
	}
	L.RawSet(lTable, lua.LString("value"), value)
	return lTable, err
}

// Gets the tag payload. Had to break this out from the main function to allow tag list recursion
func getPayload(r *bytes.Reader, tagType byte, L *lua.LState) (lua.LValue, error) {
	var err error
	switch tagType {
	case 0:
		// end tag for compound; do nothing further
		return lua.LNil, nil
	case 1:
		var i int8
		err = binary.Read(r, byteOrder, &i)
		if err != nil {
			return nil, NbtParseError{"Reading int8", err}
		}
		return lua.LNumber(i), nil
	case 2:
		var i int16
		err = binary.Read(r, byteOrder, &i)
		if err != nil {
			return nil, NbtParseError{"Reading int16", err}
		}
		return lua.LNumber(i), nil
	case 3:
		var i int32
		err = binary.Read(r, byteOrder, &i)
		if err != nil {
			return nil, NbtParseError{"Reading int32", err}
		}
		return lua.LNumber(i), nil
		/*
			case 4:
				var i int64
				err = binary.Read(r, byteOrder, &i)
				if err != nil {
					return nil, NbtParseError{"Reading int64", err}
				}
				if longAsString {
					output = fmt.Sprintf("%d", i)
				} else {
					output = longToIntPair(i)
				}
		*/
	case 5:
		var f float32
		err = binary.Read(r, byteOrder, &f)
		if err != nil {
			return nil, NbtParseError{"Reading float32", err}
		}
		return lua.LNumber(f), nil
	case 6:
		var f float64
		err = binary.Read(r, byteOrder, &f)
		if err != nil {
			return nil, NbtParseError{"Reading float64", err}
		}
		if math.IsNaN(f) {
			return lua.LNumber(math.NaN()), nil
		} else {
			return lua.LNumber(f), nil
		}
		/*
			case 7:
				var byteArray []int8
				var oneByte int8
				var numRecords int32
				err := binary.Read(r, byteOrder, &numRecords)
				if err != nil {
					return nil, NbtParseError{"Reading byte array tag length", err}
				}
				for i := int32(1); i <= numRecords; i++ {
					err = binary.Read(r, byteOrder, &oneByte)
					if err != nil {
						return nil, NbtParseError{"Reading byte in byte array tag", err}
					}
					byteArray = append(byteArray, oneByte)
				}
				output = byteArray
			case 8:
				var strLen int16
				err := binary.Read(r, byteOrder, &strLen)
				if err != nil {
					return nil, NbtParseError{"Reading string tag length", err}
				}
				utf8String := make([]byte, strLen)
				err = binary.Read(r, byteOrder, &utf8String)
				if err != nil {
					return nil, NbtParseError{"Reading string tag data", err}
				}
				output = string(utf8String[:])
			case 9:
				var tagList NbtTagList
				err = binary.Read(r, byteOrder, &tagList.TagListType)
				if err != nil {
					return nil, NbtParseError{"Reading TagType", err}
				}
				var numRecords int32
				err := binary.Read(r, byteOrder, &numRecords)
				if err != nil {
					return nil, NbtParseError{"Reading list tag length", err}
				}
				for i := int32(1); i <= numRecords; i++ {
					payload, err := getPayload(r, tagList.TagListType)
					if err != nil {
						return nil, NbtParseError{"Reading list tag item", err}
					}
					tagList.List = append(tagList.List, payload)
				}
				output = tagList
							case 10:
								var compound []json.RawMessage
								var tagType byte
								for err = binary.Read(r, byteOrder, &tagType); tagType != 0; err = binary.Read(r, byteOrder, &tagType) {
									if err != nil {
										return nil, NbtParseError{"compound: reading next tag type", err}
									}
									_, err = r.Seek(-1, 1)
									if err != nil {
										return nil, NbtParseError{"seeking back one", err}
									}
									tag, err := getTag(r, L)
									if err != nil {
										return nil, NbtParseError{"compound: reading a child tag", err}
									}
									compound = append(compound, json.RawMessage(string(tag)))
								}
								if compound == nil {
									// Explicitly give empty array else value will be null instead of []
									output = []int{}
								} else {
									output = compound
								}
					case 11:
						var intArray []int32
						var numRecords, oneInt int32
						err := binary.Read(r, byteOrder, &numRecords)
						if err != nil {
							return nil, NbtParseError{"Reading int array tag length", err}
						}
						for i := int32(1); i <= numRecords; i++ {
							err := binary.Read(r, byteOrder, &oneInt)
							if err != nil {
								return nil, NbtParseError{"Reading int in int array tag", err}
							}
							intArray = append(intArray, oneInt)
						}
						output = intArray
					case 12:
						var longArray []NbtLong
						var longStringArray []string
						var numRecords, oneInt int64
						err := binary.Read(r, byteOrder, &numRecords)
						if err != nil {
							return nil, NbtParseError{"Reading long array tag length", err}
						}
						for i := int64(1); i <= numRecords; i++ {
							err := binary.Read(r, byteOrder, &oneInt)
							if err != nil {
								return nil, NbtParseError{"Reading long in long array tag", err}
							}
							longArray = append(longArray, longToIntPair(i))
							longStringArray = append(longStringArray, fmt.Sprintf("%d", i))
						}
						if longAsString {
							output = longStringArray
						} else {
							output = longArray
						}
		*/

	default:
		return nil, NbtParseError{fmt.Sprintf("TagType %d not recognized", tagType), nil}
	}
	// return output, nil
	// return lua.LNil, errors.New("Fell to bottom of getPayload in nbt2lua.go...this shouldn't happen")
}
