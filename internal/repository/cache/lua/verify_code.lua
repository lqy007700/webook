local key = KEYS[1]
local exceptedCode = ARGV[1]
local code = redis.call("get", key)
local cntKey = key .. ":cnt"
local cnt = tonumber(redic.call("get", cntKey))

if cnt <= 0 then
    return -1
elseif exceptedCode == code then
    return 0
elseif
    return -2
end