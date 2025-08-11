package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {
	urls := []string{
		"http://localhost:8000/v1/bdspro/v2/user/product/market",
	}

	var wg sync.WaitGroup
	client := &http.Client{}
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url := urls[rand.Intn(len(urls))]

			resp, err := client.Get(url)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Printf("Status for %s: %s\n", url, resp.Status)
			resp.Body.Close()
		}()
	}

	wg.Wait()
}
