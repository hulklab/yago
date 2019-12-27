package rds

/**
 * Keys
 */

func (r *Rds) Del(keys ...interface{}) (reply interface{}, err error) {
	return r.Do("DEL", keys...)
}

func (r *Rds) Exists(key interface{}) (reply interface{}, err error) {
	return r.Do("EXISTS", key)
}

func (r *Rds) Expire(key, seconds interface{}) (reply interface{}, err error) {
	return r.Do("EXPIRE", key, seconds)
}

func (r *Rds) ExpireAt(key, timestamp interface{}) (reply interface{}, err error) {
	return r.Do("EXPIREAT", key, timestamp)
}

func (r *Rds) Keys(pattern interface{}) (reply interface{}, err error) {
	return r.Do("KEYS", pattern)
}

func (r *Rds) Persist(key interface{}) (reply interface{}, err error) {
	return r.Do("PERSIST", key)
}

func (r *Rds) PExpire(key, milliseconds interface{}) (reply interface{}, err error) {
	return r.Do("PEXPIRE", key)
}

func (r *Rds) PExpireAt(key, millisecondsTimestamp interface{}) (reply interface{}, err error) {
	return r.Do("PEXPIREAT", key)
}

func (r *Rds) PTtl(key interface{}) (reply interface{}, err error) {
	return r.Do("PTTL", key)
}

func (r *Rds) Ttl(key interface{}) (reply interface{}, err error) {
	return r.Do("TTL", key)
}

func (r *Rds) Type(key interface{}) (reply interface{}, err error) {
	return r.Do("TYPE", key)
}

// SCAN cursor [MATCH pattern] [COUNT count]
func (r *Rds) Scan(cursor interface{}, opts ...interface{}) (reply interface{}, err error) {
	args := mergeKeyAndArgs(cursor, opts...)
	return r.Do("TYPE", args...)
}

/*
 * Strings
 */

func (r *Rds) Decr(key interface{}) (reply interface{}, err error) {
	return r.Do("DECR", key)
}

func (r *Rds) DecrBy(key, decrement interface{}) (reply interface{}, err error) {
	return r.Do("DECRBY", key)
}

func (r *Rds) Get(key interface{}) (reply interface{}, err error) {
	return r.Do("GET", key)
}

func (r *Rds) GetBit(key, offset interface{}) (reply interface{}, err error) {
	return r.Do("GETBIT", key, offset)
}

func (r *Rds) GetRange(key, start, end interface{}) (reply interface{}, err error) {
	return r.Do("GETRANGE", key, start, end)
}

func (r *Rds) GetSet(key, value interface{}) (reply interface{}, err error) {
	return r.Do("GETSET", key, value)
}

func (r *Rds) Incr(key interface{}) (reply interface{}, err error) {
	return r.Do("INCR", key)
}

func (r *Rds) IncrBy(key, increment interface{}) (reply interface{}, err error) {
	return r.Do("INCRBY", key, increment)
}

func (r *Rds) IncrByFloat(key, increment interface{}) (reply interface{}, err error) {
	return r.Do("INCRBYFLOAT", key, increment)
}

func (r *Rds) MGet(keys ...interface{}) (reply interface{}, err error) {
	return r.Do("MGET", keys...)
}

// MSET key value [key value ...]
func (r *Rds) MSet(keyValues ...interface{}) (reply interface{}, err error) {
	return r.Do("MSET", keyValues...)
}

// MSETNX key value [key value ...]
func (r *Rds) MSetNx(keyValues ...interface{}) (reply interface{}, err error) {
	return r.Do("MSETNX", keyValues...)
}

// PSETEX key milliseconds value
func (r *Rds) PSetEx(key, milliseconds, value interface{}) (reply interface{}, err error) {
	return r.Do("PSETEX", key, milliseconds, value)
}

// SET key value [EX seconds] [PX milliseconds] [NX|XX]
func (r *Rds) Set(key, value interface{}, args ...interface{}) (reply interface{}, err error) {
	args2 := mergeKeyAndArgs(value, args...)
	nargs := mergeKeyAndArgs(key, args2...)
	return r.Do("SET", nargs...)
}

func (r *Rds) SetBit(key, offset, value interface{}) (reply interface{}, err error) {
	return r.Do("SETBIT", key, offset, value)
}

func (r *Rds) SetEx(key, seconds, value interface{}) (reply interface{}, err error) {
	return r.Do("SETEX", key, seconds, value)
}

