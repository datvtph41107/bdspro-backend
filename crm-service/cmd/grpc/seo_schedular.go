package cmd_grpc

import (
	"context"
	"crm/config"
	"crm/initial"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func startSeoScheduler(app *initial.InitialApp) {
	if !config.AppProperties.Seo.GenerationJob.Enabled {
		log.Println("SEO generation job disabled")
		return
	}

	s := gocron.NewScheduler(time.Local)

	_, err := s.Every(1).Day().At("00:00").Do(func() {
		if err := app.SeoDomainUsecase.TriggerSeoGenerationJob(context.Background()); err != nil {
			log.Printf("seo generation job error: %v", err)
		}
	})
	if err != nil {
		log.Printf("failed to schedule seo generation job: %v", err)
		return
	}

	s.StartAsync()
}
