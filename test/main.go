package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"io"
	"net/http"
	"strconv"
)

func main() {

	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://es-cn-2ml47q1gp0003equa.public.elasticsearch.aliyuncs.com:9200",
		},
		APIKey: "VGtxVEk1WUJ6a1BZNGNYZVRnTmQ6VTdMQzZIc19UYmFVdWVvQnBLcEpfUQ==",
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	client, err := elasticsearch.NewClient(cfg)

	keywords := []string{"title", "summary", "content"}
	searchValue := "11"
	shouldClauses := make([]map[string]interface{}, 0, len(keywords))
	for _, field := range keywords {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"wildcard": map[string]interface{}{
				field: map[string]interface{}{
					"value": fmt.Sprintf("*%s*", searchValue),
				},
			},
		})
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": shouldClauses,
			},
		},
	}

	searchResp, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex("micro_post"),
		client.Search.WithBody(esutil.NewJSONReader(query)),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
	)
	// 解析搜索结果，获取 ids
	if err != nil {
		fmt.Println(err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(searchResp.Body)
	if searchResp.IsError() {
		fmt.Println("Error:", searchResp.String())
	} else {
		// 解析查询结果
		var response map[string]interface{}
		if err := json.NewDecoder(searchResp.Body).Decode(&response); err != nil {
			fmt.Println("Error:", searchResp.String())
		}
		// 提取数据的ID
		var ids []int64
		if hits, ok := response["hits"].(map[string]interface{})["hits"].([]interface{}); ok {
			for _, hit := range hits {
				if doc, ok := hit.(map[string]interface{}); ok {
					if idStr, ok := doc["_id"].(string); ok {
						if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
							ids = append(ids, id)
						}
					}
				}
			}
		}
		fmt.Println(ids)
	}
}
