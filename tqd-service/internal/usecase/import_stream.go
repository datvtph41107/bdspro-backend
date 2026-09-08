package usecase

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type geoJSONFeature struct {
	Type       string                 `json:"type"`
	Geometry   map[string]interface{} `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

func skipJSONValue(d *json.Decoder) error {
	t, err := d.Token()
	if err != nil {
		return err
	}
	switch v := t.(type) {
	case json.Delim:
		switch v {
		case '{':
			for d.More() {
				if _, err := d.Token(); err != nil {
					return err
				}
				if err := skipJSONValue(d); err != nil {
					return err
				}
			}
			_, err = d.Token()
			return err
		case '[':
			for d.More() {
				if err := skipJSONValue(d); err != nil {
					return err
				}
			}
			_, err = d.Token()
			return err
		default:
			return fmt.Errorf("geojson: unexpected delimiter %v", v)
		}
	default:
		return nil
	}
}

func streamFeatureCollectionFeatures(r io.Reader, fn func(f geoJSONFeature, index int) error) error {
	dec := json.NewDecoder(bufio.NewReaderSize(r, 1024*1024))
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	d, ok := tok.(json.Delim)
	if !ok || d != '{' {
		return fmt.Errorf("geojson: root must be an object")
	}
	idx := 0
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return err
		}
		key, ok := keyTok.(string)
		if !ok {
			return fmt.Errorf("geojson: invalid object key")
		}
		if key == "features" {
			ft, err := dec.Token()
			if err != nil {
				return err
			}
			fd, ok := ft.(json.Delim)
			if !ok || fd != '[' {
				return fmt.Errorf("geojson: features must be an array")
			}
			for dec.More() {
				var feat geoJSONFeature
				if err := dec.Decode(&feat); err != nil {
					return err
				}
				if feat.Type != "" && feat.Type != "Feature" {
					idx++
					continue
				}
				if err := fn(feat, idx); err != nil {
					return err
				}
				idx++
			}
			endArr, err := dec.Token()
			if err != nil {
				return err
			}
			if ed, ok := endArr.(json.Delim); !ok || ed != ']' {
				return fmt.Errorf("geojson: expected end of features array")
			}
			continue
		}
		if err := skipJSONValue(dec); err != nil {
			return err
		}
	}
	endObj, err := dec.Token()
	if err != nil {
		return err
	}
	if ed, ok := endObj.(json.Delim); !ok || ed != '}' {
		return fmt.Errorf("geojson: expected end of root object")
	}
	return nil
}

func forEachNDJSONFeature(path string, fn func(f geoJSONFeature, index int) error) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	br := bufio.NewReaderSize(f, 1024*1024)
	idx := 0
	for {
		line, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			if err == io.EOF {
				break
			}
			continue
		}
		var feat geoJSONFeature
		if err := json.Unmarshal([]byte(line), &feat); err != nil {
			return fmt.Errorf("ndjson line %d: %w", idx, err)
		}
		if feat.Type != "" && feat.Type != "Feature" {
			return fmt.Errorf("ndjson line %d: expected Feature, got %q", idx, feat.Type)
		}
		if err := fn(feat, idx); err != nil {
			return err
		}
		idx++
		if err == io.EOF {
			break
		}
	}
	return nil
}

func peekFileHead(path string, max int) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := make([]byte, max)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return buf[:n], nil
}

func importFormatFromPeek(head []byte) string {
	s := string(head)
	if strings.Contains(s, "FeatureCollection") {
		return "geojson"
	}
	return ""
}

func tryImportSmallWholeFile(path string, maxSize int64, onFeature func(f geoJSONFeature, index int) error) (handled bool, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if fi.Size() > maxSize {
		return false, nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	var head struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(raw, &head); err != nil {
		return false, nil
	}
	if head.Type == "Feature" {
		var feat geoJSONFeature
		if err := json.Unmarshal(raw, &feat); err != nil {
			return false, err
		}
		return true, onFeature(feat, 0)
	}
	if head.Type == "FeatureCollection" {
		var wrapper struct {
			Features []geoJSONFeature `json:"features"`
		}
		if err := json.Unmarshal(raw, &wrapper); err != nil {
			return false, err
		}
		for i, feat := range wrapper.Features {
			if err := onFeature(feat, i); err != nil {
				return true, err
			}
		}
		return true, nil
	}
	return false, nil
}

func processImportFileByFormat(path, format string, onFeature func(f geoJSONFeature, index int) error) error {
	f := strings.ToLower(strings.TrimSpace(format))
	switch f {
	case "", "auto":
		head, err := peekFileHead(path, 256*1024)
		if err != nil {
			return err
		}
		if importFormatFromPeek(head) == "geojson" {
			r, err := os.Open(path)
			if err != nil {
				return err
			}
			defer r.Close()
			return streamFeatureCollectionFeatures(r, onFeature)
		}
		handled, err := tryImportSmallWholeFile(path, 100<<20, onFeature)
		if err != nil {
			return err
		}
		if handled {
			return nil
		}
		return forEachNDJSONFeature(path, onFeature)
	case "geojson", "json":
		head, err := peekFileHead(path, 256*1024)
		if err != nil {
			return err
		}
		if importFormatFromPeek(head) == "geojson" {
			r, err := os.Open(path)
			if err != nil {
				return err
			}
			defer r.Close()
			return streamFeatureCollectionFeatures(r, onFeature)
		}
		fi, err := os.Stat(path)
		if err != nil {
			return err
		}
		handled, err := tryImportSmallWholeFile(path, fi.Size(), onFeature)
		if err != nil {
			return err
		}
		if !handled {
			return fmt.Errorf("geojson: expected FeatureCollection or a single Feature object")
		}
		return nil
	case "ndjson":
		return forEachNDJSONFeature(path, onFeature)
	default:
		return fmt.Errorf("unsupported fileFormat %q (use auto, geojson, ndjson)", format)
	}
}

func writeUploadFileChunked(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	const bufSize = 1024 * 1024
	r := bytes.NewReader(data)
	buf := make([]byte, bufSize)
	_, err = io.CopyBuffer(out, r, buf)
	return err
}
