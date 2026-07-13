package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

// 并发控制缓存常量定义
//
// 性能优化说明：
// 原实现使用 SCAN 命令遍历独立的槽位键（concurrency:account:{id}:{requestID}），
// 在高并发场景下 SCAN 需要多次往返，且遍历大量键时性能下降明显。
//
// 新实现改用 Redis 有序集合（Sorted Set）：
// 1. 每个账号/用户只有一个键，成员为 requestID，分数为时间戳
// 2. 使用 ZCARD 原子获取并发数，时间复杂度 O(1)
// 3. 使用 ZREMRANGEBYSCORE 清理过期槽位，避免手动管理 TTL
// 4. 单次 Redis 调用完成计数，减少网络往返
const (
	// 并发槽位键前缀（有序集合）
	// 格式: concurrency:account:{accountID}
	accountSlotKeyPrefix = "concurrency:account:"
	// 格式: concurrency:user:{userID}
	userSlotKeyPrefix = "concurrency:user:"
	// 格式: concurrency:api_key:{apiKeyID}
	apiKeySlotKeyPrefix = "concurrency:api_key:"
	// 共享容量组键使用相同 hash tag，确保 Redis Cluster 下可由单个 Lua 原子读写。
	// 格式:
	// concurrency:capacity:{set:<setID>}:group
	// concurrency:capacity:{set:<setID>}:member:<accountID>
	// concurrency:capacity:{set:<setID>}:metrics
	capacitySlotKeyPrefix = "concurrency:capacity:"
	// 等待队列计数器格式: concurrency:wait:{userID}
	waitQueueKeyPrefix = "concurrency:wait:"
	// 账号级等待队列计数器格式: wait:account:{accountID}
	accountWaitKeyPrefix = "wait:account:"

	// 默认槽位过期时间（分钟），可通过配置覆盖
	defaultSlotTTLMinutes = 15
	capacityMetricsTTL    = 48 * 60 * 60
)

