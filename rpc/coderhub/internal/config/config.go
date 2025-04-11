package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	WebSite       string
	Minio         MinioConfig
	ElasticSearch ElasticSearch
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

type ElasticSearch struct {
	Endpoint []string `json:"endpoint,env=Endpoint"`
	APIKEY   string   `json:"apikey,env=APIKEY"`
}
