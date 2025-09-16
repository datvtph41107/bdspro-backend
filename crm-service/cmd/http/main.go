package cmd_http

import (
	"common"
	_jwt "common/jwt"
	"context"
	"crm/config"
	"crm/infra/db"
	"crm/initial"
	"crm/wire"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"
)

var publicRoutes = []string{}

var HttpCmd = &cobra.Command{
	Use:   "http",
	Short: "Run the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()

		parameters := common.SetupParameters{
			ConfigFile:   config.ConfigFile,
			Engine:       r,
			PublicRoutes: publicRoutes,
			Prefix:       "/crm",
		}

		common.SetupRestfulAPI(parameters, func() {
			// config.LoadProperties()
			log.Printf("Data: ", config.AppProperties)

			app, _ := wire.InitializeApp()
			db.AutoMigrate()

			r.Use(_jwt.JWTAuthMiddleware(nil))
			// app.RuntimeRouter.RuntimeRoutes(r)

			startScheduler(app)
		})
	},
}

func startScheduler(app *initial.InitialApp) {
	now := time.Now()
	nextHour := now.Truncate(time.Hour).Add(1 * time.Hour) // giờ tròn kế tiếp
	wait := time.Until(nextHour)
	app.RuleEventUsecase.TriggerJobAutoEvent(context.Background())

	time.AfterFunc(wait, func() {
		s := gocron.NewScheduler(time.Local)
		s.Every(1).Hour().Do(func() {
			app.RuleEventUsecase.TriggerJobAutoEvent(context.Background())
		})
		s.StartAsync()
	})
}
