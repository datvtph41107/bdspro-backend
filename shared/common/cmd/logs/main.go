// Lệnh đọc log dùng cùng logging.yml với process; chỉ hiển thị, không giữ state nghiệp vụ.
package main

import (
	"bufio"
	"common/logging"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type event map[string]any

func main() {
	config := flag.String("config", "logging.yml", "Đường dẫn config logging")
	service := flag.String("service", "", "Service cần xem; bỏ trống để nối toàn hệ thống")
	request := flag.String("request-id", "", "Lọc request_id hoặc operation_id")
	level := flag.String("level", "", "debug/info/warn/error")
	category := flag.String("category", "", "application/database/grpc/dependency/http/worker/lifecycle")
	limit := flag.Int("limit", 100, "Số sự kiện cuối cùng")
	color := flag.String("color", "auto", "auto/always/never; màu chỉ dành cho terminal")
	stdin := flag.Bool("stdin", false, "Đọc JSONL từ stdin, dùng với tail -F")
	flag.Parse()
	colored := *color == "always"
	if *color == "auto" {
		info, err := os.Stdout.Stat()
		colored = err == nil && info.Mode()&os.ModeCharDevice != 0 && os.Getenv("NO_COLOR") == ""
	}
	accept := func(e event) bool {
		return (*service == "" || e["service_name"] == *service+"-service") && (*request == "" || e["request_id"] == *request || e["operation_id"] == *request) && (*level == "" || strings.EqualFold(fmt.Sprint(e["level"]), *level)) && (*category == "" || e["category"] == *category)
	}
	consume := func(reader io.Reader, emit func(event)) error {
		scanner := bufio.NewScanner(reader)
		scanner.Buffer(make([]byte, 4096), 4<<20)
		for scanner.Scan() {
			var e event
			if json.Unmarshal(scanner.Bytes(), &e) == nil && accept(e) {
				emit(e)
			}
		}
		return scanner.Err()
	}
	if *stdin {
		if err := consume(os.Stdin, func(e event) { display(e, colored) }); err != nil {
			fail(err)
		}
		return
	}
	c, err := logging.Load(*config, *service)
	if err != nil {
		fail(err)
	}
	root := c.Root
	if *service != "" {
		if filepath.Base(*service) != *service || *service == ".." {
			fail(fmt.Errorf("invalid service"))
		}
		root = filepath.Join(root, *service)
	}
	if root == "" {
		fail(fmt.Errorf("logging config has no file root"))
	}
	var events []event
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(path, ".jsonl") {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		return consume(f, func(e event) {
			events = append(events, e)
			// Giới hạn memory khi đọc nhiều phiên; giữ các sự kiện mới nhất theo thời gian.
			if *limit > 0 && len(events) > 2**limit {
				sortEvents(events)
				events = events[len(events)-*limit:]
			}
		})
	})
	if err != nil {
		fail(err)
	}
	sortEvents(events)
	if *limit > 0 && len(events) > *limit {
		events = events[len(events)-*limit:]
	}
	for _, e := range events {
		display(e, colored)
	}
}
func sortEvents(events []event) {
	sort.SliceStable(events, func(i, j int) bool { return fmt.Sprint(events[i]["time"]) < fmt.Sprint(events[j]["time"]) })
}
func display(e event, color bool) {
	severity := fmt.Sprint(e["level"])
	if color {
		prefix := "\033[36m"
		switch severity {
		case "ERROR":
			prefix = "\033[31m"
		case "WARN":
			prefix = "\033[33m"
		case "INFO":
			prefix = "\033[32m"
		}
		severity = prefix + severity + "\033[0m"
	}
	fmt.Printf("%s %-5s %s [%v] %v\n", e["time"], severity, e["service_name"], e["category"], e["msg"])
	keys := []string{"request_id", "operation_id", "event_name", "grpc_method", "peer_service", "http_route", "http_status_code", "grpc_code", "handler", "field", "error_reason", "error_source", "query_source", "sql_operation", "sql", "rows", "duration_ms", "error", "source"}
	for _, key := range keys {
		if v, ok := e[key]; ok && v != "" {
			fmt.Printf("  %s=%v\n", key, v)
		}
	}
}
func fail(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
