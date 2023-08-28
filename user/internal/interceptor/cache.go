package interceptor

import (
	"common/ggEncrypts"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"strings"
	"time"
	"user/internal/dal"
	"user/internal/repo"
)

type CacheInterceptor struct {
	cache    repo.Cache
	cacheMap map[string]any
}

func New() *CacheInterceptor {
	cacheMap := make(map[string]any)
	//cacheMap["/project.service.v1.ProjectService/FindProjectByMemId"] = &project.MyProjectResponse{}
	return &CacheInterceptor{cache: dal.Rc, cacheMap: cacheMap}
}

func (c *CacheInterceptor) Cache() grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		respType := c.cacheMap[info.FullMethod]
		if respType == nil {
			return handler(ctx, req)
		}
		// 先查询缓存 有缓存直接返回 没有缓存 先请求 再放入缓存
		cacheCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		reqJson, _ := json.Marshal(req)
		cacheKey := ggEncrypts.MD5(string(reqJson))
		respJson, _ := c.cache.Get(cacheCtx, info.FullMethod+"::"+cacheKey)
		if respJson != "" {
			json.Unmarshal([]byte(respJson), &respType)
			zap.L().Info(info.FullMethod + " 走了缓存")
			return respType, nil
		}
		resp, err = handler(ctx, req)
		bytes, _ := json.Marshal(resp)
		c.cache.Put(cacheCtx, info.FullMethod+"::"+cacheKey, string(bytes), 5*time.Minute)
		return
	})
}

func (c *CacheInterceptor) CacheInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		c = New()
		respType := c.cacheMap[info.FullMethod]
		if respType == nil {
			return handler(ctx, req)
		}
		// 先查询缓存 有缓存直接返回 没有缓存 先请求 再放入缓存
		cacheCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		reqJson, _ := json.Marshal(req)
		cacheKey := ggEncrypts.MD5(string(reqJson))
		respJson, _ := c.cache.Get(cacheCtx, info.FullMethod+"::"+cacheKey)
		if respJson != "" {
			json.Unmarshal([]byte(respJson), &respType)
			zap.L().Info(info.FullMethod + " 走了缓存")
			return respType, nil
		}
		resp, err = handler(ctx, req)
		bytes, _ := json.Marshal(resp)
		c.cache.Put(cacheCtx, info.FullMethod+"::"+cacheKey, string(bytes), 5*time.Minute)
		// hash key->task field->redisKey
		if strings.HasPrefix(info.FullMethod, "/task") {
			_ = c.cache.HSet(cacheCtx, "task", info.FullMethod+"::"+cacheKey, "")
		}
		return
	}
}
