package main

import (
	_db "common/db"
	qhprorpc "common/rpc"
	"common/rpcenv"
	"context"
	"file/config"
	_ "file/docs/file"
	filehttp "file/infra/handler/filehttp"
	"file/internal/fileauthorization"
	"file/internal/fileauthorization/authgrpc"
	"file/internal/filecontent"
	"file/internal/filemedia"
	"file/internal/fileversion"
	"file/internal/processlifecycle"
	"file/internal/serviceauth"
	versionhubgrpc "file/internal/versionauth/hubgrpc"
	"file/models"
	"file/repositories"
	"file/services"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"os/signal"
	authpb "pb/types/auth"
	hubpb "pb/types/hub"
	"syscall"
	"time"
)

// @title File Service API
// @version 1.0
// @description QHPRO file content, delivery, media and version-artifact service.
// @BasePath /v1/file
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := run(); err != nil {
		log.Printf("File-service stopped with error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	if err := config.LoadConfig(); err != nil {
		return fmt.Errorf("load file config: %w", err)
	}
	runtimeCfg, err := config.CurrentRuntimeConfig()
	if err != nil {
		return fmt.Errorf("load file runtime config: %w", err)
	}

	serviceAuth, ok := serviceauth.NewVerifier(
		runtimeCfg.FileService.ServiceAuthKey,
	)
	if !ok {
		return fmt.Errorf("load file runtime config: service.auth_key is required")
	}

	authConn, err := qhprorpc.NewClient(qhprorpc.ClientConfig{
		Target:          runtimeCfg.Auth.GRPCAddress,
		BackoffMaxDelay: 5 * time.Second,
		Credentials:     insecure.NewCredentials(),
		Transport:       rpcenv.LoadTransportConfig(),
	})
	if err != nil {
		return fmt.Errorf("configure auth permission client: %w", err)
	}
	defer authConn.Close()

	permissionAuthorizer := authgrpc.New(
		authpb.NewAuthInternalServiceClient(authConn),
		runtimeCfg.Auth.Timeout,
	)

	hubConn, err := qhprorpc.NewClient(qhprorpc.ClientConfig{
		Target:          runtimeCfg.Hub.GRPCAddress,
		BackoffMaxDelay: 5 * time.Second,
		Credentials:     insecure.NewCredentials(),
		Transport:       rpcenv.LoadTransportConfig(),
	})
	if err != nil {
		return fmt.Errorf("configure hub API-key client: %w", err)
	}
	defer hubConn.Close()

	apiKeyVerifier := versionhubgrpc.New(
		hubpb.NewApiKeyServiceClient(hubConn),
		runtimeCfg.Hub.Timeout,
	)

	db, err := _db.Open(_db.DatabaseConfig{DSN: runtimeCfg.DatabaseDSN})
	if err != nil {
		return fmt.Errorf("open file database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("obtain file database connection pool: %w", err)
	}
	defer sqlDB.Close()

	schemaPolicy, err := _db.LoadSchemaPolicy("file")
	if err != nil {
		return fmt.Errorf("load File schema policy: %w", err)
	}
	log.Printf("File database schema mode=%s source=%s", schemaPolicy.Mode, schemaPolicy.Source)
	if err := _db.ApplySchemaPolicy(db, schemaPolicy, models.AutoMigrate); err != nil {
		return fmt.Errorf("apply File schema policy: %w", err)
	}

	fileRepository := repositories.NewFileRepository(db)
	accessRepository := repositories.NewFileAccessRepository(db)

	storageService, err := services.NewStorageService(
		runtimeCfg.FileService.StorageRoot,
	)
	if err != nil {
		return fmt.Errorf("configure file storage: %w", err)
	}

	mediaProcessor := filemedia.NewProcessor(storageService)
	mediaUpload := filecontent.NewMediaUploadService(
		fileRepository,
		mediaProcessor,
		runtimeCfg.FileService.XorCryptKey,
	)
	fileContent := filecontent.NewService(
		fileRepository,
		storageService,
		runtimeCfg.FileService.XorCryptKey,
	)

	signedAccess := fileauthorization.NewSignedAccess(
		accessRepository,
		runtimeCfg.FileService.XorCryptKey,
		runtimeCfg.FileService.SignatureKey,
	)
	signedReadIssuer := fileauthorization.NewSignedReadIssuer(
		runtimeCfg.FileService.XorCryptKey,
		runtimeCfg.FileService.SignatureKey,
		runtimeCfg.FileService.SignatureExpire,
	)

	fileHTTP := filehttp.NewHandler(fileContent, serviceAuth)
	fileReadHTTP := filehttp.NewReadHandler(
		serviceAuth,
		signedAccess,
		permissionAuthorizer,
		runtimeCfg.FileService.XorCryptKey,
		storageService,
	)
	fileSignatureHTTP := filehttp.NewSignatureHandler(signedReadIssuer)
	fileVideoReadHTTP := filehttp.NewVideoReadHandler(
		runtimeCfg.FileService.XorCryptKey,
		storageService,
	)
	fileVideoUploadHTTP := filehttp.NewVideoUploadHandler(mediaUpload)

	versionReference, err := fileversion.NewReferenceCodec(
		runtimeCfg.FileService.XorCryptKey,
		runtimeCfg.FileService.SignatureKey,
	)
	if err != nil {
		return fmt.Errorf("configure version reference: %w", err)
	}

	versionUploader := fileversion.NewUploader(
		fileRepository,
		storageService,
		versionReference,
	)
	versionUploadHTTP := filehttp.NewVersionUploadHandler(
		apiKeyVerifier,
		versionUploader,
	)
	versionDownloadHTTP := filehttp.NewVersionDownloadHandler(
		versionReference,
		storageService,
	)

	router := filehttp.NewRouter(
		filehttp.RouterDependencies{
			Files:           fileHTTP,
			Read:            fileReadHTTP,
			Signature:       fileSignatureHTTP,
			VideoRead:       fileVideoReadHTTP,
			VideoUpload:     fileVideoUploadHTTP,
			VersionUpload:   versionUploadHTTP,
			VersionDownload: versionDownloadHTTP,
		},
		filehttp.RouterConfig{
			AllowedOrigins: runtimeCfg.FileService.CORSAllowedOrigins,
		},
	)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	httpServer := &http.Server{
		Addr:              ":" + runtimeCfg.HTTPPort,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return processlifecycle.Run(
		ctx,
		stop,
		httpServer,
		nil,
		processlifecycle.Config{ShutdownTimeout: 10 * time.Second},
	)
}