func (r *Rds) SetNx(key, value interface{}) (reply interface{}, err error) {
	return r.Do("SETNX", key, value)
}

func (r *Rds) SetRange(key, offset, value interface{}) (reply interface{}, err error) {
	return r.Do("SETRANGE", key, offset, value)
}

func (r *Rds) StrLen(key interface{}) (reply interface{}, err error) {
	return r.Do("STRLEN", key)
}

/**
 * Lists
 */

func (r *Rds) BLPop(key, timeout interface{}) (reply interface{}, err error) {
	return r.Do("BLPOP", key, timeout)
}

func (r *Rds) BRPop(key, timeout interface{}) (reply interface{}, err error) {
	return r.Do("BRPOP", key, timeout)
}

func (r *Rds) BRPopLPush(source, destination, timeout interface{}) (reply interface{}, err error) {
	return r.Do("BRPOPLPUSH", source, destination, timeout)
}

func (r *Rds) LIndex(key, index interface{}) (reply interface{}, err error) {
	return r.Do("lindex", key, index)
}

// LINSERT key BEFORE|AFTER pivot value
func (r *Rds) LInsert(key, where, pivot, value interface{}) (reply interface{}, err error) {
	return r.Do("LINSERT", key, where, pivot, value)
}

func (r *Rds) LLen(key interface{}) (reply interface{}, err error) {
	return r.Do("LLEN", key)
}

func (r *Rds) LPop(key interface{}) (reply interface{}, err error) {
	return r.Do("LPOP", key)
}

func (r *Rds) LPush(key interface{}, values ...interface{}) (reply interface{}, err error) {
	args := mergeKeyAndArgs(key, values...)
	return r.Do("LPUSH", args...)
}

func (r *Rds) LPushX(key, value interface{}) (reply interface{}, err error) {
	return r.Do("LPUSHX", key, value)
}

func (r *Rds) LRange(key, start, stop interface{}) (reply interface{}, err error) {
	return r.Do("LRANGE", key, start, stop)
}

func (r *Rds) LRem(key, count, value interface{}) (reply interface{}, err error) {
	return r.Do("LREM", key, count, value)
}

func (r *Rds) LSet(key, index, value interface{}) (reply interface{}, err error) {
	return r.Do("LSET", key, index, value)
}

func (r *Rds) LTrim(key, start, stop interface{}) (reply interface{}, err error) {
	return r.Do("LTRIM", key, start, stop)
}

func (r *Rds) RPop(key interface{}) (reply interface{}, err error) {
	return r.Do("RPOP", key)
}

func (r *Rds) RPopLPush(source, destination interface{}) (reply interface{}, err error) {
	return r.Do("RPOPLPUSH", source, destination)
}

func (r *Rds) RPush(key interface{}, values ...interface{}) (reply interface{}, err error) {
	args := mergeKeyAndArgs(key, values...)
	return r.Do("RPUSH", args...)
}

func (r *Rds) RPushX(key, value interface{}) (reply interface{}, err error) {
	return r.Do("RPUSHX", key, value)
}

/**
 * Hashes
 */

func (r *Rds) HDel(key interface{}, fields ...interface{}) (reply interface{}, err error) {
	args := mergeKeyAndArgs(key, fields...)
	return r.Do("HDEL", args...)
}

func (r *Rds) HExists(key, field interface{}) (reply interface{}, err error) {
	return r.Do("HEXISTS", key, field)
}

func (r *Rds) HGet(key, field interface{}) (reply interface{}, err error) {
	return r.Do("HGET", key, field)
}

func (r *Rds) HGetAll(key interface{}) (reply interface{}, err error) {
	return r.Do("HGETALL", key)
}

func (r *Rds) HIncrBy(key, field, increment interface{}) (reply interface{}, err error) {
	return r.Do("HINCRBY", key, field, increment)
}

func (r *Rds) HIncrByFloat(key, field, increment interface{}) (reply interface{}, err error) {
	return r.Do("HINCRBYFLOAT", key, field, increment)
}

func (r *Rds) HKeys(key interface{}) (reply interface{}, err error) {
	return r.Do("HKEYS", key)
}

func (r *Rds) HLen(key interface{}) (reply interface{}, err error) {
	return r.Do("HLEN", key)
}

func (r *Rds) HMGet(key interface{}, fields ...interface{}) (reply interface{}, err error) {
	args := mergeKeyAndArgs(key, fields...)
	return r.Do("HMGET", args...)
}

