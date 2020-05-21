package nlua

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math"
	"path/filepath"
	"runtime"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestNbt2Lua(t *testing.T) {

	numberTags := []struct {
		tagType lua.LNumber
		value   lua.LNumber
		nbt     []byte
	}{
		{1, math.MaxInt8, []byte{1, 0, 0, 0x7f}},
		{1, math.MinInt8, []byte{1, 0, 0, 0x80}},
		{2, math.MaxInt16, []byte{2, 0, 0, 0xff, 0x7f}},
		{2, math.MinInt16, []byte{2, 0, 0, 0x00, 0x80}},
		{3, math.MaxInt32, []byte{3, 0, 0, 0xff, 0xff, 0xff, 0x7f}},
		{3, math.MinInt32, []byte{3, 0, 0, 0x00, 0x00, 0x00, 0x80}},
		{5, 0, []byte{5, 0, 0, 0x00, 0x00, 0x00, 0x00}},
		{5, math.MaxFloat32, []byte{5, 0, 0, 0xff, 0xff, 0x7f, 0x7f}},
		{5, math.SmallestNonzeroFloat32, []byte{5, 0, 0, 0x01, 0x00, 0x00, 0x00}},
		{6, 0, []byte{6, 0, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{6, math.MaxFloat64, []byte{6, 0, 0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xef, 0x7f}},
		{6, math.SmallestNonzeroFloat64, []byte{6, 0, 0, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}

	UseBedrockEncoding()
	L := NewState()
	defer L.Close()
	for _, tag := range numberTags {
		err := Nbt2Lua(tag.nbt, L)
		if err != nil {
			t.Error(fmt.Sprintf("Error processing nbt: `%s` NBT hex dump:\n%s", err.Error(), hex.Dump(tag.nbt)))
		} else {
			lNbt := L.GetGlobal("nbt")
			if lNbtTable, ok := lNbt.(*lua.LTable); ok {
				lNbtTable.ForEach(func(k lua.LValue, v lua.LValue) {
					if lTag, ok := v.(*lua.LTable); ok {
						name := L.RawGet(lTag, lua.LString("name"))
						if sName, ok := name.(lua.LString); ok {
							if sName.String() != "" {
								t.Error(fmt.Sprintf("Name not empty string: %s", sName))
							}
						} else {
							t.Error("Name is not a string:", name.Type())
						}
						tagType := L.RawGet(lTag, lua.LString("tagType"))
						if tagTypeNumber, ok := tagType.(lua.LNumber); ok {
							if tagTypeNumber != tag.tagType {
								t.Error(fmt.Sprintf("Expected %v, got %v", tag.tagType, tagTypeNumber))
							}
						} else {
							t.Error("tagType is not a number", tagType.Type())
						}
						value := L.RawGet(lTag, lua.LString("value"))
						if valueTyped, ok := value.(lua.LNumber); ok {
							if valueTyped != tag.value {
								t.Error(fmt.Sprintf("Value expected %v, got %v", tag.value, valueTyped))
							}
						} else {
							t.Error(fmt.Sprintf("Value type expected lua.LNumber, got %v", value))
						}
					}
				})
			}
		}
	}
}

func TestLua2Nbt(t *testing.T) {
	L := NewState()
	defer L.Close()
	var bedrockSig, javaSig []byte
	// get filename of current file; will use relative path from here for test data input
	_, filename, _, _ := runtime.Caller(0)
	luaFile := filepath.Dir(filename) + "/test_data/testnbt.lua"
	if err := L.DoFile(luaFile); err != nil {
		t.Fatal("Error running lua script: ", err)
	}

	// read sha1bedrock sig from lua
	lv := L.GetGlobal("sha1bedrock")
	if sig, ok := lv.(*lua.LTable); ok {
		sig.ForEach(func(_ lua.LValue, v lua.LValue) {
			if n, ok := v.(lua.LNumber); ok {
				bedrockSig = append(bedrockSig, byte(n))
			} else {
				t.Errorf("sha1bedrock array element expected LNumber, got %v", v)
			}
		})
	} else {
		t.Errorf("sha1bedrock expected table, got %v", lv)
	}
	// read sha1java sig from lua
	lv = L.GetGlobal("sha1java")
	if sig, ok := lv.(*lua.LTable); ok {
		sig.ForEach(func(_ lua.LValue, v lua.LValue) {
			if n, ok := v.(lua.LNumber); ok {
				javaSig = append(javaSig, byte(n))
			} else {
				t.Errorf("sha1java array element expected LNumber, got %v", v)
			}
		})
	} else {
		t.Errorf("sha1java expected table, got %v", lv)
	}

	// Bedrock Lua2Nbt sha1 check
	UseBedrockEncoding()
	nbtOut, err := Lua2Nbt(L)
	if err != nil {
		t.Error("Bedrock conversion: ", err)
	}
	s := sha1.Sum(nbtOut)
	if !bytes.Equal(s[:], bedrockSig) {
		t.Errorf("bedrock signature expected %v, got %v", bedrockSig, s)
	}

	// Round trip check
	if err := Nbt2Lua(nbtOut, L); err == nil {
		if nbtOut, err = Lua2Nbt(L); err == nil {
			s := sha1.Sum(nbtOut)
			if !bytes.Equal(s[:], bedrockSig) {
				t.Errorf("bedrock round trip reconversion signature expected %v, got %v", bedrockSig, s)
			}
		} else {
			t.Error("Bedrock round trip reconversion: ", err)
		}
	} else {
		t.Error("Bedrock sha1 round trip: ", err)
	}

	// re-run script for fresh start with next tests
	if err := L.DoFile(luaFile); err != nil {
		t.Fatal("Error running lua script: ", err)
	}

	// Java Lua2Nbt sha1 check
	UseJavaEncoding()
	nbtOut, err = Lua2Nbt(L)
	if err != nil {
		t.Error("Java conversion: ", err)
	}
	s = sha1.Sum(nbtOut)
	if !bytes.Equal(s[:], javaSig) {
		t.Errorf("Java signature expected %v, got %v", javaSig, s)
	}
	// Round trip check
	if err := Nbt2Lua(nbtOut, L); err == nil {
		if nbtOut, err = Lua2Nbt(L); err == nil {
			s := sha1.Sum(nbtOut)
			if !bytes.Equal(s[:], javaSig) {
				t.Errorf("Java round trip reconversion signature expected %v, got %v", javaSig, s)
			}
		} else {
			t.Error("Java round trip reconversion: ", err)
		}
	} else {
		t.Error("Java sha1 round trip: ", err)
	}
}
