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
func Lua2Nbt(L *lua.LState) ([]byte, error) {
	nbtOut := new(bytes.Buffer)
	nbtArray := L.GetGlobal("nbt")
	var forEachErr error
	if nbtLuaTable, ok := nbtArray.(*lua.LTable); ok {
		// TODO: decide how to handle empty nbt table; for now it will return null byte array
		//  This is not expected to be a sane situation to use this function, although it's technically correct
		// TODO: decide how to handle non-numeric nbt table keys; for now they will process the same as numeric ones
		//   nbt table produced by this package should not have non-numeric keys in nbt, so might want to ignore or error
		nbtLuaTable.ForEach(func(k lua.LValue, v lua.LValue) {
			fmt.Println(k)
			if nbtLuaTag, ok := v.(*lua.LTable); ok {
				err := writeTag(nbtOut, nbtLuaTag, L)
				if err != nil {
					// FIXME: How do I propagate an error out from the anonymous function inside a loop?
					// return nil, err
					fmt.Println("Error in top loop:", err.Error())
					// currently putting the first error in the function exit return
					if forEachErr == nil {
						forEachErr = err
					}
				}
			}
		})
	} else {
		// throw error if `nbt` is not a lua table
		return nil, LuaNbtError{fmt.Sprintf("Global nbt type, expected %T, got %T", lua.LTable{}, nbtArray), nil}
	}
	return nbtOut.Bytes(), forEachErr
}

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
		/*
			case 4:
				if int64Map, ok := v.(map[string]interface{}); ok {
					var nbtLong NbtLong
					var vl, vm float64
					if vl, ok = int64Map["valueLeast"].(float64); !ok {
						return LuaNbtError{fmt.Sprintf("Error reading valueLeast of '%v'", int64Map["valueLeast"]), nil}
					}
					nbtLong.ValueLeast = uint32(vl)
					if vm, ok = int64Map["valueMost"].(float64); !ok {
						return LuaNbtError{fmt.Sprintf("Error reading valueMost of '%v'", int64Map["valueMost"]), nil}
					}
					nbtLong.ValueMost = uint32(vm)
					err = binary.Write(w, byteOrder, int64(intPairToLong(nbtLong)))
					if err != nil {
						return LuaNbtError{"Error writing int64 (from uint32 pair) payload:", err}
					}
					} else if int64String, ok := v.(string); ok {
						i, err := strconv.ParseInt(int64String, 10, 64)
						if err != nil {
							return LuaNbtError{"Error converting long as string payload:", err}
						}
						err = binary.Write(w, byteOrder, i)
						if err != nil {
							return LuaNbtError{"Error writing int64 (from string) payload:", err}
						}
						if err != nil {
							return LuaNbtError{fmt.Sprintf("Tag 4 Long value string field '%s' not an integer", int64String), err}
						}
						} else {
							return LuaNbtError{fmt.Sprintf("Tag 4 Long value field '%v' not an object", v), err}
						}
		*/
	case 5:
		if f, ok := v.(lua.LNumber); ok {
			// Comparing to smallest/max values is causing false errors as the nbt value comes out right even if this comparison doesn't
			// Instead, will check for positive/negative infinity. Not sure what happens if f is too small for float32, but is likely edge case
			// if f != 0 && (math.Abs(f) < math.SmallestNonzeroFloat32 || math.Abs(f) > math.MaxFloat32) {
			if math.IsInf(float64(float32(f)), 0) {
				return LuaNbtError{fmt.Sprintf("%g is out of range for tag 5 - Float", f), nil}
			}
			err = binary.Write(w, byteOrder, float32(f))
			if err != nil {
				return LuaNbtError{"Error writing float32 payload", err}
			}
		} else {
			// TODO: find if lua.LNumber can represent NaN and if it convierts
			//   If so, this else block may need to throw an error
			// If NaN is valid for double, maybe it's valid for float?
			// return LuaNbtError{fmt.Sprintf("Tag 5 Float value field '%v' not a number", v), err}
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
			// TODO: find if lua.LNumber can represent NaN and if it convierts
			//   If so, this else block may need to throw an error
			// Apparently NaN is a valid value in Minecraft for double?
			// return LuaNbtError{fmt.Sprintf("Tag 6 Double value field '%v' not a number", v), err}
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
		// := were keeping it in a lower scope and zeroing it out.
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
					// for _, value := range values {
					// fakeTag := make(map[string]interface{})
					// fakeTag["value"] = value
					err = writePayload(w, n, tagListType, L)
					if err != nil && forEachErr == nil {
						forEachErr = LuaNbtError{"While writing tag 9 list of type " + strconv.Itoa(int(tagListType)), err}
					}
				})
				if forEachErr != nil {
					return forEachErr
				}
				/* TODO: I think lua nbt empty tables will work, but check
				} else if lTable["list"] == nil {
					// NBT lists can be null / nil and therefore aren't represented as an array in JSON
					err = binary.Write(w, byteOrder, int32(0))
					if err != nil {
						return LuaNbtError{"While writing tag 9 list null size", err}
					}
					return nil
				*/
			} else {
				return LuaNbtError{fmt.Sprintf("Tag 9 List's value field '%v' not an array or null", lv), err}
			}

		} else {
			return LuaNbtError{fmt.Sprintf("Tag 9 List value field '%v' not an object", v), err}
		}
		/*
			case 10:
				if values, ok := v.([]interface{}); ok {
					for _, value := range values {
						err = writeTag(w, value)
						if err != nil {
							return LuaNbtError{"While writing Compound tags", err}
						}
					}
					// write the end tag which is just a single byte 0
					err = binary.Write(w, byteOrder, byte(0))
					if err != nil {
						return LuaNbtError{"Writing End tag", err}
					}
				} else {
					return LuaNbtError{fmt.Sprintf("Tag 10 Compound value field '%v' not an array", v), err}
				}
			case 11:
				if values, ok := v.([]interface{}); ok {
					err = binary.Write(w, byteOrder, int32(len(values)))
					if err != nil {
						return LuaNbtError{"Error writing int32 array length", err}
					}
					for _, value := range values {
						if i, ok := value.(float64); ok {
							if i < math.MinInt32 || i > math.MaxInt32 {
								return LuaNbtError{fmt.Sprintf("%v is out of range for Int in tag 11 - Int Array", i), nil}
							}
							err = binary.Write(w, byteOrder, int32(i))
							if err != nil {
								return LuaNbtError{"Error writing element of int32 array", err}
							}
						} else {
							return LuaNbtError{fmt.Sprintf("Tag 11 Int Array element value field '%v' not an integer", value), err}
						}
					}
				} else {
					return LuaNbtError{fmt.Sprintf("Tag Int Array value field '%v' not an array", v), err}
				}
			case 12:
				if values, ok := v.([]interface{}); ok {
					err = binary.Write(w, byteOrder, int64(len(values)))
					if err != nil {
						return LuaNbtError{"Error writing int64 array length", err}
					}
					for _, value := range values {
						if int64Map, ok := value.(map[string]interface{}); ok {
							var nbtLong NbtLong
							var vl, vm float64
							if vl, ok = int64Map["valueLeast"].(float64); !ok {
								return LuaNbtError{fmt.Sprintf("Error reading valueLeast of '%v'", int64Map["valueLeast"]), nil}
							}
							nbtLong.ValueLeast = uint32(vl)
							if vm, ok = int64Map["valueMost"].(float64); !ok {
								return LuaNbtError{fmt.Sprintf("Error reading valueMost of '%v'", int64Map["valueMost"]), nil}
							}
							nbtLong.ValueMost = uint32(vm)
							// if i, ok := value.(float64); ok {
							err = binary.Write(w, byteOrder, int64(intPairToLong(nbtLong)))
							if err != nil {
								return LuaNbtError{"Error writing element of int64 array", err}
							}
						} else if int64String, ok := value.(string); ok {
							i, err := strconv.ParseInt(int64String, 10, 64)
							if err != nil {
								return LuaNbtError{"Error converting long array element as string payload:", err}
							}
							err = binary.Write(w, byteOrder, i)
							if err != nil {
								return LuaNbtError{"Error writing int64 array element (from string) payload", err}
							}
							if err != nil {
								return LuaNbtError{fmt.Sprintf("Tag 4 Long Array element value string field '%s' not an integer", int64String), err}
							}
						} else {
							return LuaNbtError{fmt.Sprintf("Tag Long Array element value field '%v' not an object", value), err}
						}
					}
				} else {
					return LuaNbtError{fmt.Sprintf("Tag 12 Long Array element value field '%v' not an array", v), err}
				}
		*/
	default:
		return LuaNbtError{fmt.Sprintf("tagType '%v' is not recognized", tagType), err}
	}
	return err
}