// HMSET key field value [field value ...]
func (r *Rds) HMSet(key interface{}, fieldValues ...interface{}) (reply interface{}, err error) {
	args := mergeKeyAndArgs(key, fieldValues...)
	return r.Do("HMSET", args...)
}

func (r *Rds) HSet(key, field, value interface{}) (reply interface{}, err error) {
	return r.Do("HSET", key, field, value)
}

func (r *Rds) HSetNx(key, field, value interface{}) (reply interface{}, err error) {
	return r.Do("HSETNX", key, field, value)
}

func (r *Rds) HStrLen(key, field interface{}) (reply interface{}, err error) {
	return r.Do("HSTRLEN", key, field)
}

func (r *Rds) HVals(key interface{}) (reply interface{}, err error) {
	return r.Do("HVALS", key)
}

func (r *Rds) HScan(key interface{}, cursor interface{}, opts ...interface{}) (reply interface{}, err error) {
	args1 := mergeKeyAndArgs(cursor, opts...)
	args := mergeKeyAndArgs(key, args1...)
	return r.Do("HSCAN", args...)
}

/**
 * Sets
 */
func (r *Rds) SAdd(key interface{}, members ...interface{}) (reply interface{}, err error) {
	args := mergeKeyAndArgs(key, members...)
	return r.Do("SADD", args...)
}

func (r *Rds) SCard(key interface{}) (reply interface{}, err error) {
	return r.Do("SCARD", key)
}

func (r *Rds) SDiff(key interface{}, keys ...interface{}) (reply interface{}, err error) {
	args := mergeKeyAndArgs(key, keys...)
	return r.Do("SDIFF", args...)
}

func (r *Rds) SDiffStore(destination, key interface{}, keys ...interface{}) (reply interface{}, err error) {
	args1 := mergeKeyAndArgs(key, keys...)
	args := mergeKeyAndArgs(destination, args1...)
	return r.Do("SDIFFSTORE", args...)
}

func (r *Rds) SInter(key interface{}, keys ...interface{}) (reply interface{}, err error) {
	args := mergeKeyAndArgs(key, keys...)
	return r.Do("SINTER", args...)
}

func (r *Rds) SInterStore(destination, key interface{}, keys ...interface{}) (reply interface{}, err error) {
	args1 := mergeKeyAndArgs(key, keys...)
	args := mergeKeyAndArgs(destination, args1...)
	return r.Do("SINTERSTORE", args...)
}

func (r *Rds) SIsMember(key, member interface{}) (reply interface{}, err error) {
	return r.Do("SISMEMBER", key, member)
}

func (r *Rds) SMembers(key interface{}) (reply interface{}, err error) {
	return r.Do("SMEMBERS", key)
}

func (r *Rds) SMove(source, destination, member interface{}) (reply interface{}, err error) {
	return r.Do("SMOVE", source, destination, member)
}

func (r *Rds) SPop(key interface{}) (reply interface{}, err error) {
	return r.Do("SPOP", key)
}

func (r *Rds) SRem(key interface{}, members ...interface{}) (reply interface{}, err error) {
	args := mergeKeyAndArgs(key, members...)
	return r.Do("SREM", args...)
}

func (r *Rds) SUnion(key interface{}, keys ...interface{}) (reply interface{}, err error) {
	args := mergeKeyAndArgs(key, keys...)
	return r.Do("SUNION", args...)
}

func (r *Rds) SUnionStore(key interface{}, keys ...interface{}) (reply interface{}, err error) {
	args := mergeKeyAndArgs(key, keys...)
	return r.Do("SUNIONSTORE", args...)
}

func (r *Rds) SScan(key interface{}, cursor interface{}, options ...interface{}) (reply interface{}, err error) {
	args1 := mergeKeyAndArgs(cursor, options...)
	args := mergeKeyAndArgs(key, args1...)
	return r.Do("SSCAN", args...)
}

/*
 * Sorted Sets
 */

// zadd key [nx|xx] [ch] [incr] score member [score member...]
func (r *Rds) ZAdd(key interface{}, options ...interface{}) (reply interface{}, err error) {
	args := mergeKeyAndArgs(key, options...)
	return r.Do("ZADD", args...)
}

func (r *Rds) ZCard(key interface{}) (reply interface{}, err error) {
	return r.Do("ZCARD", key)
}

