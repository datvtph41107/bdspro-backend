package identity

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestActorHTTPContextIsolation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		profileID := request.Header.Get("X-Test-Profile-ID")
		if profileID == "" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		ctx, err := BindActor(request.Context(), Actor{ProfileID: 42, Role: "user", TokenType: "ACCESS"})
		if err != nil {
			http.Error(writer, "actor bind failed", http.StatusInternalServerError)
			return
		}
		actor, ok := ActorFromContext(ctx)
		if !ok {
			http.Error(writer, "actor missing", http.StatusInternalServerError)
			return
		}
		_, _ = fmt.Fprintf(writer, "%d:%s:%s", actor.ProfileID, actor.Role, actor.TokenType)
	}))
	defer server.Close()

	response, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("anonymous request error = %v", err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		t.Fatalf("anonymous status = %d", response.StatusCode)
	}

	const requests = 100
	var wait sync.WaitGroup
	errors := make(chan error, requests)
	for index := 0; index < requests; index++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
			request.Header.Set("X-Test-Profile-ID", "42")
			result, requestErr := http.DefaultClient.Do(request)
			if requestErr != nil {
				errors <- requestErr
				return
			}
			body, readErr := io.ReadAll(result.Body)
			result.Body.Close()
			if readErr != nil {
				errors <- readErr
				return
			}
			if string(body) != "42:user:ACCESS" {
				errors <- fmt.Errorf("body = %q", body)
			}
		}()
	}
	wait.Wait()
	close(errors)
	for err := range errors {
		t.Fatal(err)
	}
}