var (
	// acquireScript 使用有序集合计数并在未达上限时添加槽位
	// 使用 Redis TIME 命令获取服务器时间，避免多实例时钟不同步问题
	// KEYS[1] = 有序集合键 (concurrency:account:{id} / concurrency:user:{id})
	// ARGV[1] = maxConcurrency
	// ARGV[2] = TTL（秒）
	// ARGV[3] = requestID
	acquireScript = redis.NewScript(`
		-- Redis 3.2-4.x compat: opt into effects replication so redis.call('TIME')
		-- replicates correctly. No-op on Redis 5.0+ (effects replication is default).
		redis.replicate_commands()
		local key = KEYS[1]
		local maxConcurrency = tonumber(ARGV[1])
		local ttl = tonumber(ARGV[2])
		local requestID = ARGV[3]

		-- 使用 Redis 服务器时间，确保多实例时钟一致
		local timeResult = redis.call('TIME')
		local now = tonumber(timeResult[1])
		local expireBefore = now - ttl

		-- 清理过期槽位
		redis.call('ZREMRANGEBYSCORE', key, '-inf', expireBefore)

		-- 检查是否已存在（支持重试场景刷新时间戳）
		local exists = redis.call('ZSCORE', key, requestID)
		if exists ~= false then
			redis.call('ZADD', key, now, requestID)
			redis.call('EXPIRE', key, ttl)
			return 1
		end

		-- 检查是否达到并发上限
		local count = redis.call('ZCARD', key)
		if count < maxConcurrency then
			redis.call('ZADD', key, now, requestID)
			redis.call('EXPIRE', key, ttl)
			return 1
		end

		return 0
	`)

	capacityAcquireScript = redis.NewScript(`
		redis.replicate_commands()
		local groupKey = KEYS[1]
		local memberKey = KEYS[2]
		local metricsKey = KEYS[3]
		local groupLimit = tonumber(ARGV[1])
		local memberHardLimit = tonumber(ARGV[2])
		local ttl = tonumber(ARGV[3])
		local requestID = ARGV[4]
		local memberSoftShare = tonumber(ARGV[5])
		local metricsTTL = tonumber(ARGV[6])

		local timeResult = redis.call('TIME')
		local now = tonumber(timeResult[1])
		local expireBefore = now - ttl
		local minuteBucket = math.floor(now / 60)

		redis.call('ZREMRANGEBYSCORE', groupKey, '-inf', expireBefore)
		redis.call('ZREMRANGEBYSCORE', memberKey, '-inf', expireBefore)

		local groupExists = redis.call('ZSCORE', groupKey, requestID)
		local memberExists = redis.call('ZSCORE', memberKey, requestID)
		if groupExists ~= false and memberExists ~= false then
			redis.call('ZADD', groupKey, now, requestID)
			redis.call('ZADD', memberKey, now, requestID)
			redis.call('EXPIRE', groupKey, ttl)
			redis.call('EXPIRE', memberKey, ttl)
			return {1, 0}
		end

		if groupExists ~= false then
			redis.call('ZREM', groupKey, requestID)
		end
		if memberExists ~= false then
			redis.call('ZREM', memberKey, requestID)
		end

		local groupCount = redis.call('ZCARD', groupKey)
		if groupCount >= groupLimit then
			redis.call('HINCRBY', metricsKey, 'group_full:' .. minuteBucket, 1)
			redis.call('EXPIRE', metricsKey, metricsTTL)
			return {2, 0}
		end

		local memberCount = redis.call('ZCARD', memberKey)
		if memberHardLimit > 0 and memberCount >= memberHardLimit then
			redis.call('HINCRBY', metricsKey, 'member_full:' .. minuteBucket, 1)
			redis.call('EXPIRE', metricsKey, metricsTTL)
			return {3, 0}
		end

		local borrowed = 0
		if memberSoftShare > 0 and memberCount >= memberSoftShare then
			borrowed = 1
			redis.call('HINCRBY', metricsKey, 'borrowed_slot:' .. minuteBucket, 1)
			redis.call('EXPIRE', metricsKey, metricsTTL)
		end

		redis.call('ZADD', groupKey, now, requestID)
		redis.call('ZADD', memberKey, now, requestID)
		redis.call('EXPIRE', groupKey, ttl)
		redis.call('EXPIRE', memberKey, ttl)
		return {1, borrowed}
	`)

	capacityReleaseScript = redis.NewScript(`
		redis.call('ZREM', KEYS[1], ARGV[1])
		redis.call('ZREM', KEYS[2], ARGV[1])
		return 1
	`)

	capacityCountScript = redis.NewScript(`
		redis.replicate_commands()
		local timeResult = redis.call('TIME')
		local now = tonumber(timeResult[1])
		local ttl = tonumber(ARGV[1])
		local expireBefore = now - ttl
		redis.call('ZREMRANGEBYSCORE', KEYS[1], '-inf', expireBefore)
		redis.call('ZREMRANGEBYSCORE', KEYS[2], '-inf', expireBefore)
		return {redis.call('ZCARD', KEYS[1]), redis.call('ZCARD', KEYS[2])}
	`)

	// getCountScript 统计有序集合中的槽位数量并清理过期条目
	// 使用 Redis TIME 命令获取服务器时间
	// KEYS[1] = 有序集合键
	// ARGV[1] = TTL（秒）
	getCountScript = redis.NewScript(`
		-- Redis 3.2-4.x compat: opt into effects replication so redis.call('TIME')
		-- replicates correctly. No-op on Redis 5.0+ (effects replication is default).
		redis.replicate_commands()
		local key = KEYS[1]
		local ttl = tonumber(ARGV[1])

		-- 使用 Redis 服务器时间
		local timeResult = redis.call('TIME')
		local now = tonumber(timeResult[1])
		local expireBefore = now - ttl

		redis.call('ZREMRANGEBYSCORE', key, '-inf', expireBefore)
		return redis.call('ZCARD', key)
	`)

	// trackSlotScript 记录 stats-only 槽位，不做并发上限判断。
	// KEYS[1] = 有序集合键
	// ARGV[1] = TTL（秒）
	// ARGV[2] = requestID
	trackSlotScript = redis.NewScript(`
		-- Redis 3.2-4.x compat: opt into effects replication so redis.call('TIME')
		-- replicates correctly. No-op on Redis 5.0+ (effects replication is default).
		redis.replicate_commands()
		local key = KEYS[1]
		local ttl = tonumber(ARGV[1])
		local requestID = ARGV[2]

		local timeResult = redis.call('TIME')
		local now = tonumber(timeResult[1])
		local expireBefore = now - ttl

		redis.call('ZREMRANGEBYSCORE', key, '-inf', expireBefore)
		redis.call('ZADD', key, now, requestID)
		redis.call('EXPIRE', key, ttl)
		return 1
	`)

	// incrementWaitScript - refreshes TTL on each increment to keep queue depth accurate
	// KEYS[1] = wait queue key
	// ARGV[1] = maxWait
	// ARGV[2] = TTL in seconds
	incrementWaitScript = redis.NewScript(`
		local current = redis.call('GET', KEYS[1])
		if current == false then
			current = 0
		else
			current = tonumber(current)
		end

		if current >= tonumber(ARGV[1]) then
			return 0
		end

		local newVal = redis.call('INCR', KEYS[1])

		-- Refresh TTL so long-running traffic doesn't expire active queue counters.
		redis.call('EXPIRE', KEYS[1], ARGV[2])

			return 1
		`)

	// incrementAccountWaitScript - account-level wait queue count (refresh TTL on each increment)
	incrementAccountWaitScript = redis.NewScript(`
			local current = redis.call('GET', KEYS[1])
			if current == false then
				current = 0
			else
				current = tonumber(current)
			end

			if current >= tonumber(ARGV[1]) then
				return 0
			end

			local newVal = redis.call('INCR', KEYS[1])

			-- Refresh TTL so long-running traffic doesn't expire active queue counters.
			redis.call('EXPIRE', KEYS[1], ARGV[2])

			return 1
		`)

	// decrementWaitScript - same as before
	decrementWaitScript = redis.NewScript(`
			local current = redis.call('GET', KEYS[1])
			if current ~= false and tonumber(current) > 0 then
				redis.call('DECR', KEYS[1])
			end
			return 1
		`)

	// cleanupExpiredSlotsScript 清理单个账号/用户有序集合中过期槽位
	// KEYS[1] = 有序集合键
	// ARGV[1] = TTL（秒）
	cleanupExpiredSlotsScript = redis.NewScript(`
		-- Redis 3.2-4.x compat: opt into effects replication so redis.call('TIME')
		-- replicates correctly. No-op on Redis 5.0+ (effects replication is default).
		redis.replicate_commands()
		local key = KEYS[1]
		local ttl = tonumber(ARGV[1])
		local timeResult = redis.call('TIME')
		local now = tonumber(timeResult[1])
		local expireBefore = now - ttl
		redis.call('ZREMRANGEBYSCORE', key, '-inf', expireBefore)
		if redis.call('ZCARD', key) == 0 then
			redis.call('DEL', key)
		else
			redis.call('EXPIRE', key, ttl)
		end
		return 1
	`)

	// startupCleanupScript 清理非当前进程前缀的槽位成员。
	// KEYS 是有序集合键列表，ARGV[1] 是当前进程前缀，ARGV[2] 是槽位 TTL。
	// 遍历每个 KEYS[i]，移除前缀不匹配的成员，清空后删 key，否则刷新 EXPIRE。
	startupCleanupScript = redis.NewScript(`
		local activePrefix = ARGV[1]
		local slotTTL = tonumber(ARGV[2])
		local removed = 0
		for i = 1, #KEYS do
			local key = KEYS[i]
			local members = redis.call('ZRANGE', key, 0, -1)
			for _, member in ipairs(members) do
				if string.sub(member, 1, string.len(activePrefix)) ~= activePrefix then
					removed = removed + redis.call('ZREM', key, member)
				end
			end
			if redis.call('ZCARD', key) == 0 then
				redis.call('DEL', key)
			else
				redis.call('EXPIRE', key, slotTTL)
			end
		end
		return removed
	`)
)