func (r *Rds) ZCount(key, min, max interface{}) (reply interface{}, err error) {
	return r.Do("ZCOUNT", key, min, max)
}

func (r *Rds) ZIncrBy(key, increment, member interface{}) (reply interface{}, err error) {
	return r.Do("ZINCRBY", key, increment, member)
}

// zinterstore destination numkeys key [key ...] [weights weight]
func (r *Rds) ZInterStore(destination, numkeys, key1, key2 interface{}, opts ...interface{}) (reply interface{}, err error) {
	args1 := mergeKeyAndArgs(key2, opts...)
	args2 := mergeKeyAndArgs(key1, args1...)
	args3 := mergeKeyAndArgs(numkeys, args2...)
	args := mergeKeyAndArgs(destination, args3...)

	return r.Do("ZINTERSTORE", args...)
}

// zrange key start stop [withscores]
func (r *Rds) ZRange(key, start, stop interface{}, opts ...interface{}) (reply interface{}, err error) {
	args1 := mergeKeyAndArgs(stop, opts...)
	args2 := mergeKeyAndArgs(start, args1...)
	args := mergeKeyAndArgs(key, args2...)
	return r.Do("ZRANGE", args...)
}

// zrangebyscore key min max [withscores] [limit offset count]
func (r *Rds) ZRangeByScore(key, min, max interface{}, opts ...interface{}) (reply interface{}, err error) {
	args1 := mergeKeyAndArgs(max, opts...)
	args2 := mergeKeyAndArgs(min, args1...)
	args := mergeKeyAndArgs(key, args2...)
	return r.Do("ZRANGEBYSCORE", args...)
}

func (r *Rds) ZRank(key, member interface{}) (reply interface{}, err error) {
	return r.Do("ZRANK", key, member)
}

func (r *Rds) ZRem(key, member interface{}, members ...interface{}) (reply interface{}, err error) {
	args1 := mergeKeyAndArgs(member, members...)
	args := mergeKeyAndArgs(key, args1...)
	return r.Do("ZREM", args...)
}

func (r *Rds) ZRemRangeByRank(key, start, stop interface{}) (reply interface{}, err error) {
	return r.Do("ZREMRANGEBYRANK", key, start, stop)
}

func (r *Rds) ZRemRangeByScore(key, min, max interface{}) (reply interface{}, err error) {
	return r.Do("ZREMRANGEBYSCORE", key, min, max)
}

func (r *Rds) ZRevRange(key, start, stop interface{}) (reply interface{}, err error) {
	return r.Do("ZREVRANGE", key, start, stop)
}

// zrevrangebyscore key max min [withscores] [limit offset count]
func (r *Rds) ZRevRangeByScore(key, max, min interface{}, opts ...interface{}) (reply interface{}, err error) {
	args1 := mergeKeyAndArgs(min, opts...)
	args2 := mergeKeyAndArgs(max, args1...)
	args := mergeKeyAndArgs(key, args2...)
	return r.Do("ZREVRANGEBYSCORE", args...)
}

func (r *Rds) ZRevRank(key, member interface{}) (reply interface{}, err error) {
	return r.Do("ZREVRANK", key, member)
}

func (r *Rds) ZScore(key, member interface{}) (reply interface{}, err error) {
	return r.Do("ZSCORE", key, member)
}

// zunionstore destination numkeys key [key ...] [weights weight] [sum|min|mix]
func (r *Rds) ZUnionStore(destination, numkeys, key1, key2 interface{}, opts ...interface{}) (reply interface{}, err error) {
	args1 := mergeKeyAndArgs(key2, opts...)
	args2 := mergeKeyAndArgs(key1, args1...)
	args3 := mergeKeyAndArgs(numkeys, args2...)
	args := mergeKeyAndArgs(destination, args3...)

	return r.Do("ZUNIONSTORE", args...)
}

func (r *Rds) ZScan(key, cursor interface{}, opts ...interface{}) (reply interface{}, err error) {
	args1 := mergeKeyAndArgs(cursor, opts...)
	args := mergeKeyAndArgs(key, args1...)

	return r.Do("ZSCAN", args...)
}

func mergeKeyAndArgs(key interface{}, args ...interface{}) []interface{} {
	args2 := make([]interface{}, 0)
	args2 = append(args2, key)
	if len(args) > 0 {
		for _, arg := range args {
			args2 = append(args2, arg)
		}
	}
	return args2
}
