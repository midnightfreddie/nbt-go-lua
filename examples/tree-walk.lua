-- top level is an array of concatenated nbt top-tags (Often just the one compound tag)
for k,v in ipairs(nbt) do
    print("=== A top-of-hierarchy tag ===")
    -- usually a compound tag
    for k2,v2 in pairs(v) do
        if k2=="value" and type(v2)=="table" then
            for k3,v3 in pairs(v2) do
            if k3=="value" and type(v3)=="table" then
                print("--- A second-tier tag ---")
                print("Key", k3, "Value", v3)
            else
                print(k3, v3)
            end
        end
    else
        print(k2,v2)
    end
    end

end
