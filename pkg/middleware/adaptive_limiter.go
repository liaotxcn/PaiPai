package middleware

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"math"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// 自适应限流模块/中间件
//  支持 Redis 全局限流 + 本地令牌桶兜底
//  根据 CPU 使用率、Redis 延迟实时调整 QPS（PID 控制器）
//  指标通过 go-zero/stat 上报

// ------------------------------------------------------------------
// 配置相关
// ------------------------------------------------------------------

// DynamicLimitConfig 限流器所有可调参数
type DynamicLimitConfig struct {
	BaseRate     int           // 初始 QPS
	MinRate      int           // 最低保护 QPS
	MaxRate      int           // 最高上限 QPS
	WindowSize   int           // 滑动窗口秒数（限流算法用）
	RedisKey     string        // Redis key 前缀
	CpuThreshold float64       // CPU 告警阈值，0~1
	RedisTimeout time.Duration // Redis ping 超时阈值
	BurstPercent int           // 突发流量 = rate * BurstPercent/100
}

// ------------------------------------------------------------------
// 自适应限流器主体
// ------------------------------------------------------------------

// AdaptiveLimiter 组合了 Redis + 本地令牌桶 + PID 调节
type AdaptiveLimiter struct {
	config        DynamicLimitConfig
	redisClient   *redis.Client
	metrics       *stat.Metrics     // go-zero 自带指标收集器
	cpuPercent    int32             // 原子：最近一次 CPU 使用率 %
	redisLatency  int32             // 原子：最近一次 Redis ping 延迟（纳秒）
	tokenBuckets  *collection.Cache // 本地令牌桶缓存（按 key）
	pidController *PIDController    // PID 控制器实例
	activeKeys    sync.Map          // 记录哪些 key 正在使用（优雅退出时清理）
	limiterLock   syncx.SpinLock    // 本地桶并发保护
}

// ------------------------------------------------------------------
// PID 控制器实现
// ------------------------------------------------------------------

// PIDController 简易位置式 PID
type PIDController struct {
	Kp, Ki, Kd float64 // 比例/积分/微分系数
	setpoint   float64 // 目标健康度（固定为 1.0）
	integral   float64 // 积分累计
	lastError  float64 // 上一次误差
}

// NewPIDController 创建一个 PID 控制器
func NewPIDController(Kp, Ki, Kd float64) *PIDController {
	return &PIDController{
		Kp: Kp,
		Ki: Ki,
		Kd: Kd,
	}
}

// Update 根据当前健康度输出调整系数（>0 放大 QPS，<0 缩小）
func (p *PIDController) Update(current float64) float64 {
	err := p.setpoint - current
	p.integral += err
	derivative := err - p.lastError
	out := p.Kp*err + p.Ki*p.integral + p.Kd*derivative
	p.lastError = err
	return out
}

// ------------------------------------------------------------------
// 本地令牌桶
// ------------------------------------------------------------------

// TokenBucket 单机内存令牌桶
type TokenBucket struct {
	rate       int        // 每秒产生的令牌数
	capacity   int        // 桶最大容量
	tokens     int        // 当前剩余令牌
	lastUpdate time.Time  // 上次更新时间
	mu         sync.Mutex // 并发锁
}

// NewTokenBucket 创建令牌桶实例
func NewTokenBucket(rate, burst int) *TokenBucket {
	return &TokenBucket{
		rate:       rate,
		capacity:   burst,
		tokens:     burst,
		lastUpdate: time.Now(),
	}
}

// Allow 判断本次请求能否拿到令牌
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastUpdate)
	tb.lastUpdate = now

	// 根据时间差补充令牌
	newTokens := int(float64(elapsed) / float64(time.Second) * float64(tb.rate))
	tb.tokens += newTokens
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// ------------------------------------------------------------------
// 构造器与生命周期
// ------------------------------------------------------------------

// NewAdaptiveLimiter 创建并初始化自适应限流器
func NewAdaptiveLimiter(c DynamicLimitConfig, redisClient *redis.Client) *AdaptiveLimiter {
	// 防御式默认值
	if c.MinRate <= 0 {
		c.MinRate = 100
	}
	if c.WindowSize <= 0 {
		c.WindowSize = 10
	}

	l := &AdaptiveLimiter{
		config:        c,
		redisClient:   redisClient,
		metrics:       stat.NewMetrics("limiter"),
		pidController: NewPIDController(0.8, 0.2, 0.1),
	}

	// 初始化本地缓存：1 分钟过期，最多 1 万个 key
	l.tokenBuckets, _ = collection.NewCache(time.Minute, collection.WithLimit(10000))

	// 后台协程：定期采集系统指标
	go l.collectMetrics()

	// 优雅退出：清理资源
	proc.AddShutdownListener(l.cleanup)

	return l
}

// collectMetrics 每 5 秒采样 CPU 与 Redis 延迟
func (l *AdaptiveLimiter) collectMetrics() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// CPU 使用率
		cpu := int32(math.Round(float64(stat.CpuUsage() * 100)))
		atomic.StoreInt32(&l.cpuPercent, cpu)

		// Redis ping 延迟
		ctx, cancel := context.WithTimeout(context.Background(), l.config.RedisTimeout)
		start := time.Now()
		if err := l.redisClient.Ping(ctx).Err(); err == nil {
			atomic.StoreInt32(&l.redisLatency, int32(time.Since(start).Nanoseconds()))
		}
		cancel()
	}
}

