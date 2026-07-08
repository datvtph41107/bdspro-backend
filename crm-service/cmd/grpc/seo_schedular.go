package cmd_grpc

import (
	"context"
	"crm/config"
	"crm/initial"
	"crm/internal/dto"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func startSeoScheduler(app *initial.InitialApp) {
	if !config.AppProperties.Seo.GenerationJob.Enabled {
		log.Println("SEO generation job disabled")
		return
	}

	if app == nil || app.SeoWorkerUsecase == nil {
		log.Println("SEO generation job disabled: SeoWorkerUsecase is nil")
		return
	}

	s := gocron.NewScheduler(time.Local)

	_, err := s.Every(1).Day().At("00:00").Do(func() {
		batchSize := config.AppProperties.Seo.GenerationJob.BatchSize
		if batchSize <= 0 {
			batchSize = 50
		}

		result, err := app.SeoWorkerUsecase.GeneratePendingSeoPages(context.Background(), &dto.SeoWorkerGenerateRequest{
			Limit:       batchSize,
			TriggerType: "scheduler",
			Reason:      "daily seo generation job",
		})
		if err != nil {
			log.Printf("seo generation job error: %v", err)
			return
		}

		if result != nil {
			log.Printf(
				"seo generation job done scanned=%d success=%d failed=%d skipped=%d",
				result.Scanned,
				result.Success,
				result.Failed,
				result.Skipped,
			)
		}
	})
	if err != nil {
		log.Printf("failed to schedule seo generation job: %v", err)
		return
	}

	s.StartAsync()
}
