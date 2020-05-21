This repo contains a Go module/pakage to make NBT (named binary tags) readable
and editable in a Lua script environment, and a command-line interface similar
to native lua to run Lua scripts which can load, modify, and write NBT data.

## Lua NBT functions

- `use_bedrock_encoding()` - Sets future NBT encoding decoding using the Bedrock
Edition (little endian) format
- `use_java_encoding()` - Sets future NBT encoding decoding using the Java
Edition (little endian) format
- `loadnbt(path)` - Where `path` is a path to an NBT file, it will auto-detect
whether it's compressed and populate the `nbt` variable with its data
- `savenbt(path, compress)` - Converts `nbt` back to NBT and writes to `path`.
`compress` is `true` for compressed output and ommitted or `false` for
uncompressed output.

## Format of `nbt` variable in Lua

- lua's global `nbt` is a table `{}` in which each top-level nbt tag is
- in many cases there is only one top-level nbt compound tag, so `nbt[1]` is that tag, and `nbt[1][1]`, `nbt[1][2]`... are the first-tier tags you're looking for. Try `nbt[1][1].name` or the equivalent `nbt[1][1]["name"]`
- All tags (except tag 0 / end) are added as tables, and they have a `tagType`, `value`, and `name`
- Compound and list tags' values are again tables of the values beginning with `[1]`

## Lua examples

See /examples folder for example lua scripts.

## Go code example

See examples/go-example.go for how to use the package

### Go exported functions

- `func Nbt2Lua(b []byte, L *lua.LState) error` - pass it an uncompressed nbt byte array and the gopher-lua state variable, and it will populate the `nbt` global variable in Lua with a table hierarchy representing the nbt data
- `func Lua2Nbt(L *lua.LState) ([]byte, error)` - pass it the gopher-lua state variable, and it will convert the `nbt` global variable into an nbt byte array and return it
- `func UseBedrockEncoding()` - This makes any future conversions read/write the nbt usable by Minecraft Bedrock Edition (little endian). This is the default state when the package is loaded.
- `func UseJavaEncoding()` - This makes any future conversions read/write the nbt usable by Minecraft Java Edition (big endian)
- `func NewState() *lua.LState` - This can be used in place of calling lua.NewState for one less include in the client program, and it calls Nlua before returing LState
- `func Nlua(L *lua.LState)` - Nlua injects `loadnbt()` and (future) `savenbt()` functions into a lua environment