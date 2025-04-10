package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type ElasticsearchImpl interface {
	CreateIndex(index string, body io.Reader) error
	SearchByFields(index string, fields map[string]interface{}) ([]int64, error)
	DeleteByID(index string, ids int64) error
}

type ElasticSearchClient struct {
	Client *elasticsearch.Client
}

func NewElasticSearchClient(cfg *elasticsearch.Config) (*ElasticSearchClient, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	})
	if err != nil {
		return nil, err
	}
	return &ElasticSearchClient{Client: client}, nil
}

// CreateIndex 创建索引
func (c *ElasticSearchClient) CreateIndex(index string, docID int64, body io.Reader) error {
	// 创建索引
	result, err := c.Client.Index(index, body, c.Client.Index.WithDocumentID(strconv.FormatInt(docID, 10)))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(result.Body)
	if result.IsError() {
		return fmt.Errorf("error creating index in Elasticsearch: %s", result.String())
	}
	return nil
}

// SearchByFields 根据传入的字段查询数据，只返回数据的ID
func (c *ElasticSearchClient) SearchByFields(index string, fields map[string]interface{}) ([]int64, error) {
	// 构建查询
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{},
			},
		},
	}

	for field, value := range fields {
		// 过滤空值
		if value == nil ||
			(reflect.TypeOf(value).Kind() == reflect.String && value.(string) == "") ||
			(reflect.TypeOf(value).Kind() == reflect.Int64 && value.(int64) == 0) ||
			(reflect.TypeOf(value).Kind() == reflect.Int32 && value.(int32) == 0) ||
			(reflect.TypeOf(value).Kind() == reflect.Float64 && value.(float64) == 0) ||
			(reflect.TypeOf(value).Kind() == reflect.Float32 && value.(float32) == 0) ||
			(reflect.TypeOf(value).Kind() == reflect.Bool && !value.(bool)) {
			continue
		}

		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			map[string]interface{}{
				"match": map[string]interface{}{
					field: value,
				},
			},
		)
	}

	// 执行搜索
	result, err := c.Client.Search(
		c.Client.Search.WithIndex(index),
		c.Client.Search.WithBody(esutil.NewJSONReader(query)),
		c.Client.Search.WithTrackTotalHits(true),
		c.Client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(result.Body)

	if result.IsError() {
		return nil, fmt.Errorf("error searching Elasticsearch: %s", result.String())
	}

	// 解析响应以提取 ID
	var response map[string]interface{}
	if err := json.NewDecoder(result.Body).Decode(&response); err != nil {
		return nil, err
	}

	hits := response["hits"].(map[string]interface{})["hits"].([]interface{})
	ids := make([]int64, 0, len(hits))
	for _, hit := range hits {
		id := hit.(map[string]interface{})["_id"].(string)
		s, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, s)
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
