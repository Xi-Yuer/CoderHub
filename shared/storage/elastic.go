package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type ElasticsearchImpl interface {
	CreateIndex(index string, docID int64, body io.Reader) error
	SearchByFields(index string, fields map[string]interface{}, page, pageSize int64) ([]int64, error)
	DeleteByID(index string, ids int64) error
}

type ElasticSearchClient struct {
	Client *elasticsearch.Client
}

func NewElasticSearchClient(cfg *elasticsearch.Config) (*ElasticSearchClient, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Addresses,
		APIKey:    cfg.APIKey,
	})
	if err != nil {
		return nil, err
	}
	return &ElasticSearchClient{Client: client}, nil
}

// CreateIndex 创建索引
func (c *ElasticSearchClient) CreateIndex(index string, docID int64, body io.Reader) error {
	// 检查索引是否存在
	exists, err := c.Client.Indices.Exists([]string{index})
	if err != nil {
		return err
	}
	if exists.StatusCode == 404 {
		// 创建索引
		_, err = c.Client.Indices.Create(index)
	}
	// 执行索引操作
	_, err = c.Client.Index(index, body, c.Client.Index.WithDocumentID(strconv.FormatInt(docID, 10)))
	if err != nil {
		return err
	}
	return nil
}

// SearchByFields 根据传入的字段查询数据，只返回数据的ID
func (c *ElasticSearchClient) SearchByFields(index string, searchValue string, keywords []string, page, pageSize int64) ([]int64, error) {
	// 提取数据的ID
	var ids []int64
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
	from := (page - 1) * pageSize
	searchResp, err := c.Client.Search(
		c.Client.Search.WithContext(context.Background()),
		c.Client.Search.WithIndex(index),
		c.Client.Search.WithBody(esutil.NewJSONReader(query)),
		c.Client.Search.WithFrom(int(from)),
		c.Client.Search.WithSize(int(pageSize)),
		c.Client.Search.WithTrackTotalHits(true),
		c.Client.Search.WithPretty(),
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
	}
	return ids, nil
}

// DeleteByID 根据ID删除索引
func (c *ElasticSearchClient) DeleteByID(index string, id int64) error {
	// 删除索引
	result, err := c.Client.Delete(index, strconv.FormatInt(id, 10))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(result.Body)

	if result.IsError() {
		return fmt.Errorf("error deleting Elasticsearch index: %s", result.String())
	}
	return nil
}
