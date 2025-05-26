package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"search/config"
	"search/models"
)

const indexName = "documents"

func IndexDocument(doc models.Document) error {
	// Normalize text về lowercase
	doc.Text = strings.ToLower(doc.Text)

	// Bước 1: Tìm xem đã tồn tại document có cùng text chưa
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"text.keyword": doc.Text, // exact match với keyword field
			},
		},
		"size": 1,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return fmt.Errorf("encode query failed: %w", err)
	}

	res, err := config.ES.Search(
		config.ES.Search.WithContext(context.Background()),
		config.ES.Search.WithIndex(indexName),
		config.ES.Search.WithBody(&buf),
	)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode == 404 {
		// return fmt.Errorf("search response error: %s", res.String())
	} else {
		var searchResult map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
			return fmt.Errorf("decode search result failed: %w", err)
		}

		hits := searchResult["hits"].(map[string]interface{})["hits"].([]interface{})
		if len(hits) > 0 {
			// Đã tồn tại document này, không cần lưu lại
			return nil
		}
	}

	// Bước 2: Lưu mới
	body, _ := json.Marshal(doc)
	_, err = config.ES.Index(
		indexName,
		bytes.NewReader(body),
		config.ES.Index.WithContext(context.Background()),
	)
	return err
}

func DeleteDocument(id string) error {
	_, err := config.ES.Delete(indexName, id)
	return err
}
func SearchDocuments(words string) ([]models.Document, error) {
	keyword := strings.ToLower(words)

	matching := map[string]interface{}{
		"size": 10,
		"query": map[string]interface{}{
			"function_score": map[string]interface{}{
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"should": []interface{}{
							map[string]interface{}{
								"match_phrase_prefix": map[string]interface{}{
									"text": map[string]interface{}{
										"query": keyword,
										"boost": 5,
									},
								},
							},
							map[string]interface{}{
								"match": map[string]interface{}{
									"text": map[string]interface{}{
										"query": keyword,
										"boost": 2,
									},
								},
							},
						},
					},
				},
				"field_value_factor": map[string]interface{}{
					"field":    "popularity",
					"factor":   1,
					"modifier": "log1p",
					"missing":  0,
				},
				"boost_mode": "sum",
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(matching); err != nil {
		return nil, fmt.Errorf("encode query error: %w", err)
	}

	res, err := config.ES.Search(
		config.ES.Search.WithContext(context.Background()),
		config.ES.Search.WithIndex(indexName),
		config.ES.Search.WithBody(&buf),
		config.ES.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("search error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		// Nếu trả về lỗi, decode lỗi để lấy message
		var errResp map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("search error, decode failed: %w", err)
		}

		message := ""
		if errInfo, ok := errResp["error"].(map[string]interface{}); ok {
			if reason, ok := errInfo["reason"].(string); ok {
				message = reason
			}
		}

		return nil, fmt.Errorf("elasticsearch error: %s (status %d)", message, res.StatusCode)
	}

	// Decode response thành công vào biến r
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	result := []models.Document{}

	// Duyệt các hit trong kết quả
	hits, ok := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response structure")
	}

	for _, hit := range hits {
		source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
		if !ok {
			continue // skip nếu không đúng kiểu
		}

		doc := models.Document{}

		if v, ok := source["text"].(string); ok {
			doc.Text = v
		}
		if v, ok := source["targetId"].(float64); ok {
			doc.TargetID = uint64(v)
		}
		if v, ok := source["type"].(float64); ok {
			doc.Type = uint(v)
		}
		if v, ok := source["popularity"].(float64); ok {
			doc.Popularity = int64(v)
		}

		result = append(result, doc)
	}

	return result, nil
}

func SaveKeywordsBatchBulk(keywords []string) error {
	var buf bytes.Buffer

	for _, keyword := range keywords {
		meta := map[string]map[string]string{
			"index": {"_index": indexName},
		}
		metaLine, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		buf.Write(metaLine)
		buf.WriteByte('\n')

		doc := models.Document{
			Text:       keyword,
			TargetID:   0,
			Type:       0,
			Popularity: 0,
		}
		docLine, err := json.Marshal(doc)
		if err != nil {
			return err
		}
		buf.Write(docLine)
		buf.WriteByte('\n')
	}

	res, err := config.ES.Bulk(
		bytes.NewReader(buf.Bytes()),
		config.ES.Bulk.WithContext(context.Background()),
		config.ES.Bulk.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk request error: %s", res.String())
	}
	return nil
}
