package storage

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type ElasticsearchImpl interface {
	SearchByFields(index string, fields map[string]interface{}) ([]int64, error)
}

type ElasticSearchClient struct {
	Client *elasticsearch.Client
}

func NewElasticSearchClient(cfg *elasticsearch.Config) (*ElasticSearchClient, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &ElasticSearchClient{Client: client}, nil
}

// 根据传入的字段查询数据，只返回数据的ID
func (c *ElasticSearchClient) SearchByFields(index string, fields map[string]interface{}) ([]int64, error) {
	// Construct the query using the fields
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{},
			},
		},
	}

	for field, value := range fields {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			map[string]interface{}{
				"match": map[string]interface{}{
					field: value,
				},
			},
		)
	}

	// Perform the search
	res, err := c.Client.Search(
		c.Client.Search.WithIndex(index),
		c.Client.Search.WithBody(esutil.NewJSONReader(query)),
		c.Client.Search.WithTrackTotalHits(true),
		c.Client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching Elasticsearch: %s", res.String())
	}

	// Parse the response to extract IDs
	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
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
