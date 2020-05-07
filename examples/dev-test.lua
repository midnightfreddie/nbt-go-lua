-- This script is for whatever I need to test in the moment

-- print("Hi from test script")

-- print(nbt[1].tagType)

-- overwrite nbt for testing purposes
nbt = {}

-- create tags to test Lua2Nbt conversion
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
nbt[#nbt+1] = {
    tagType = 4,
    name = "long",
    value = {
        least = 0xffffffff,
        most = 1,
    },
}
nbt[#nbt+1] = {
    tagType = 5,
    name = "float32",
    value = 5,
}
nbt[#nbt+1] = {
    tagType = 6,
    name = "float64",
    value = 5,
}
nbt[#nbt+1] = {
    tagType = 7,
    name = "byteArray",
    value = { 5, 6, 7, 8 },
}
nbt[#nbt+1] = {
    tagType = 8,
    name = "stringName",
    value = "stringString",
}
nbt[#nbt+1] = {
    tagType = 9,
    name = "taglist",
    value = {
        tagListType = 1,
        list = { 5, 6, 7, 8 },
    }
}
nbt[#nbt+1] = {
    tagType = 10,
    name = "compound",
    value = {
        {
            tagType = 1,
            name = "byte",
            value = 5,
        },
        {
            tagType = 2,
            name = "short",
            value = 5,
        },
        {
            tagType = 3,
            name = "int",
            value = 5,
        },
        {
            tagType = 4,
            name = "long",
            value = {
                least = 0xffffffff,
                most = 1,
            },
        },
    }
}
nbt[#nbt+1] = {
    tagType = 11,
    name = "intArray",
    value = { 5, 6, 7, 8 },
}