// cleanup 程序退出时统一清理 Redis key 与本地缓存
func (l *AdaptiveLimiter) cleanup() {
	var count int
	l.activeKeys.Range(func(key, _ interface{}) bool {
		l.tokenBuckets.Del(key.(string))
		count++
		return true
	})
	logx.Infof("清理令牌桶数量: %d", count)
}

// ------------------------------------------------------------------
// 动态速率计算
// ------------------------------------------------------------------

// currentRate 根据系统健康度实时计算当前 QPS 上限
func (l *AdaptiveLimiter) currentRate() int {
	cpuLoad := float64(atomic.LoadInt32(&l.cpuPercent)) / 100
	redisLatency := time.Duration(atomic.LoadInt32(&l.redisLatency))

	// 健康度：越接近 1 越健康
	health := 1.0 - math.Max(
		cpuLoad/l.config.CpuThreshold,
		float64(redisLatency)/float64(l.config.RedisTimeout),
	)

	// PID 输出：>0 表示扩大，<0 表示缩小
	adjustment := l.pidController.Update(health)
	rate := float64(l.config.BaseRate) * (1 + adjustment)

	// 边界保护
	rate = math.Min(rate, float64(l.config.MaxRate))
	rate = math.Max(rate, float64(l.config.MinRate))

	return int(rate)
}

// ------------------------------------------------------------------
// key 生成与中间件入口
// ------------------------------------------------------------------

// limitKey 把请求特征变成唯一字符串，用于 Redis 与本地桶
// 格式：prefix:path:ip:uid
func (l *AdaptiveLimiter) limitKey(r *http.Request) string {
	ip := httpx.GetRemoteAddr(r)
	uid, _ := r.Context().Value("uid").(string)
	if uid == "" {
		uid = "unknown"
	}
	key := fmt.Sprintf("%s:%s:%s:%s", l.config.RedisKey, r.URL.Path, ip, uid)
	l.activeKeys.Store(key, struct{}{}) // 记录活跃 key，便于清理
	return key
}

// Handle 返回一个标准 http.HandlerFunc，用于 go-zero 路由注册
func (l *AdaptiveLimiter) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 实时计算速率与突发容量
		rate := l.currentRate()
		burst := rate * l.config.BurstPercent / 100
		key := l.limitKey(r)

		// 1. 先走 Redis 分布式限流
		allowed, err := l.allowRedis(key, rate, burst)
		if err != nil {
			logx.Errorf("Redis限流失败: %v", err)
			// Redis 故障时退化到本地令牌桶
			allowed = l.allowLocal(key, rate, burst)
		}

		// 2. 拒绝请求
		if !allowed {
			l.metrics.Add(stat.Task{
				Drop:        true,
				Duration:    time.Since(start),
				Description: "rejected",
			})
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"code": 429, "msg": "请求过于频繁"}`))
			return
		}

		// 3. 通过请求
		l.metrics.Add(stat.Task{
			Drop:        false,
			Duration:    time.Since(start),
			Description: "passed",
		})
		next(w, r)
	}
}

// ------------------------------------------------------------------
// 限流算法实现
// ------------------------------------------------------------------

// allowRedis 使用 Lua 脚本在 Redis 中做令牌桶限流
// KEYS[1]  -> key
// ARGV[1]  -> rate
// ARGV[2]  -> burst
// ARGV[3]  -> now (unix)
// ARGV[4]  -> window size
func (l *AdaptiveLimiter) allowRedis(key string, rate, burst int) (bool, error) {
	script := `
	local key = KEYS[1]
	local rate = tonumber(ARGV[1])
	local burst = tonumber(ARGV[2])
	local now  = tonumber(ARGV[3])
	local window = tonumber(ARGV[4])

	local tokens = redis.call("HGET", key, "tokens") or burst
	local last_time = redis.call("HGET", key, "last_time") or now

	local elapsed = now - last_time
	local newTokens = elapsed * rate / window
	tokens = math.min(tokens + newTokens, burst)

	if tokens < 1 then
		return 0
	end

	tokens = tokens - 1
	redis.call("HSET", key, "tokens", tokens, "last_time", now)
	redis.call("EXPIRE", key, window * 2)
	return 1
	`
	now := time.Now().Unix()
	res, err := l.redisClient.Eval(
		context.Background(),
		script,
		[]string{key},
		rate, burst, now, l.config.WindowSize,
	).Int()
	return res == 1, err
}

// allowLocal 本地内存令牌桶兜底
func (l *AdaptiveLimiter) allowLocal(key string, rate, burst int) bool {
	l.limiterLock.Lock()
	defer l.limiterLock.Unlock()

	val, ok := l.tokenBuckets.Get(key)
	if !ok {
		val = NewTokenBucket(rate, burst)
		l.tokenBuckets.Set(key, val)
	}
	return val.(*TokenBucket).Allow()
}
