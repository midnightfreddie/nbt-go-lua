-- pass it a table and optinoally hierarchy number
function printNbtTables(t, n)
    -- default hierarchy n to 0 if not set
    n = n or 0
    -- indent string, length based on hierarchy
    local s = string.rep("  ", n)
    -- loop through all key/value pairs of the table
    for k,v in pairs(t) do
        -- if value is a table, recurse and increase indentation
        if type(v) == "table" then
            print(s..k..":")
            printNbtTables(v,n+1)
        else
            -- otherwise print the key/value pair which are usually "tagType" and "name"
            print(s..k, v)
        end
    end
end

printNbtTables(nbt)
