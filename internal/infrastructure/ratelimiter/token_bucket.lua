-- KEYS[1] - ключ
-- ARGV[1] - лимит (максимум токенов)
-- ARGV[2] - скорость (токенов в секунду)
-- ARGV[3] - текущее время (ms)
-- ARGV[4] - стоимость запроса

local key = KEYS[1]
local limit = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local cost = tonumber(ARGV[4])

local data = redis.call("HMGET", key, "tokens", "last_refill")
local tokens = tonumber(data[1])
local last_refill = tonumber(data[2])

if tokens == nil then
  tokens = limit
  last_refill = now
end

local delta = math.max(0, now - last_refill) / 1000.0 * rate
tokens = math.min(limit, tokens + delta)
last_refill = now

local allowed = tokens >= cost
if allowed then
  tokens = tokens - cost
end

redis.call("HMSET", key, "tokens", tokens, "last_refill", last_refill)
redis.call("PEXPIRE", key, 60000) -- TTL для очистки

return { allowed and 1 or 0, tokens, last_refill }
