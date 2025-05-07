package common

import (
	"bytes"
	_db "common/db"
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFile "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type SetupParameters struct {
	ConfigFile    []byte
	Engine        *gin.Engine
	PublicRoutes  []string
	Prefix        string
	IgnoreSwagger bool
}

func ProfileIDInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req) // Không có metadata thì tiếp tục
	}

	profileIDList := md.Get("profileid")
	if len(profileIDList) > 0 {
		profileIDStr := profileIDList[0]

		profileID, err := strconv.ParseUint(profileIDStr, 10, 64)
		if err == nil {
			log.Printf("Gán profileID", profileID)
			// Gán profileID vào context mới
			ctx = context.WithValue(ctx, "profileId", profileID)
		}
	}

	// Gọi handler tiếp theo (controller thực sự)
	return handler(ctx, req)
}

func SetupGRPCServer(parameters SetupParameters, setup func(*grpc.Server)) {
	LoadConfig(parameters.ConfigFile)
	// _db.NewDB()
	if !parameters.IgnoreSwagger {
		parameters.Engine.GET(parameters.Prefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFile.Handler))
	}

	// models.AutoMigrate()
	// parameters.Engine.Use(cors.New(cors.Config{
	// 	AllowOrigins: []string{"http://localhost:3000", "http://14.225.210.29:8000"}, // Cho phép gọi từ frontend
	// 	AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	// 	// ExposeHeaders:    []string{"Content-Length"},
	// 	ExposeHeaders:    []string{"Content-Length", "Set-Cookie"}, // Cho phép frontend thấy Set-Cookie
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	port := fmt.Sprintf(":%s", viper.GetString("server.tcp_port"))

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(ProfileIDInterceptor),
	)
	setup(s)

	log.Printf("Server TCP listen on %v", port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Println("TCP server is running on", port)
	// parameters.Engine.Run(port)
}

func SetupRestfulAPI(parameters SetupParameters, setup func()) {
	LoadConfig(parameters.ConfigFile)
	_db.NewDB()
	parameters.Engine.GET(parameters.Prefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFile.Handler))
	// models.AutoMigrate()
	// parameters.Engine.Use(cors.New(cors.Config{
	// 	AllowOrigins: []string{"http://localhost:3000", "http://14.225.210.29:8000"}, // Cho phép gọi từ frontend
	// 	AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	// 	// ExposeHeaders:    []string{"Content-Length"},
	// 	ExposeHeaders:    []string{"Content-Length", "Set-Cookie"}, // Cho phép frontend thấy Set-Cookie
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	setup()

	port := fmt.Sprintf(":%s", viper.GetString("server.port"))
	log.Println("Chat server is running on", port)
	parameters.Engine.Run(port)
}

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig(configFile []byte) {
	// viper.SetConfigName("config")   // Tên file không có phần mở rộng
	viper.SetConfigType("yaml") // Loại file config
	// viper.AddConfigPath("./config") // Đọc từ thư mục gốc

	// Đọc biến môi trường (nếu có)
	viper.AutomaticEnv()

	// Đọc file config
	err := viper.ReadConfig(bytes.NewReader(configFile))
	if err != nil {
		log.Fatalf("Lỗi khi đọc file config: %v", err)
	}
}