type concurrencyCache struct {
	rdb                 *redis.Client
	slotTTLSeconds      int // 槽位过期时间（秒）
	waitQueueTTLSeconds int // 等待队列过期时间（秒）
}

// NewConcurrencyCache 创建并发控制缓存
// slotTTLMinutes: 槽位过期时间（分钟），0 或负数使用默认值 15 分钟
// waitQueueTTLSeconds: 等待队列过期时间（秒），0 或负数使用 slot TTL
func NewConcurrencyCache(rdb *redis.Client, slotTTLMinutes int, waitQueueTTLSeconds int) service.ConcurrencyCache {
	if slotTTLMinutes <= 0 {
		slotTTLMinutes = defaultSlotTTLMinutes
	}
	if waitQueueTTLSeconds <= 0 {
		waitQueueTTLSeconds = slotTTLMinutes * 60
	}
	return &concurrencyCache{
		rdb:                 rdb,
		slotTTLSeconds:      slotTTLMinutes * 60,
		waitQueueTTLSeconds: waitQueueTTLSeconds,
	}
}

// Helper functions for key generation
func accountSlotKey(accountID int64) string {
	return fmt.Sprintf("%s%d", accountSlotKeyPrefix, accountID)
}

func userSlotKey(userID int64) string {
	return fmt.Sprintf("%s%d", userSlotKeyPrefix, userID)
}

