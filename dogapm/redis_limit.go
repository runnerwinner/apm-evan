package dogapm

import (
	"context"

	"github.com/redis/go-redis/v9"
)
type redisLimit struct {}

const limitScript = `
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local expire = tonumber(ARGV[2])

local count = redis.call("INCR", key)
if count == 1 then
	redis.call("EXPIRE", key, expire)
end	
if count > limit then
	return 1
end
return 0
`	

var RedisLimiter = &redisLimit{}

func (r *redisLimit) IsLimit(client *redis.Client, key string, limitCnt int64, expireSeconds int) bool {
	sha,err := client.ScriptLoad(context.Background(), limitScript).Result()
	if err != nil {
		return false
	}
	result,err := client.EvalSha(context.Background(), sha, []string{key}, limitCnt, expireSeconds).Result()
	if err != nil {
		return false
	}
	if result.(int64) == 1 {
		return true
	}
	return false
}