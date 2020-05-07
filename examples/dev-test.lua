-- This script is for whatever I need to test in the moment

print("Hi from test script")

print(nbt[1].tagType)

-- overwrite nbt for testing purposes
nbt = {}
nbt[#nbt+1] = {
    tagType = 1,
    name = "byte",
    value = 5,
}
nbt[#nbt+1] = {
    tagType = 2,
    name = "short",
    value = 5,
}
nbt[#nbt+1] = {
    tagType = 3,
    name = "int",
    value = 5,
}