func apiKeySlotKey(apiKeyID int64) string {
	return fmt.Sprintf("%s%d", apiKeySlotKeyPrefix, apiKeyID)
}

func capacityGroupHashTag(groupID int64) string {
	return fmt.Sprintf("{set:%d}", groupID)
}

func capacityGroupSlotKey(groupID int64) string {
	return capacitySlotKeyPrefix + capacityGroupHashTag(groupID) + ":group"
}

func capacityMemberSlotKey(groupID, accountID int64) string {
	return fmt.Sprintf("%s%s:member:%d", capacitySlotKeyPrefix, capacityGroupHashTag(groupID), accountID)
}

func capacityMetricsKey(groupID int64) string {
	return capacitySlotKeyPrefix + capacityGroupHashTag(groupID) + ":metrics"
}

func waitQueueKey(userID int64) string {
	return fmt.Sprintf("%s%d", waitQueueKeyPrefix, userID)
}

func accountWaitKey(accountID int64) string {
	return fmt.Sprintf("%s%d", accountWaitKeyPrefix, accountID)
}

// Account slot operations

func (c *concurrencyCache) AcquireAccountSlot(ctx context.Context, accountID int64, maxConcurrency int, requestID string) (bool, error) {
	key := accountSlotKey(accountID)
	// 时间戳在 Lua 脚本内使用 Redis TIME 命令获取，确保多实例时钟一致
	result, err := acquireScript.Run(ctx, c.rdb, []string{key}, maxConcurrency, c.slotTTLSeconds, requestID).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func (c *concurrencyCache) ReleaseAccountSlot(ctx context.Context, accountID int64, requestID string) error {
	key := accountSlotKey(accountID)
	return c.rdb.ZRem(ctx, key, requestID).Err()
}

func (c *concurrencyCache) AcquireCapacitySlot(
	ctx context.Context,
	groupID int64,
	accountID int64,
	groupLimit int,
	memberHardLimit int,
	memberSoftShare int,
	requestID string,
) (service.CapacitySlotAcquireResult, error) {
	if groupID <= 0 || accountID <= 0 || groupLimit <= 0 || requestID == "" {
		return service.CapacitySlotAcquireResult{}, errors.New("invalid shared capacity slot input")
	}
	keys := []string{
		capacityGroupSlotKey(groupID),
		capacityMemberSlotKey(groupID, accountID),
		capacityMetricsKey(groupID),
	}
	result, err := capacityAcquireScript.Run(
		ctx,
		c.rdb,
		keys,
		groupLimit,
		memberHardLimit,
		c.slotTTLSeconds,
		requestID,
		memberSoftShare,
		capacityMetricsTTL,
	).Slice()
	if err != nil {
		return service.CapacitySlotAcquireResult{}, err
	}
	if len(result) != 2 {
		return service.CapacitySlotAcquireResult{}, fmt.Errorf("unexpected shared capacity acquire result length: %d", len(result))
	}
	status, err := redisResultInt(result[0])
	if err != nil {
		return service.CapacitySlotAcquireResult{}, fmt.Errorf("parse shared capacity acquire status: %w", err)
	}
	borrowed, err := redisResultInt(result[1])
	if err != nil {
		return service.CapacitySlotAcquireResult{}, fmt.Errorf("parse shared capacity borrowed flag: %w", err)
	}
	return service.CapacitySlotAcquireResult{
		Status:   service.CapacityAcquireStatus(status),
		Borrowed: borrowed == 1,
	}, nil
}

func (c *concurrencyCache) ReleaseCapacitySlot(ctx context.Context, groupID, accountID int64, requestID string) error {
	if groupID <= 0 || accountID <= 0 || requestID == "" {
		return nil
	}
	_, err := capacityReleaseScript.Run(
		ctx,
		c.rdb,
		[]string{capacityGroupSlotKey(groupID), capacityMemberSlotKey(groupID, accountID)},
		requestID,
	).Result()
	return err
}

func (c *concurrencyCache) GetCapacitySlotCounts(ctx context.Context, groupID, accountID int64) (service.CapacitySlotCounts, error) {
	if groupID <= 0 || accountID <= 0 {
		return service.CapacitySlotCounts{}, nil
	}
	result, err := capacityCountScript.Run(
		ctx,
		c.rdb,
		[]string{capacityGroupSlotKey(groupID), capacityMemberSlotKey(groupID, accountID)},
		c.slotTTLSeconds,
	).Slice()
	if err != nil {
		return service.CapacitySlotCounts{}, err
	}
	if len(result) != 2 {
		return service.CapacitySlotCounts{}, fmt.Errorf("unexpected shared capacity count result length: %d", len(result))
	}
	groupCount, err := redisResultInt(result[0])
	if err != nil {
		return service.CapacitySlotCounts{}, fmt.Errorf("parse shared capacity group count: %w", err)
	}
	memberCount, err := redisResultInt(result[1])
	if err != nil {
		return service.CapacitySlotCounts{}, fmt.Errorf("parse shared capacity member count: %w", err)
	}
	return service.CapacitySlotCounts{
		GroupConcurrency:  groupCount,
		MemberConcurrency: memberCount,
	}, nil
}

func (c *concurrencyCache) GetCapacityMetrics(ctx context.Context, groupID int64, since time.Time) (service.CapacityMetrics, error) {
	if groupID <= 0 {
		return service.CapacityMetrics{}, nil
	}
	values, err := c.rdb.HGetAll(ctx, capacityMetricsKey(groupID)).Result()
	if err != nil {
		return service.CapacityMetrics{}, err
	}
	minBucket := since.Unix() / 60
	metrics := service.CapacityMetrics{}
	for field, raw := range values {
		parts := strings.SplitN(field, ":", 2)
		if len(parts) != 2 {
			continue
		}
		bucket, parseErr := strconv.ParseInt(parts[1], 10, 64)
		if parseErr != nil || bucket < minBucket {
			continue
		}
		value, parseErr := strconv.Atoi(raw)
		if parseErr != nil {
			continue
		}
		switch parts[0] {
		case "group_full":
			metrics.GroupFullCount += value
		case "member_full":
			metrics.MemberFullCount += value
		case "borrowed_slot":
			metrics.BorrowedSlotCount += value
		}
	}
	return metrics, nil
}

func redisResultInt(value any) (int, error) {
	switch typed := value.(type) {
	case int64:
		return int(typed), nil
	case string:
		return strconv.Atoi(typed)
	case []byte:
		return strconv.Atoi(string(typed))
	default:
		return 0, fmt.Errorf("unsupported redis result type %T", value)
	}
}

func (c *concurrencyCache) GetAccountConcurrency(ctx context.Context, accountID int64) (int, error) {
	key := accountSlotKey(accountID)
	// 时间戳在 Lua 脚本内使用 Redis TIME 命令获取
	result, err := getCountScript.Run(ctx, c.rdb, []string{key}, c.slotTTLSeconds).Int()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c *concurrencyCache) GetAccountConcurrencyBatch(ctx context.Context, accountIDs []int64) (map[int64]int, error) {
	if len(accountIDs) == 0 {
		return map[int64]int{}, nil
	}

	now, err := c.rdb.Time(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis TIME: %w", err)
	}
	cutoffTime := now.Unix() - int64(c.slotTTLSeconds)

	pipe := c.rdb.Pipeline()
	type accountCmd struct {
		accountID int64
		zcardCmd  *redis.IntCmd
	}
	cmds := make([]accountCmd, 0, len(accountIDs))
	for _, accountID := range accountIDs {
		slotKey := accountSlotKeyPrefix + strconv.FormatInt(accountID, 10)
		pipe.ZRemRangeByScore(ctx, slotKey, "-inf", strconv.FormatInt(cutoffTime, 10))
		cmds = append(cmds, accountCmd{
			accountID: accountID,
			zcardCmd:  pipe.ZCard(ctx, slotKey),
		})
	}

	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("pipeline exec: %w", err)
	}

	result := make(map[int64]int, len(accountIDs))
	for _, cmd := range cmds {
		result[cmd.accountID] = int(cmd.zcardCmd.Val())
	}
	return result, nil
}

