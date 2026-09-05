package connect

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// ProbeTimeout mặc định khi gateway check TCP upstream.
const ProbeTimeout = 800 * time.Millisecond

// ListFromViper đọc service.port.* từ config gateway rồi probe TCP.
// Dùng bởi gateway HTTP handler GET /internal/services.
func ListFromViper(envRuntime string) DiscoverResponse {
	raw := viper.GetStringMap("service.port")
	type item struct {
		name string
		port string
	}
	items := make([]item, 0, len(raw))
	for name, val := range raw {
		name = strings.TrimSpace(name)
		port := stringifyPort(val)
		if name == "" || port == "" || name == "gateway" {
			continue
		}
		items = append(items, item{name: name, port: port})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].name < items[j].name })

	results := make([]ServiceStatus, len(items))
	var wg sync.WaitGroup
	for i, it := range items {
		wg.Add(1)
		go func(i int, it item) {
			defer wg.Done()
			addr := fmt.Sprintf("%s:%s", HostOf(envRuntime, it.name), strings.TrimPrefix(it.port, ":"))
			st := ServiceStatus{Name: it.name, Address: addr}
			if err := probeTCP(addr, ProbeTimeout); err != nil {
				st.Up = false
				st.Error = err.Error()
			} else {
				st.Up = true
			}
			results[i] = st
		}(i, it)
	}
	wg.Wait()

	return DiscoverResponse{Services: results}
}

func stringifyPort(v interface{}) string {
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return strconv.FormatInt(int64(t), 10)
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}

func probeTCP(addr string, timeout time.Duration) error {
	d := net.Dialer{Timeout: timeout}
	conn, err := d.DialContext(context.Background(), "tcp", addr)
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}
