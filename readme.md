An attempt to expose NBT to a lua environment as a table for user scripting.

## Work in progress

- cmd/nbtlua and/or nbt2lua_test.go are creating a gopher-lua state (Lua v5.1) and calling Nbt2Lua with a hard-coded NBT byte array
- Nbt2Lua is creating `nbt` table in lua's global namespace
- cmd/nbtlua currently provides interactive readline and ability to run lua from an argument or a file to read `nbt`
- nbt2lua_test.go is testing tags 1, 2, 3, 5, and 6 for accurate conversion

## Status

- Ready to try to read real NBT data, but there's no convenient way to do it yet

## Vision

- The Go code will make the nbt easily accessible from lua
- Some basic helper funcions will be made available via lua
- Lua code will read and/or alter the data
- Go code will write the modified NBT to a file
- The library will be accessible for other Go projects, like my [MCPE Tool](https://github.com/midnightfreddie/McpeTool)
- The `"github.com/midnightfreddie/nbt-go-lua"` module will be kept simple and only decode/encode between nbt and lua
- Lua scripts and other projects can add more complex features
- cmd/nbtlua will eventually read and write optionally-compressed nbt files and handle the file reads & writes
