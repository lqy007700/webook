
-- 限流对象
local key = KEYS[1]
-- 窗口大小
local window = tonumber(ARGV[1])
-- 阈值
local threshold = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

-- 窗口起始时间
local min = now - window

redis.call("zremrangebyscore", key, "-inf", min)
local cnt = redis.call("zcount", key, '-inf', '+inf')


-- 限流
if cnt >= threshold then
    return "true"
else
    redis.call("zadd", key, now, 1)
    redis.call("expire", key, window)
    return "false"
end
