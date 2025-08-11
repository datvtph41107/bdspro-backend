package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	inputRoot := "../../../protobuf/docs" // thư mục gốc chứa các thư mục con
	outputFile := "../../../gateway-service/docs/total/swagger.json"

	// Cấu trúc bản merge mặc định
	merged := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":   "Merged API",
			"version": "1.0.0",
		},
		"paths":      map[string]interface{}{},
		"components": map[string]interface{}{},
	}

	// Duyệt đệ quy tất cả file trong tree
	err := filepath.WalkDir(inputRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // lỗi truy cập file/dir
		}

		// Chỉ xử lý file thường có đuôi .json (hoặc .swagger.json)
		if d.Type().IsRegular() && filepath.Ext(d.Name()) == ".json" {
			if err := mergeFile(path, merged); err != nil {
				return fmt.Errorf("merge %s: %w", path, err)
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	// Tạo thư mục cha nếu chưa có
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		panic(err)
	}
	if err := EnsureFile(outputFile); err != nil {
		panic(err)
	}
	// Ghi kết quả
	if err := writeMerged(outputFile, merged); err != nil {
		panic(err)
	}

	fmt.Printf("✅ Đã gộp swagger thành công → %s\n", outputFile)
}

// EnsureFile checks if the file at `path` exists. If not, it creates an empty file.
func EnsureFile(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// Tạo file mới rỗng
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		fmt.Printf("📄 Created new file: %s\n", path)
		return nil
	}
	if err != nil {
		// Lỗi khác (ví dụ permission denied)
		return err
	}
	// File đã tồn tại
	return nil
}

// mergeFile đọc một file swagger JSON và gộp vào map merged
func mergeFile(path string, merged map[string]interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}

	// paths
	if paths, ok := doc["paths"].(map[string]interface{}); ok {
		for k, v := range paths {
			merged["paths"].(map[string]interface{})[k] = v
		}
	}

	// components
	if comps, ok := doc["components"].(map[string]interface{}); ok {
		for compType, compVal := range comps {
			dst := merged["components"].(map[string]interface{})
			if _, exists := dst[compType]; !exists {
				dst[compType] = compVal
				continue
			}

			// deep‑merge: components.<type>.*  (schemas, responses, …)
			srcMap, srcOK := compVal.(map[string]interface{})
			dstMap, dstOK := dst[compType].(map[string]interface{})
			if srcOK && dstOK {
				for k, v := range srcMap {
					dstMap[k] = v
				}
			}
		}
	}
	return nil
}

func writeMerged(outPath string, merged map[string]interface{}) error {
	out, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, out, 0o644)
}