// User slot operations

func (c *concurrencyCache) AcquireUserSlot(ctx context.Context, userID int64, maxConcurrency int, requestID string) (bool, error) {
	key := userSlotKey(userID)
	// 时间戳在 Lua 脚本内使用 Redis TIME 命令获取，确保多实例时钟一致
	result, err := acquireScript.Run(ctx, c.rdb, []string{key}, maxConcurrency, c.slotTTLSeconds, requestID).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func (c *concurrencyCache) ReleaseUserSlot(ctx context.Context, userID int64, requestID string) error {
	key := userSlotKey(userID)
	return c.rdb.ZRem(ctx, key, requestID).Err()
}

func (c *concurrencyCache) GetUserConcurrency(ctx context.Context, userID int64) (int, error) {
	key := userSlotKey(userID)
	// 时间戳在 Lua 脚本内使用 Redis TIME 命令获取
	result, err := getCountScript.Run(ctx, c.rdb, []string{key}, c.slotTTLSeconds).Int()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c *concurrencyCache) TrackAPIKeySlot(ctx context.Context, apiKeyID int64, requestID string) error {
	key := apiKeySlotKey(apiKeyID)
	_, err := trackSlotScript.Run(ctx, c.rdb, []string{key}, c.slotTTLSeconds, requestID).Result()
	return err
}

func (c *concurrencyCache) ReleaseAPIKeySlot(ctx context.Context, apiKeyID int64, requestID string) error {
	key := apiKeySlotKey(apiKeyID)
	return c.rdb.ZRem(ctx, key, requestID).Err()
}

func (c *concurrencyCache) GetAPIKeyConcurrencyBatch(ctx context.Context, apiKeyIDs []int64) (map[int64]int, error) {
	if len(apiKeyIDs) == 0 {
		return map[int64]int{}, nil
	}

	now, err := c.rdb.Time(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis TIME: %w", err)
	}
	cutoffTime := now.Unix() - int64(c.slotTTLSeconds)

	pipe := c.rdb.Pipeline()
	type apiKeyCmd struct {
		apiKeyID int64
		zcardCmd *redis.IntCmd
	}
	cmds := make([]apiKeyCmd, 0, len(apiKeyIDs))
	for _, apiKeyID := range apiKeyIDs {
		slotKey := apiKeySlotKeyPrefix + strconv.FormatInt(apiKeyID, 10)
		pipe.ZRemRangeByScore(ctx, slotKey, "-inf", strconv.FormatInt(cutoffTime, 10))
		cmds = append(cmds, apiKeyCmd{
			apiKeyID: apiKeyID,
			zcardCmd: pipe.ZCard(ctx, slotKey),
		})
	}

	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("pipeline exec: %w", err)
	}

	result := make(map[int64]int, len(apiKeyIDs))
	for _, cmd := range cmds {
		result[cmd.apiKeyID] = int(cmd.zcardCmd.Val())
	}
	return result, nil
}

