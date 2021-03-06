package nlua

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

// Lua2Nbt converts lua's global nbt table variable to uncompressed NBT byte array
//   Note: arrays/lists will iterate all keys in the table, even non-numeric even though Nbt2Lua will not make those
//   Note: A nil lua nbt will return an error, but an nbt empty table will return an empty byte array
func Lua2Nbt(L *lua.LState) ([]byte, error) {
	nbtOut := new(bytes.Buffer)
	nbtArray := L.GetGlobal("nbt")
	var forEachErr error
	if nbtLuaTable, ok := nbtArray.(*lua.LTable); ok {
		nbtLuaTable.ForEach(func(_ lua.LValue, v lua.LValue) {
			if nbtLuaTag, ok := v.(*lua.LTable); ok {
				err := writeTag(nbtOut, nbtLuaTag, L)
				if err != nil {
					if forEachErr == nil {
						forEachErr = err
					}
				}
			}
		})
	} else {
		return nil, LuaNbtError{fmt.Sprintf("Global nbt type, expected %T, got %T", lua.LTable{}, nbtArray), nil}
	}
	return nbtOut.Bytes(), forEachErr
}

// called by Lua2Nbt for each LTable representing an nbt tag; also called from writePayload for compound tags
func writeTag(w io.Writer, nbtLuaTag *lua.LTable, L *lua.LState) error {
	var err error
	var lValue lua.LValue
	lValue = nbtLuaTag.RawGet(lua.LString("tagType"))
	if tagType, ok := lValue.(lua.LNumber); ok {
		if tagType == 0 {
			// not expecting a 0 tag, but if it occurs just ignore it
			return nil
		}
		err = binary.Write(w, byteOrder, byte(tagType))
		if err != nil {
			return LuaNbtError{"Error writing tagType" + string(byte(tagType)), err}
		}
		lValue = nbtLuaTag.RawGet(lua.LString("name"))
		if name, ok := lValue.(lua.LString); ok {
			err = binary.Write(w, byteOrder, int16(len(name)))
			if err != nil {
				return LuaNbtError{"Error writing name length", err}
			}
			err = binary.Write(w, byteOrder, []byte(name))
			if err != nil {
				return LuaNbtError{"Error converting name", err}
			}
		} else {
			return LuaNbtError{fmt.Sprintf("name field '%v' not a string", lValue), err}
		}
		lValue = nbtLuaTag.RawGet(lua.LString("value"))
		err = writePayload(w, lValue, tagType, L)
		if err != nil {
			return err
		}

	} else {
		return LuaNbtError{fmt.Sprintf("tagType '%v' is not an integer", lValue), err}
	}
	return err
}

