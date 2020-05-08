This repo is a Go module/pakage to make NBT (named binary tags) readable and editable in a Lua script environment.

## Work in progress / status

- The Go module is working and successfully decoding/encoding NBT (named binary tags, mostly used in Minecraft) into a Lua 5.1 environment in the global `nbt` variable
- Lua scripts in examples and test_data can give you an idea of how to access or modify it
- The cmd/nbtlua executable is not really in a distributable state yet; I've been hard coding things while getting the nlua package finished

## Format of `nbt`

- lua's global `nbt` is a table `{}` in which each top-level nbt tag is
- in many cases there is only one top-level nbt compound tag, so `nbt[1]` is that tag, and `nbt[1][1]`, `nbt[1][2]`... are the first-tier tags you're looking for. Try `nbt[1][1].name` or the equivalent `nbt[1][1]["name"]`
- All tags (except tag 0 / end) are added as tables, and they have a `tagType`, `value`, and `name`
- Compound and list tags' values are again tables of the values beginning with `[1]`

## Vision

- The Go code will make the nbt easily accessible from lua
- Some basic helper funcions will be made available via lua
- Lua code will read and/or alter the data
- Go code will write the modified NBT to a file
- The library will be accessible for other Go projects, like my [MCPE Tool](https://github.com/midnightfreddie/McpeTool)
- The `"github.com/midnightfreddie/nbt-go-lua"` module will be kept simple and only decode/encode between nbt and lua
- Lua scripts and other projects can add more complex features
- cmd/nbtlua will eventually read and write optionally-compressed nbt files and handle the file reads & writes

## Go code example

See examples/go-example.go for how to use the package

### Go exported functions

- `func Nbt2Lua(b []byte, L *lua.LState) error` - pass it an uncompressed nbt byte array and the gopher-lua state variable, and it will populate the `nbt` global variable in Lua with a table hierarchy representing the nbt data
- `func Lua2Nbt(L *lua.LState) ([]byte, error)` - pass it the gopher-lua state variable, and it will convert the `nbt` global variable into an nbt byte array and return it
- `func UseBedrockEncoding()` - This makes any future conversions read/write the nbt usable by Minecraft Bedrock Edition (little endian). This is the default state when the package is loaded.
- `func UseJavaEncoding()` - This makes any future conversions read/write the nbt usable by Minecraft Java Edition (big endian)