// Wait queue operations

func (c *concurrencyCache) IncrementWaitCount(ctx context.Context, userID int64, maxWait int) (bool, error) {
	key := waitQueueKey(userID)
	result, err := incrementWaitScript.Run(ctx, c.rdb, []string{key}, maxWait, c.waitQueueTTLSeconds).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func (c *concurrencyCache) DecrementWaitCount(ctx context.Context, userID int64) error {
	key := waitQueueKey(userID)
	_, err := decrementWaitScript.Run(ctx, c.rdb, []string{key}).Result()
	return err
}

// Account wait queue operations

func (c *concurrencyCache) IncrementAccountWaitCount(ctx context.Context, accountID int64, maxWait int) (bool, error) {
	key := accountWaitKey(accountID)
	result, err := incrementAccountWaitScript.Run(ctx, c.rdb, []string{key}, maxWait, c.waitQueueTTLSeconds).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func (c *concurrencyCache) DecrementAccountWaitCount(ctx context.Context, accountID int64) error {
	key := accountWaitKey(accountID)
	_, err := decrementWaitScript.Run(ctx, c.rdb, []string{key}).Result()
	return err
}

func (c *concurrencyCache) GetAccountWaitingCount(ctx context.Context, accountID int64) (int, error) {
	key := accountWaitKey(accountID)
	val, err := c.rdb.Get(ctx, key).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, err
	}
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return val, nil
}

func (c *concurrencyCache) GetAccountsLoadBatch(ctx context.Context, accounts []service.AccountWithConcurrency) (map[int64]*service.AccountLoadInfo, error) {
	if len(accounts) == 0 {
		return map[int64]*service.AccountLoadInfo{}, nil
	}

	// 使用 Pipeline 替代 Lua 脚本，兼容 Redis Cluster（Lua 内动态拼 key 会 CROSSSLOT）。
	// 每个账号执行 3 个命令：ZREMRANGEBYSCORE（清理过期）、ZCARD（并发数）、GET（等待数）。
	now, err := c.rdb.Time(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis TIME: %w", err)
	}
	cutoffTime := now.Unix() - int64(c.slotTTLSeconds)

	pipe := c.rdb.Pipeline()

	type accountCmds struct {
		id             int64
		maxConcurrency int
		zcardCmd       *redis.IntCmd
		groupZCardCmd  *redis.IntCmd
		capacity       *service.AccountCapacityScope
		getCmd         *redis.StringCmd
	}
	cmds := make([]accountCmds, 0, len(accounts))
	groupZCardCmds := make(map[int64]*redis.IntCmd)
	for _, acc := range accounts {
		slotKey := accountSlotKeyPrefix + strconv.FormatInt(acc.ID, 10)
		var groupZCardCmd *redis.IntCmd
		if acc.Capacity != nil && acc.Capacity.GroupID > 0 && acc.Capacity.GroupLimit > 0 {
			slotKey = capacityMemberSlotKey(acc.Capacity.GroupID, acc.ID)
			groupZCardCmd = groupZCardCmds[acc.Capacity.GroupID]
			if groupZCardCmd == nil {
				groupKey := capacityGroupSlotKey(acc.Capacity.GroupID)
				pipe.ZRemRangeByScore(ctx, groupKey, "-inf", strconv.FormatInt(cutoffTime, 10))
				groupZCardCmd = pipe.ZCard(ctx, groupKey)
				groupZCardCmds[acc.Capacity.GroupID] = groupZCardCmd
			}
		}
		waitKey := accountWaitKeyPrefix + strconv.FormatInt(acc.ID, 10)
		pipe.ZRemRangeByScore(ctx, slotKey, "-inf", strconv.FormatInt(cutoffTime, 10))
		ac := accountCmds{
			id:             acc.ID,
			maxConcurrency: acc.MaxConcurrency,
			zcardCmd:       pipe.ZCard(ctx, slotKey),
			groupZCardCmd:  groupZCardCmd,
			capacity:       acc.Capacity,
			getCmd:         pipe.Get(ctx, waitKey),
		}
		cmds = append(cmds, ac)
	}

	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("pipeline exec: %w", err)
	}

	loadMap := make(map[int64]*service.AccountLoadInfo, len(accounts))
	for _, ac := range cmds {
		currentConcurrency := int(ac.zcardCmd.Val())
		waitingCount := 0
		if v, err := ac.getCmd.Int(); err == nil {
			waitingCount = v
		}
		loadRate := 0
		memberLimit := ac.maxConcurrency
		if ac.capacity != nil {
			if ac.capacity.MemberHardLimit > 0 {
				memberLimit = ac.capacity.MemberHardLimit
			} else if ac.capacity.MemberSoftShare > 0 {
				memberLimit = ac.capacity.MemberSoftShare
			}
		}
		if memberLimit > 0 {
			loadRate = (currentConcurrency + waitingCount) * 100 / memberLimit
		}
		if ac.capacity != nil && ac.capacity.GroupLimit > 0 && ac.groupZCardCmd != nil {
			groupLoadRate := int(ac.groupZCardCmd.Val()) * 100 / ac.capacity.GroupLimit
			if groupLoadRate > loadRate {
				loadRate = groupLoadRate
			}
		}
		loadMap[ac.id] = &service.AccountLoadInfo{
			AccountID:          ac.id,
			CurrentConcurrency: currentConcurrency,
			WaitingCount:       waitingCount,
			LoadRate:           loadRate,
		}
	}

	return loadMap, nil
}

