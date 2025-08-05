package interceptor

import (
	"PaiPai/pkg/wuid"
	"PaiPai/pkg/xerr"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type Idempotent interface {
	// Identify 获取请求的标识
	Identify(ctx context.Context, method string) string
	// IsIdempotentMethod 是够支持幂等性
	IsIdempotentMethod(fullMethod string) bool
	// TryAcquire 幂等性的验证
	TryAcquire(ctx context.Context, id string) (resp any, isAcquire bool)
	// SaveResp 执行之后结果的保存
	SaveResp(ctx context.Context, id string, resp any, respErr error) error
}

var (
	// TKey 请求任务标识
	TKey = "pai-pai-idempotence-task-id"
	// DKey 设置rpc调度中的rpc请求的标识
	DKey = "pai-pai-idempotence-dispatch-key"
)

//  设值客户端、服务端拦截器：

// ContextWithVal 添加到上下文方便客户端获取
func ContextWithVal(ctx context.Context) context.Context {
	return context.WithValue(ctx, TKey, wuid.GenerateID())
}

func NewIdempotenceClient(idempotent Idempotent) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 获取唯一的key
		identify := idempotent.Identify(ctx, method)
		// 在rpc请求中防止头部信息
		ctx = metadata.NewOutgoingContext(ctx, map[string][]string{
			DKey: []string{identify},
		})
		// 请求
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func NewIdempotenceServer(idempotent Idempotent) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		identify := metadata.ValueFromIncomingContext(ctx, DKey)
		if len(identify) == 0 || !idempotent.IsIdempotentMethod(info.FullMethod) {
			// 不进行幂等处理
			return handler(ctx, req)
		}

		fmt.Println("----", "请求进行幂等处理", identify)
		r, isAcquire := idempotent.TryAcquire(ctx, identify[0])
		if isAcquire {
			resp, err = handler(ctx, req)
			fmt.Println("---- 执行任务")
			if err := idempotent.SaveResp(ctx, identify[0], resp, err); err != nil {
				return resp, err
			}
			return resp, err
		}
		// 任务在执行或者已经执行完了
		if r != nil {
			fmt.Println("----", "任务已经执行完了")
			return r, nil
		}
		return nil, errors.WithStack(xerr.New(int(codes.DeadlineExceeded), fmt.Sprintf("存在其他任务在执行"+
			"id %v", identify[0])))
	}
}

var (
	DefaultIdempotent       = new(defaultIdempotent)
	DefaultIdempotentClient = NewIdempotenceClient(DefaultIdempotent)
)

type defaultIdempotent struct {
	// 获取和设置请求的id
	Redis *redis.Redis
	// 注意存储
	Cache *collection.Cache
	// 设置方法时对幂等的支持
	method map[string]bool
}

func NewDefaultIdempotent(c redis.RedisConf) Idempotent {
	cache, err := collection.NewCache(60 * 60)
	if err != nil {
		panic(err)
	}
	return &defaultIdempotent{
		Redis: redis.MustNewRedis(c),
		Cache: cache,
		method: map[string]bool{
			"/social.social/GroupCreate": true,
		},
	}
}

func (d *defaultIdempotent) Identify(ctx context.Context, method string) string {
	id := ctx.Value(TKey)
	// 让其生成请求id
	rpcId := fmt.Sprintf("%v.%s", id, method)
	return rpcId
}

func (d *defaultIdempotent) IsIdempotentMethod(fullMethod string) bool {
	return d.method[fullMethod]
}

func (d *defaultIdempotent) TryAcquire(ctx context.Context, id string) (resp any, isAcquire bool) {
	// 基于redis实现
	// 如果存在这个键就返回false
	retry, err := d.Redis.SetnxEx(id, "1", 60*60)
	if err != nil {
		return nil, false
	}
	if retry {
		return nil, true
	}
	resp, _ = d.Cache.Get(id)
	return resp, false
}

func (d *defaultIdempotent) SaveResp(ctx context.Context, id string, resp any, respErr error) error {
	d.Cache.Set(id, resp)
	return nil
}
