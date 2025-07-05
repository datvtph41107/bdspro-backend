package config

import (
	"crm/infra/router"
	services "crm/infra/services"
	"crm/internal/usecase"
	_ "embed"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
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
	RuntimeRouter            *router.RuntimeRouter
	RuleUsecase              *usecase.RuleUsecase
	RuleEventUsecase         *usecase.RuleEventUsecase
	PipelineService          *services.PipelineService
	ContactService           *services.ContactService
	LeadService              *services.LeadService
	StageService             *services.StageService
	RuleService              *services.RuleService
	FollowService            *services.FollowService
	BlockService             *services.BlockService
	FriendGroupService       *services.FriendGroupService
	FriendService            *services.FriendService
	SharingAccessService     *services.SharingAccessService
	InvitationInstallService *services.InvitationInstallService
	CrmInternalService       *services.CrmInternalService
	Logger                   *flogging.FabricLogger
}

func NewConfigApp(
	runtimeRouter *router.RuntimeRouter,
	ruleUsecase *usecase.RuleUsecase,
	ruleEventUsecase *usecase.RuleEventUsecase,
	pipelineService *services.PipelineService,
	contactService *services.ContactService,
	customerService *services.LeadService,
	stageService *services.StageService,
	ruleService *services.RuleService,
	followService *services.FollowService,
	blockService *services.BlockService,
	friendGroupService *services.FriendGroupService,
	friendService *services.FriendService,
	sharingAccessService *services.SharingAccessService,
	invitationInstallService *services.InvitationInstallService,
	crmInternalService *services.CrmInternalService,
) *ConfigApp {
	return &ConfigApp{
		RuntimeRouter:            runtimeRouter,
		RuleUsecase:              ruleUsecase,
		RuleEventUsecase:         ruleEventUsecase,
		PipelineService:          pipelineService,
		ContactService:           contactService,
		LeadService:              customerService,
		StageService:             stageService,
		RuleService:              ruleService,
		FollowService:            followService,
		BlockService:             blockService,
		FriendGroupService:       friendGroupService,
		FriendService:            friendService,
		SharingAccessService:     sharingAccessService,
		InvitationInstallService: invitationInstallService,
		CrmInternalService:       crmInternalService,
		Logger:                   flogging.MustGetLogger(viper.GetString("server.name")),
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