func (c *concurrencyCache) GetUsersLoadBatch(ctx context.Context, users []service.UserWithConcurrency) (map[int64]*service.UserLoadInfo, error) {
	if len(users) == 0 {
		return map[int64]*service.UserLoadInfo{}, nil
	}

	// 使用 Pipeline 替代 Lua 脚本，兼容 Redis Cluster。
	now, err := c.rdb.Time(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis TIME: %w", err)
	}
	cutoffTime := now.Unix() - int64(c.slotTTLSeconds)

	pipe := c.rdb.Pipeline()

	type userCmds struct {
		id             int64
		maxConcurrency int
		zcardCmd       *redis.IntCmd
		getCmd         *redis.StringCmd
	}
	cmds := make([]userCmds, 0, len(users))
	for _, u := range users {
		slotKey := userSlotKeyPrefix + strconv.FormatInt(u.ID, 10)
		waitKey := waitQueueKeyPrefix + strconv.FormatInt(u.ID, 10)
		pipe.ZRemRangeByScore(ctx, slotKey, "-inf", strconv.FormatInt(cutoffTime, 10))
		uc := userCmds{
			id:             u.ID,
			maxConcurrency: u.MaxConcurrency,
			zcardCmd:       pipe.ZCard(ctx, slotKey),
			getCmd:         pipe.Get(ctx, waitKey),
		}
		cmds = append(cmds, uc)
	}

	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("pipeline exec: %w", err)
	}

	loadMap := make(map[int64]*service.UserLoadInfo, len(users))
	for _, uc := range cmds {
		currentConcurrency := int(uc.zcardCmd.Val())
		waitingCount := 0
		if v, err := uc.getCmd.Int(); err == nil {
			waitingCount = v
		}
		loadRate := 0
		if uc.maxConcurrency > 0 {
			loadRate = (currentConcurrency + waitingCount) * 100 / uc.maxConcurrency
		}
		loadMap[uc.id] = &service.UserLoadInfo{
			UserID:             uc.id,
			CurrentConcurrency: currentConcurrency,
			WaitingCount:       waitingCount,
			LoadRate:           loadRate,
		}
	}

	return loadMap, nil
}

