-- This script is for whatever I need to test in the moment

print("Hi from test script")

print(nbt[1].tagType)

-- overwrite nbt for testing purposes
nbt = {}
nbt[1] = {
    tagType = 1,
    name = "HelloByte",
    value = 5,
}

