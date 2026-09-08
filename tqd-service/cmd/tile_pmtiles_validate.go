package cmd

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/protomaps/go-pmtiles/pmtiles"
)

// validatePMTilesDir kiểm tra header + toàn bộ directory tree mỗi file .pmtiles trước khi serve.
// Trả danh sách tileset lỗi (vd. vietnam-v3 corrupt leaf → gzip panic trong go-pmtiles).
func validatePMTilesDir(tilesDir string, logger *log.Logger) []string {
	entries, err := os.ReadDir(tilesDir)
	if err != nil {
		logger.Printf("⚠️  không đọc được thư mục tiles %s: %v", tilesDir, err)
		return nil
	}

	var invalid []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".pmtiles") {
			continue
		}
		path := filepath.Join(tilesDir, e.Name())
		if err := validatePMTilesFile(path); err != nil {
			name := strings.TrimSuffix(e.Name(), ".pmtiles")
			invalid = append(invalid, name)
			logger.Printf("❌ pmtiles invalid: %s — %v", e.Name(), err)
		}
	}
	if len(invalid) == 0 {
		logger.Printf("✅ đã kiểm tra %d file .pmtiles trong %s", countPMTiles(entries), tilesDir)
	} else {
		logger.Printf("⚠️  %d file .pmtiles lỗi — tileset này sẽ trả 503: %v", len(invalid), invalid)
	}
	return invalid
}

func countPMTiles(entries []os.DirEntry) int {
	n := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".pmtiles") {
			n++
		}
	}
	return n
}

func validatePMTilesFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}
	size := info.Size()
	if size < int64(pmtiles.HeaderV3LenBytes) {
		return fmt.Errorf("file too small (%d bytes)", size)
	}

	headerBuf := make([]byte, pmtiles.HeaderV3LenBytes)
	if _, err := io.ReadFull(f, headerBuf); err != nil {
		return fmt.Errorf("read header: %w", err)
	}
	header, err := pmtiles.DeserializeHeader(headerBuf)
	if err != nil {
		return err
	}

	if header.RootLength == 0 {
		return fmt.Errorf("empty root directory")
	}
	if size < int64(header.RootOffset+header.RootLength) {
		return fmt.Errorf("truncated file: need %d bytes, have %d", header.RootOffset+header.RootLength, size)
	}

	return walkPMTilesDirectories(f, size, header)
}

func walkPMTilesDirectories(f *os.File, fileSize int64, header pmtiles.HeaderV3) error {
	compression := header.InternalCompression

	var walk func(dirOffset, dirLength uint64) error
	walk = func(dirOffset, dirLength uint64) error {
		if dirLength == 0 {
			return fmt.Errorf("directory at offset %d has zero length", dirOffset)
		}
		end := dirOffset + dirLength
		if end > uint64(fileSize) {
			return fmt.Errorf("directory at offset %d length %d exceeds file size %d", dirOffset, dirLength, fileSize)
		}

		buf := make([]byte, dirLength)
		if _, err := f.ReadAt(buf, int64(dirOffset)); err != nil {
			return fmt.Errorf("read directory at %d: %w", dirOffset, err)
		}

		entries, err := deserializeEntriesSafe(buf, compression)
		if err != nil {
			return fmt.Errorf("directory at offset %d: %w", dirOffset, err)
		}

		for _, entry := range entries {
			if entry.RunLength > 0 {
				continue
			}
			leafOffset := header.LeafDirectoryOffset + entry.Offset
			if err := walk(leafOffset, uint64(entry.Length)); err != nil {
				return err
			}
		}
		return nil
	}

	return walk(header.RootOffset, header.RootLength)
}

// deserializeEntriesSafe giống pmtiles.DeserializeEntries nhưng trả lỗi thay vì panic khi gzip/varint hỏng.
func deserializeEntriesSafe(data []byte, compression pmtiles.Compression) ([]pmtiles.EntryV3, error) {
	var reader io.Reader = bytes.NewReader(data)
	if compression == pmtiles.Gzip {
		gz, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("gzip decompress failed: %w", err)
		}
		defer gz.Close()
		reader = gz
	} else if compression != pmtiles.NoCompression {
		return nil, fmt.Errorf("unsupported internal compression %d", compression)
	}

	byteReader := bufio.NewReader(reader)
	numEntries, err := binary.ReadUvarint(byteReader)
	if err != nil {
		return nil, fmt.Errorf("entry count: %w", err)
	}

	entries := make([]pmtiles.EntryV3, numEntries)
	lastID := uint64(0)
	for i := uint64(0); i < numEntries; i++ {
		tmp, err := binary.ReadUvarint(byteReader)
		if err != nil {
			return nil, fmt.Errorf("tile id[%d]: %w", i, err)
		}
		entries[i].TileID = lastID + tmp
		lastID = entries[i].TileID
	}

	for i := uint64(0); i < numEntries; i++ {
		runLength, err := binary.ReadUvarint(byteReader)
		if err != nil {
			return nil, fmt.Errorf("run length[%d]: %w", i, err)
		}
		entries[i].RunLength = uint32(runLength)
	}

	for i := uint64(0); i < numEntries; i++ {
		length, err := binary.ReadUvarint(byteReader)
		if err != nil {
			return nil, fmt.Errorf("length[%d]: %w", i, err)
		}
		entries[i].Length = uint32(length)
	}

	for i := uint64(0); i < numEntries; i++ {
		tmp, err := binary.ReadUvarint(byteReader)
		if err != nil {
			return nil, fmt.Errorf("offset[%d]: %w", i, err)
		}
		if i > 0 && tmp == 0 {
			entries[i].Offset = entries[i-1].Offset + uint64(entries[i-1].Length)
		} else {
			entries[i].Offset = tmp - 1
		}
	}

	return entries, nil
}