// called by writeTag to return the value (v) field for the given tag type
func writePayload(w io.Writer, v lua.LValue, tagType lua.LNumber, L *lua.LState) error {
	var err error

	switch tagType {
	case 1:
		if i, ok := v.(lua.LNumber); ok {
			if i < math.MinInt8 || i > math.MaxInt8 {
				return LuaNbtError{fmt.Sprintf("%v is out of range for tag 1 - Byte", i), nil}
			}
			err = binary.Write(w, byteOrder, int8(i))
			if err != nil {
				return LuaNbtError{"Error writing byte payload", err}
			}
		} else {
			return LuaNbtError{fmt.Sprintf("Tag 1 Byte value field '%v' not an integer", v), err}
		}
	case 2:
		if i, ok := v.(lua.LNumber); ok {
			if i < math.MinInt16 || i > math.MaxInt16 {
				return LuaNbtError{fmt.Sprintf("%v is out of range for tag 2 - Short", i), nil}
			}
			err = binary.Write(w, byteOrder, int16(i))
			if err != nil {
				return LuaNbtError{"Error writing short payload", err}
			}
		} else {
			return LuaNbtError{fmt.Sprintf("Tag 2 Short value field '%v' not an integer", v), err}
		}
	case 3:
		if i, ok := v.(lua.LNumber); ok {
			if i < math.MinInt32 || i > math.MaxInt32 {
				return LuaNbtError{fmt.Sprintf("%v is out of range for tag 3 - Int", i), nil}
			}
			err = binary.Write(w, byteOrder, int32(i))
			if err != nil {
				return LuaNbtError{"Error writing int32 payload", err}
			}
		} else {
			return LuaNbtError{fmt.Sprintf("Tag 3 Int value field '%v' not an integer", v), err}
		}
	case 4:
		if lValue, ok := v.(*lua.LTable); ok {
			var vl, vm lua.LNumber
			var lv lua.LValue
			lv = L.RawGet(lValue, lua.LString("least"))
			if vl, ok = lv.(lua.LNumber); !ok {
				return LuaNbtError{fmt.Sprintf("Error reading valueLeast of '%v'", lv), nil}
			}
			lv = L.RawGet(lValue, lua.LString("most"))
			if vm, ok = lv.(lua.LNumber); !ok {
				return LuaNbtError{fmt.Sprintf("Error reading valueMost of '%v'", lv), nil}
			}
			err = binary.Write(w, byteOrder, int64(intPairToLong(uint32(vl), uint32(vm))))
			if err != nil {
				return LuaNbtError{"Error writing int64 (from uint32 pair) payload:", err}
			}
		} else {
			return LuaNbtError{fmt.Sprintf("Tag 4 Long value field '%v' not an object", v), err}
		}
	case 5:
		if f, ok := v.(lua.LNumber); ok {
			if f != 0 && (math.Abs(float64(f)) < math.SmallestNonzeroFloat32 || math.Abs(float64(f)) > math.MaxFloat32) {
				return LuaNbtError{fmt.Sprintf("%g is out of range for tag 5 - Float", f), nil}
			}
			err = binary.Write(w, byteOrder, float32(f))
			if err != nil {
				return LuaNbtError{"Error writing float32 payload", err}
			}
		} else {
			// will write NaN which is needed for true NaN
			err = binary.Write(w, byteOrder, float32(math.NaN()))
			if err != nil {
				return LuaNbtError{"Error writing float64 payload", err}
			}
		}
	case 6:
		if f, ok := v.(lua.LNumber); ok {
			err = binary.Write(w, byteOrder, f)
			if err != nil {
				return LuaNbtError{"Error writing float64 payload", err}
			}
		} else {
			// will write NaN which is needed for true NaN
			err = binary.Write(w, byteOrder, math.NaN())
			if err != nil {
				return LuaNbtError{"Error writing float64 payload", err}
			}
		}
	case 7:
		if values, ok := v.(*lua.LTable); ok {
			err = binary.Write(w, byteOrder, int32(values.Len()))
			if err != nil {
				return LuaNbtError{"Error writing byte array length", err}
			}
			var forEachErr error
			values.ForEach(func(_ lua.LValue, n lua.LValue) {
				if i, ok := n.(lua.LNumber); ok {
					if i < math.MinInt8 || i > math.MaxInt8 && forEachErr == nil {
						forEachErr = LuaNbtError{fmt.Sprintf("%v is out of range for Byte in tag 7 - Byte Array", i), nil}
					}
					err = binary.Write(w, byteOrder, int8(i))
					if err != nil && forEachErr == nil {
						forEachErr = LuaNbtError{"Error writing element of byte array", err}
					}
				} else if forEachErr == nil {
					forEachErr = LuaNbtError{fmt.Sprintf("Tag 7 Byte Array element value field '%v' not an integer", v), err}
				}
			})
			if forEachErr != nil {
				return LuaNbtError{"Error in byte array loop:", forEachErr}
			}
		} else {
			return LuaNbtError{fmt.Sprintf("Tag 7 Byte Array value field '%v' not a table", v), err}
		}
	case 8:
		if s, ok := v.(lua.LString); ok {
			err = binary.Write(w, byteOrder, int16(len([]byte(s))))
			if err != nil {
				return LuaNbtError{"Error writing string length", err}
			}
			err = binary.Write(w, byteOrder, []byte(s))
			if err != nil {
				return LuaNbtError{"Error writing string payload", err}
			}
		} else {
			return LuaNbtError{fmt.Sprintf("Tag 8 String value field '%v' not a string", v), err}
		}
	case 9:
		// important: tagListType needs to be in scope to be passed to writePayload
		var tagListType lua.LNumber
		var lv lua.LValue
		if lTable, ok := v.(*lua.LTable); ok {
			lv = L.RawGet(lTable, lua.LString("tagListType"))
			if tagListType, ok = lv.(lua.LNumber); ok {
				err = binary.Write(w, byteOrder, byte(tagListType))
				if err != nil {
					return LuaNbtError{"While writing tag 9 list type", err}
				}
			}
			lv = L.RawGet(lTable, lua.LString("list"))
			if values, ok := lv.(*lua.LTable); ok {
				err = binary.Write(w, byteOrder, int32(values.Len()))
				if err != nil {
					return LuaNbtError{"While writing tag 9 list size", err}
				}
				var forEachErr error
				values.ForEach(func(_ lua.LValue, n lua.LValue) {
					err = writePayload(w, n, tagListType, L)
					if err != nil && forEachErr == nil {
						forEachErr = LuaNbtError{"While writing tag 9 list of type " + strconv.Itoa(int(tagListType)), err}
					}
				})
				if forEachErr != nil {
					return forEachErr
				}
			} else {
				return LuaNbtError{fmt.Sprintf("Tag 9 List's value field '%v' not an array or null", lv), err}
			}

		} else {
			return LuaNbtError{fmt.Sprintf("Tag 9 List value field '%v' not an object", v), err}
		}
	case 10:
		if values, ok := v.(*lua.LTable); ok {
			var forEachErr error
			values.ForEach(func(_ lua.LValue, t lua.LValue) {
				if tag, ok := t.(*lua.LTable); ok {
					err = writeTag(w, tag, L)
					if err != nil && forEachErr == nil {
						forEachErr = LuaNbtError{"While writing Compound tags", err}
					}
				} else if forEachErr == nil {
					forEachErr = LuaNbtError{fmt.Sprintf("In tag type 10, expected table but got: %v", t), nil}
				}
			})
			// write the end tag which is just a single byte 0
			err = binary.Write(w, byteOrder, byte(0))
			if err != nil {
				return LuaNbtError{"Writing End tag", err}
			}
		} else {
			return LuaNbtError{fmt.Sprintf("Tag 10 Compound value field '%v' not an array", v), err}
		}
	case 11:
		if values, ok := v.(*lua.LTable); ok {
			err = binary.Write(w, byteOrder, int32(values.Len()))
			if err != nil {
				return LuaNbtError{"Error writing int32 array length", err}
			}
			var forEachErr error
			values.ForEach(func(_ lua.LValue, n lua.LValue) {
				if i, ok := n.(lua.LNumber); ok {
					if i < math.MinInt32 || i > math.MaxInt32 && forEachErr == nil {
						forEachErr = LuaNbtError{fmt.Sprintf("%v is out of range for Int in tag 11 - Int Array", i), nil}
					}
					err = binary.Write(w, byteOrder, int32(i))
					if err != nil && forEachErr == nil {
						forEachErr = LuaNbtError{"Error writing element of int32 array", err}
					}
				} else if forEachErr == nil {
					forEachErr = LuaNbtError{fmt.Sprintf("Tag 11 Int Array element value field '%v' not a number", n), err}
				}
			})
		} else {
			return LuaNbtError{fmt.Sprintf("Tag Int Array value field '%v' not an array", v), err}
		}
	case 12:
		if values, ok := v.(*lua.LTable); ok {
			err = binary.Write(w, byteOrder, int64(values.Len()))
			if err != nil {
				return LuaNbtError{"Error writing int64 array length", err}
			}
			var forEachErr error
			values.ForEach(func(_ lua.LValue, n lua.LValue) {
				if li, ok := n.(*lua.LTable); ok {
					var vl, vm lua.LNumber
					var lv lua.LValue
					lv = L.RawGet(li, lua.LString("least"))
					if vl, ok = lv.(lua.LNumber); !ok && forEachErr == nil {
						forEachErr = LuaNbtError{fmt.Sprintf("Error reading valueLeast of '%v'", lv), nil}
					}
					lv = L.RawGet(li, lua.LString("most"))
					if vm, ok = lv.(lua.LNumber); !ok && forEachErr == nil {
						forEachErr = LuaNbtError{fmt.Sprintf("Error reading valueMost of '%v'", lv), nil}
					}
					err = binary.Write(w, byteOrder, int64(intPairToLong(uint32(vl), uint32(vm))))
					if err != nil && forEachErr == nil {
						forEachErr = LuaNbtError{"Error writing int64 (from uint32 pair) payload:", err}
					}
				} else {
					forEachErr = LuaNbtError{fmt.Sprintf("Tag 4 Long value element field '%v' not a table", v), err}
				}
			})
			if forEachErr != nil {
				return forEachErr
			}
		} else {
			return LuaNbtError{fmt.Sprintf("Tag 12 Long Array element value field '%v' not an array", v), err}
		}
	default:
		return LuaNbtError{fmt.Sprintf("tagType '%v' is not recognized", tagType), err}
	}
	return err
}
