package rateLimit

// 这是固定窗口的限流实现
// 栗子 对ip做限流 3秒内只能访问100次
// 固定窗口就是 第1秒访问了1次，第三秒访问了99次，那么第3秒仍然可访问100次
func fixedWindowScript() string {
	script := `
		local expireTime = ARGV[1] 
		local limitNum = ARGV[2]
		local key = KEYS[1]

		local visitNum = redis.call('incr', key)
		if visitNum == 1 then
				redis.call('expire', key, expireTime)
		end
		
		if visitNum > tonumber(limitNum) then
				return 0
		end
		
		return 1;
    `
	return script
}

// 滑动窗口的实现
// 借助于 sorted set 实现
// 有序 set 中，可以根据 score进行排序
// 将 score 中存放时间戳，当前时间戳到（当前时间戳-过期时间）之间的范围即为滑动窗口的正常范围
func slidingWindowScript() string {
	script := `
		-- 过期时间，单位秒
		local expireTime = ARGV[1]
		-- 限制的最大数量
		local maxThreads = ARGV[2]
		-- 排序的分数
		local score = ARGV[3]
		-- 添加的成员
		local randomValue = ARGV[4]

		-- 有序集合的key
		local key = KEYS[1]
		-- 当前计算开始的分数
		local beginScore = tonumber(score)-tonumber(expireTime)

		local visitNum = redis.call('ZCOUNT', key, beginScore, tonumber(score))
		if visitNum >= tonumber(maxThreads) then
			return 0
		end

		redis.call('ZADD', key, score, randomValue)

		-- 设置过期时间
		if visitNum == 0 then
			redis.call('EXPIRE', key, expireTime)
		end
		-- 删除不在范围内的成员
		if visitNum > 1 then
			redis.call('ZREMRANGEBYSCORE', key, 1, beginScore)
		end

		return 1;
    `
	return script
}

// 令牌桶算法的实现
// 每次访问的时候，会用（当前时间 -  最后一次访问的时间）* 速率 计算下当前应该下发的令牌数量，作为目前令牌桶中的令牌数
// 如果当前令牌桶中的令牌数为0，表示令牌不够，需要进行限流了
// 如果当前的令牌有剩余，下发令牌，然后更新最后一次访问的时间和当前的剩余令牌数
func tokenBucketScript() string {
	script := `
		-- 令牌桶中的key
		local key = KEYS[1]

		-- 令牌的速率 个数/s
		local rate = tonumber(ARGV[1])
		-- 限制的最大数量
		local maxThreads = ARGV[2]
		-- 当前的时间
		local timeNow = tonumber(ARGV[3])

		local rateLimitInfo = redis.pcall("HMGET", key, "lastMillSecond", "currPermits")

		-- 上次添加令牌的时间
		local lastMillSecond = tonumber(rateLimitInfo[1])
		-- 桶里当前令牌数
		local currPermits = tonumber(rateLimitInfo[2])

		local isFirst = false
		if lastMillSecond == nil then
			lastMillSecond = timeNow
			currPermits = maxThreads
			isFirst = true
		end

		-- 每次计算下当前可以下发令牌的数量
		local addBucket = (timeNow - lastMillSecond) * rate
		currPermits = currPermits + addBucket

		if currPermits == 0 then
			redis.pcall("HSET", key, "lastMillSecond", timeNow)
			return 0
		end

		redis.pcall("HSET", key, "lastMillSecond", timeNow, "currPermits",currPermits-1)  

		if isFirst == true then
			-- 防止 key 之后不用，设置一个过期的时间，根据业务决定
			redis.call('EXPIRE', key, 60*60*24)
		end

		return 1;
    `
	return script
}
