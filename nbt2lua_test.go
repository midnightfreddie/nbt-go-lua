package nlua

import (
	"fmt"
	"math"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestValueConversions(t *testing.T) {

	tags := []struct {
		tagType lua.LNumber
		lType   lua.LValueType
		value   lua.LValue
		nbt     []byte
	}{
		{1, lua.LTNumber, lua.LNumber(math.MaxInt8), []byte{1, 0, 0, 0x7f}},
		{1, lua.LTNumber, lua.LNumber(math.MinInt8), []byte{1, 0, 0, 0x80}},
		{2, lua.LTNumber, lua.LNumber(math.MaxInt16), []byte{2, 0, 0, 0xff, 0x7f}},
		{2, lua.LTNumber, lua.LNumber(math.MinInt16), []byte{2, 0, 0, 0x00, 0x80}},
		{3, lua.LTNumber, lua.LNumber(math.MaxInt32), []byte{3, 0, 0, 0xff, 0xff, 0xff, 0x7f}},
		{3, lua.LTNumber, lua.LNumber(math.MinInt32), []byte{3, 0, 0, 0x00, 0x00, 0x00, 0x80}},
	}

	UseBedrockEncoding()
	L := lua.NewState()
	for _, tag := range tags {
		err := Nbt2Lua(tag.nbt, L)
		if err != nil {
			t.Error(fmt.Sprintf("foo: %s", err.Error()))
		}
	}
}
