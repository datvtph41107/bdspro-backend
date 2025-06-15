package config

import (
	grpc_service "crm/infra/grpc"
	"crm/infra/router"
	"crm/internal/usecase"
	_ "embed"
)

const C_SESSION_ID = "C_SESSION_ID"

//go:embed app.yml
var ConfigFile []byte

// UserProperties chứa cấu hình của ứng dụng
type UserProperties struct {
	Jwt struct {
		SecretKey           string   `mapstructure:"key-generate"`
		TokenPrefix         string   `mapstructure:"tokenPrefix"`
		TokenExpirationDays int      `mapstructure:"tokenExpirationAfterDays"`
		AccessExpMinutes    int      `mapstructure:"accessExpAfterMinutes"`
		RefreshExpMinutes   uint64   `mapstructure:"refreshExpAfterMinutes"`
		ListPermit          []string `mapstructure:"listPermit"`
		AuthorizationHeader string   `mapstructure:"authorizationHeader"`
	} `mapstructure:"jwt"`
}

var AppProperties UserProperties

type ConfigApp struct {
	RuntimeRouter        *router.RuntimeRouter
	RuleUsecase          *usecase.RuleUsecase
	RuleEventUsecase     *usecase.RuleEventUsecase
	PipelineService      *grpc_service.PipelineService
	ContactService       *grpc_service.ContactService
	CustomerService      *grpc_service.CustomerService
	StageService         *grpc_service.StageService
	HistoryService       *grpc_service.HistoryService
	RuleService          *grpc_service.RuleService
	FollowService        *grpc_service.FollowService
	BlockService         *grpc_service.BlockService
	FriendGroupService   *grpc_service.FriendGroupService
	FriendService        *grpc_service.FriendService
	SharingAccessService *grpc_service.SharingAccessService
}

func NewConfigApp(
	runtimeRouter *router.RuntimeRouter,
	ruleUsecase *usecase.RuleUsecase,
	ruleEventUsecase *usecase.RuleEventUsecase,
	pipelineService *grpc_service.PipelineService,
	contactService *grpc_service.ContactService,
	customerService *grpc_service.CustomerService,
	stageService *grpc_service.StageService,
	historyService *grpc_service.HistoryService,
	ruleService *grpc_service.RuleService,
	followService *grpc_service.FollowService,
	blockService *grpc_service.BlockService,
	friendGroupService *grpc_service.FriendGroupService,
	friendService *grpc_service.FriendService,
	sharingAccessService *grpc_service.SharingAccessService,
) *ConfigApp {
	return &ConfigApp{
		RuntimeRouter:        runtimeRouter,
		RuleUsecase:          ruleUsecase,
		RuleEventUsecase:     ruleEventUsecase,
		PipelineService:      pipelineService,
		ContactService:       contactService,
		CustomerService:      customerService,
		StageService:         stageService,
		HistoryService:       historyService,
		RuleService:          ruleService,
		FollowService:        followService,
		BlockService:         blockService,
		FriendGroupService:   friendGroupService,
		FriendService:        friendService,
		SharingAccessService: sharingAccessService,
	}
}

// // LoadConfig đọc cấu hình từ file `.env` hoặc `config.yaml`
// func LoadProperties() {
// 	log.Printf("Đọc config")

// 	// var config UserProperties
// 	if err := viper.Unmarshal(&AppProperties); err != nil {
// 		log.Fatalf("Lỗi parse config: %v", err)
// 		// return nil, err
// 	}

// 	// return &config, nil
// }
