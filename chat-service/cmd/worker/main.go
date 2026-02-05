package worker

import "time"

func main() {
	for {
		events := repo.FetchPendingOutbox(100)

		for _, e := range events {
			redisStream.Publish("chat.events", e.Payload)
			repo.MarkSent(e.ID)
		}

		time.Sleep(200 * time.Millisecond)
	}
}
