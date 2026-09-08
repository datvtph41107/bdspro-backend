package cmd

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

type failingListener struct{ err error }

func (l failingListener) Accept() (net.Conn, error) { return nil, l.err }
func (l failingListener) Close() error              { return nil }
func (l failingListener) Addr() net.Addr            { return testAddr("failing-listener") }

type testAddr string

func (a testAddr) Network() string { return "test" }
func (a testAddr) String() string  { return string(a) }

func TestServeHTTPPropagatesServeFailure(t *testing.T) {
	rootErr := errors.New("accept failed")
	readiness := newReadinessState()
	err := serveHTTP(
		context.Background(),
		&http.Server{Handler: http.NewServeMux()},
		failingListener{err: rootErr},
		time.Second,
		readiness,
	)
	if !errors.Is(err, rootErr) {
		t.Fatalf("serveHTTP() error=%v, want wrapped %v", err, rootErr)
	}
}

func TestServeHTTPDrainsOnContextCancellation(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/livez", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "alive")
	})

	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: time.Second,
	}

	readiness := newReadinessState()

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)

	go func() {
		done <- serveHTTP(
			ctx,
			server,
			listener,
			time.Second,
			readiness,
		)
	}()

	client := &http.Client{
		Timeout: 250 * time.Millisecond,
	}

	url := "http://" + listener.Addr().String() + "/livez"

	deadline := time.Now().Add(2 * time.Second)

	for {
		resp, requestErr := client.Get(url)
		if requestErr == nil {
			_ = resp.Body.Close()
			break
		}

		if time.Now().After(deadline) {
			t.Fatalf(
				"HTTP server did not become reachable: %v",
				requestErr,
			)
		}

		time.Sleep(10 * time.Millisecond)
	}

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf(
				"serveHTTP() cancellation error=%v",
				err,
			)
		}

		if readiness.isReady() {
			t.Fatal(
				"Gateway remains ready after shutdown begins",
			)
		}

	case <-time.After(2 * time.Second):
		t.Fatal(
			"serveHTTP() did not return after cancellation",
		)
	}
}
