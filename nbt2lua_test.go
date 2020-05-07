package nlua

import (
	"encoding/hex"
	"fmt"
	"math"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestValueConversions(t *testing.T) {

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
	L := lua.NewState()
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
