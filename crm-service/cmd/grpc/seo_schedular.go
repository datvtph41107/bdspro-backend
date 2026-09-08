package cmd_grpc

import (
	"context"
	"crm/config"
	"crm/initial"
	"crm/internal/dto"
	"log"
	"time"
)

func runSeoScheduler(ctx context.Context, app *initial.InitialApp) {
	if !config.AppProperties.Seo.GenerationJob.Enabled {
		log.Println("SEO generation job disabled")
		<-ctx.Done()
		return
	}
	if app == nil || app.SeoWorkerUsecase == nil {
		log.Println("SEO generation job disabled: SeoWorkerUsecase is nil")
		<-ctx.Done()
		return
	}

	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		timer := time.NewTimer(time.Until(next))
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return
		case <-timer.C:
			batchSize := config.AppProperties.Seo.GenerationJob.BatchSize
			if batchSize <= 0 {
				batchSize = 50
			}
			result, err := app.SeoWorkerUsecase.GeneratePendingSeoPages(ctx, &dto.SeoWorkerGenerateRequest{
				Limit: batchSize, TriggerType: "scheduler", Reason: "daily seo generation job",
			})
			if err != nil {
				log.Printf("seo generation job error: %v", err)
				continue
			}
			if result != nil {
				log.Printf("seo generation job done scanned=%d success=%d failed=%d skipped=%d", result.Scanned, result.Success, result.Failed, result.Skipped)
			}
		}
	}
}