func (c *concurrencyCache) CleanupExpiredAccountSlots(ctx context.Context, accountID int64) error {
	key := accountSlotKey(accountID)
	_, err := cleanupExpiredSlotsScript.Run(ctx, c.rdb, []string{key}, c.slotTTLSeconds).Result()
	return err
}

func (c *concurrencyCache) CleanupStaleProcessSlots(ctx context.Context, activeRequestPrefix string) error {
	if activeRequestPrefix == "" {
		return nil
	}

	// 1. 清理有序集合中非当前进程前缀的成员
	slotPatterns := []string{accountSlotKeyPrefix + "*", userSlotKeyPrefix + "*", apiKeySlotKeyPrefix + "*"}
	for _, pattern := range slotPatterns {
		if err := c.cleanupSlotsByPattern(ctx, pattern, activeRequestPrefix); err != nil {
			return err
		}
	}
	capacitySlotPatterns := []string{
		capacitySlotKeyPrefix + "*:group",
		capacitySlotKeyPrefix + "*:member:*",
	}
	for _, pattern := range capacitySlotPatterns {
		if err := c.cleanupSlotsByPatternOneByOne(ctx, pattern, activeRequestPrefix); err != nil {
			return err
		}
	}

	// 2. 删除所有等待队列计数器（重启后计数器失效）
	waitPatterns := []string{accountWaitKeyPrefix + "*", waitQueueKeyPrefix + "*"}
	for _, pattern := range waitPatterns {
		if err := c.deleteKeysByPattern(ctx, pattern); err != nil {
			return err
		}
	}

	return nil
}

func (c *concurrencyCache) cleanupSlotsByPatternOneByOne(ctx context.Context, pattern, activePrefix string) error {
	const scanCount = 200
	var cursor uint64
	for {
		keys, nextCursor, err := c.rdb.Scan(ctx, cursor, pattern, scanCount).Result()
		if err != nil {
			return fmt.Errorf("scan %s: %w", pattern, err)
		}
		for _, key := range keys {
			if _, err := startupCleanupScript.Run(ctx, c.rdb, []string{key}, activePrefix, c.slotTTLSeconds).Result(); err != nil {
				return fmt.Errorf("cleanup slot %s: %w", key, err)
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			return nil
		}
	}
}

// cleanupSlotsByPattern 扫描匹配 pattern 的有序集合键，批量调用 Lua 脚本清理非当前进程成员。
func (c *concurrencyCache) cleanupSlotsByPattern(ctx context.Context, pattern, activePrefix string) error {
	const scanCount = 200
	var cursor uint64
	for {
		keys, nextCursor, err := c.rdb.Scan(ctx, cursor, pattern, scanCount).Result()
		if err != nil {
			return fmt.Errorf("scan %s: %w", pattern, err)
		}
		if len(keys) > 0 {
			_, err := startupCleanupScript.Run(ctx, c.rdb, keys, activePrefix, c.slotTTLSeconds).Result()
			if err != nil {
				return fmt.Errorf("cleanup slots %s: %w", pattern, err)
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

// deleteKeysByPattern 扫描匹配 pattern 的键并删除。
func (c *concurrencyCache) deleteKeysByPattern(ctx context.Context, pattern string) error {
	const scanCount = 200
	var cursor uint64
	for {
		keys, nextCursor, err := c.rdb.Scan(ctx, cursor, pattern, scanCount).Result()
		if err != nil {
			return fmt.Errorf("scan %s: %w", pattern, err)
		}
		if len(keys) > 0 {
			if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("del %s: %w", pattern, err)
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}
