-- pass it a table and optinoally hierarchy number
function printNbtTables(t, n)
    n = n or 0
    local s = string.rep("  ", n)
    for k,v in pairs(t) do
        if type(v) == "table" then
            print(s, k, ":")
            printNbtTables(v,n+1)
        else
            print(s, k, v)
        end
    end
end

printNbtTables(nbt)
