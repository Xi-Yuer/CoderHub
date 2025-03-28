package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	WebSite string
	Minio   MinioConfig
}

type MinioConfig struct {
	BASEURL         string
	Endpoint        string
	AccessKey       string
	SecretKey       string
	UseSSL          bool
	Bucket          string
	Region          string
	ThumbnailBucket string
	ThumbnailWidth  uint
}
