package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/metadata"
	"strconv"
)

func GetUserID(ctx context.Context) (int64, error) {
	value := ctx.Value("userId")
	if value == nil {
		return 0, fmt.Errorf("userId not found")
	}
	// 判断实际类型并处理
	switch v := value.(type) {
	case string:
		// 如果是 string 类型，直接转换为 int64
		userId, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("userId is not a valid int64")
		}
		return userId, nil
	case json.Number:
		// 如果是 json.Number 类型，转换为 int64
		userId, err := v.Int64()
		if err != nil {
			return 0, fmt.Errorf("userId is not a valid int64")
		}
		return userId, nil
	default:
		// 如果类型不匹配，返回错误
		return 0, fmt.Errorf("unexpected type for userId: %T", value)
	}
}

func SetUserMetaData(ctx context.Context) context.Context {
	value := ctx.Value("userId")
	md := metadata.New(
		map[string]string{
			"userId": fmt.Sprintf("%v", value), // 将 json.Number 转换为 string
		})
	return metadata.NewOutgoingContext(ctx, md)
}

func GetUserMetaData(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata not found")
	}
	value := md.Get("userId")
	if len(value) == 0 {
		return "", fmt.Errorf("key not found")
	}
	return value[0], nil
}
