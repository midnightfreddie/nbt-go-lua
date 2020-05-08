-- overwrite nbt for testing purposes
nbt = {}

-- create tags to test Lua2Nbt conversion
nbt[#nbt+1] = {
    tagType = 0,
    name = "endTag",
    value = "This tag should be ignored as it is meaningless and shouldn't be in the lua representation",
}nbt[#nbt+1] = {
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
nbt[#nbt+1] = {
    tagType = 12,
    name = "longArray",
    value = {
        {
            least = 0xffffffff,
            most = 1,
        },
        {
            least = 0xffffffff,
            most = 0x7fffffff,
        },
        {
            least = 0,
            most = 0x80000000,
        },

    },
}

-- sha1 signatures of the nbt output of the above
-- sha1bedrock = { 116, 110, 134, 107, 34, 229, 46, 1, 116, 206, 74, 130, 62, 200, 45, 242, 228, 223, 207, 130, }
-- sha1java = { 158, 160, 55, 155, 162, 162, 219, 130, 111, 181, 84, 136, 17, 34, 222, 77, 174, 238, 224, 152, }

sha1bedrock = { 183, 26, 31, 26, 110, 33, 191, 178, 15, 197, 19, 165, 71, 31, 223, 242, 206, 152, 249, 212, }
sha1java = { 137, 230, 39, 223, 1, 53, 121, 63, 213, 86, 196, 239, 214, 171, 194, 132, 222, 16, 188, 97, }
