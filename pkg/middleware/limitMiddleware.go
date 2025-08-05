package middleware

import (
	"github.com/zeromicro/go-zero/core/limit" // limit限流库
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

// limit限流

type LimitMiddleware struct {
	redisCfg redis.RedisConf
	*limit.TokenLimiter
}

func (m *LimitMiddleware) TokenLimitHandler(rate, burst int) rest.Middleware {
	m.TokenLimiter = limit.NewTokenLimiter(rate, burst, redis.MustNewRedis(m.redisCfg),
		"REDIS_TOKEN_LIMIT_KEY")

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if m.TokenLimiter.AllowCtx(r.Context()) {
				next(w, r)
				return
			}
			w.WriteHeader(http.StatusTooManyRequests)
		}
	}

}
