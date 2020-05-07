An attempt to expose NBT to a lua environment as a table for user scripting.

## Work in progress

- cmd/nbtlua and/or nbt2lua_test.go are creating a gopher-lua state (Lua v5.1) and calling Nbt2Lua with a hard-coded NBT byte array
- Nbt2Lua is creating `nbt` table in lua's global namespace
- cmd/nbtlua currently provides interactive readline and ability to run lua from an argument or a file to read `nbt`
- nbt2lua_test.go is testing tags 1, 2, 3, 5, and 6 for accurate conversion

## Status

- Successfully reading real NBT data, but currently from a hard-coded file location
- examples/tree-walk.lua works! `nbtlua tree-walk` works if you have an uncompressed nbt file named "player.dat" in the local directory

## Format of `nbt`

- lua's global `nbt` is a table `{}` in which each top-level nbt tag is
- in many cases there is only one top-level nbt tag, so `nbt[1]` is that tag
- All tags (except tag 0 / end) are added as tables, and they have a `tagType` and `value`, and many have a `name`
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
