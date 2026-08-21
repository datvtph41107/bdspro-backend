package application

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"payment/internal/commerce/domain"
)

func TestOrderServiceScopesCommandAndFailsClosed(t *testing.T) {
	now := time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC)
	store := newMemoryStore()
	service := NewOrderService(store, staticReference{value: "QHPORDER1"}, func() time.Time { return now })
	command := CreateOrderCommand{Subject: subject(), Terms: terms(), CommandKey: "checkout-1", ExpiresAt: now.Add(time.Hour)}

	first, created, err := service.CreateOrder(context.Background(), command)
	if err != nil || !created || first.ID == 0 {
		t.Fatalf("first order = %+v created=%v err=%v", first, created, err)
	}
	replay, created, err := service.CreateOrder(context.Background(), command)
	if err != nil || created || replay.ID != first.ID {
		t.Fatalf("replay order = %+v created=%v err=%v", replay, created, err)
	}

	changed := command
	changed.Terms = terms()
	changed.Terms.Price.AmountMinor++
	if _, _, err := service.CreateOrder(context.Background(), changed); !errors.Is(err, domain.ErrOrderCommandConflict) {
		t.Fatalf("changed terms err=%v", err)
	}

	otherProduct := command
	otherProduct.Terms = terms()
	otherProduct.Terms.ProductCode = "another-product"
	otherProduct.CommandKey = command.CommandKey
	otherProduct.Terms.PlanCode = "another.pro"
	otherProduct.Terms.PlanVersionID = 12
	service2 := NewOrderService(store, staticReference{value: "QHPORDER2"}, func() time.Time { return now })
	second, created, err := service2.CreateOrder(context.Background(), otherProduct)
	if err != nil || !created || second.ID == first.ID {
		t.Fatalf("scoped command order = %+v created=%v err=%v", second, created, err)
	}

	invalid := command
	invalid.CommandKey = "missing-tier"
	invalid.Terms = terms()
	invalid.Terms.TierRank = 0
	if _, _, err := service.CreateOrder(context.Background(), invalid); !errors.Is(err, domain.ErrInvalidCommercialTerms) {
		t.Fatalf("missing tier err=%v", err)
	}
}

func TestOrderServiceConcurrentSameCommandHasOneCreator(t *testing.T) {
	now := time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC)
	store := newMemoryStore()
	service := NewOrderService(store, staticReference{value: "QHP-CONCURRENT"}, func() time.Time { return now })
	command := CreateOrderCommand{Subject: subject(), Terms: terms(), CommandKey: "same-command", ExpiresAt: now.Add(time.Hour)}

	type result struct {
		order   domain.Order
		created bool
		err     error
	}
	start := make(chan struct{})
	results := make(chan result, 2)
	var wg sync.WaitGroup
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			order, created, err := service.CreateOrder(context.Background(), command)
			results <- result{order: order, created: created, err: err}
		}()
	}
	close(start)
	wg.Wait()
	close(results)

	var rows []result
	createdCount := 0
	for item := range results {
		if item.err != nil {
			t.Fatalf("concurrent create err=%v", item.err)
		}
		if item.created {
			createdCount++
		}
		rows = append(rows, item)
	}
	if createdCount != 1 || len(rows) != 2 || rows[0].order.ID == 0 || rows[0].order.ID != rows[1].order.ID {
		t.Fatalf("concurrent results=%+v created_count=%d", rows, createdCount)
	}
}
