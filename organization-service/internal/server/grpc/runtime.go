package grpc

import (
	_middleware "common/middleware"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	organizationpb "pb/types/organization"

	"organization/env"
	"organization/internal/interface/handler"
	"organization/internal/server/grpc/interceptors"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/oklog/run"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	organizationHandler           *handler.OrganizationHandler
	organizationMemberHandler     *handler.OrganizationMemberHandler
	organizationRoleHandler       *handler.OrganizationRoleHandler
	organizationPermissionHandler *handler.OrganizationPermissionHandler
	organizationBranchHandler     *handler.OrganizationBranchHandler

	groupDocumentHandler *handler.GroupDocumentHandler

	groupHandler             *handler.GroupHandler
	groupDealHandler         *handler.GroupDealHandler
	groupLogActivityHandler  *handler.GroupLogActivityHandler
	groupNotificationHandler *handler.GroupNotificationHandler
	groupSettingHandler      *handler.GroupSettingHandler
	groupMemberHandler       *handler.GroupMemberHandler

	organizationLogActivityHandler *handler.OrganizationLogActivityHandler
	logger                         *flogging.FabricLogger
}

func NewGRPCServer(
	organizationHandler *handler.OrganizationHandler,
	organizationMemberHandler *handler.OrganizationMemberHandler,
	organizationRoleHandler *handler.OrganizationRoleHandler,
	organizationPermissionHandler *handler.OrganizationPermissionHandler,
	organizationBranchHandler *handler.OrganizationBranchHandler,

	groupDocumentHandler *handler.GroupDocumentHandler,
	groupHandler *handler.GroupHandler,
	groupDealHandler *handler.GroupDealHandler,
	groupLogActivityHandler *handler.GroupLogActivityHandler,
	groupNotificationHandler *handler.GroupNotificationHandler,
	groupSettingHandler *handler.GroupSettingHandler,
	groupMemberHandler *handler.GroupMemberHandler,

	organizationLogActivityHandler *handler.OrganizationLogActivityHandler,
) *GRPCServer {
	return &GRPCServer{
		organizationHandler:           organizationHandler,
		organizationMemberHandler:     organizationMemberHandler,
		organizationRoleHandler:       organizationRoleHandler,
		organizationPermissionHandler: organizationPermissionHandler,
		organizationBranchHandler:     organizationBranchHandler,

		groupDocumentHandler: groupDocumentHandler,

		groupHandler:             groupHandler,
		groupDealHandler:         groupDealHandler,
		groupLogActivityHandler:  groupLogActivityHandler,
		groupNotificationHandler: groupNotificationHandler,
		groupSettingHandler:      groupSettingHandler,
		groupMemberHandler:       groupMemberHandler,

		organizationLogActivityHandler: organizationLogActivityHandler,
		logger:                         flogging.MustGetLogger(env.SERVICE_NAME),
	}
}

func (s *GRPCServer) Start() {
	g := &run.Group{}

	start, stop := buildGRPCServer(s.logger,
		s.organizationHandler,
		s.organizationMemberHandler,
		s.organizationRoleHandler,
		s.organizationPermissionHandler,
		s.organizationBranchHandler,
		s.groupDocumentHandler,
		s.groupHandler,
		s.groupDealHandler,
		s.groupLogActivityHandler,
		s.groupNotificationHandler,
		s.groupSettingHandler,
		s.groupMemberHandler,
		s.organizationLogActivityHandler,
	)

	g.Add(start, stop)

	signalStart, signalStop := buildSignalHandler(s.logger)
	g.Add(signalStart, signalStop)

	if err := g.Run(); err != nil {
		s.logger.Errorf("Server exited: %v", err)
	}
}

func buildGRPCServer(logger *flogging.FabricLogger,
	organizationHandler *handler.OrganizationHandler,
	organizationMemberHandler *handler.OrganizationMemberHandler,
	organizationRoleHandler *handler.OrganizationRoleHandler,
	organizationPermissionHandler *handler.OrganizationPermissionHandler,
	organizationBranchHandler *handler.OrganizationBranchHandler,

	groupDocumentHandler *handler.GroupDocumentHandler,
	groupHandler *handler.GroupHandler,
	groupDealHandler *handler.GroupDealHandler,
	groupLogActivityHandler *handler.GroupLogActivityHandler,
	groupNotificationHandler *handler.GroupNotificationHandler,
	groupSettingHandler *handler.GroupSettingHandler,
	groupMemberHandler *handler.GroupMemberHandler,

	organizationLogActivityHandler *handler.OrganizationLogActivityHandler,
) (func() error, func(error)) {
	unaryInts := []grpc.UnaryServerInterceptor{
		interceptors.UnaryLoggerInterceptor(logger),
		// interceptors.UnaryAuthInterceptor(logger),
		interceptors.UnaryRecoveryInterceptor(logger),
		_middleware.ParseGrpcMetadataContextMiddleware,
	}
	streamInts := []grpc.StreamServerInterceptor{}
	logger.Info(fmt.Sprintf("Starting GRPC server at port %s", env.GRPC_PORT))
	_, start, stop, err := BuildGRPCServer(
		fmt.Sprintf(":%s", env.GRPC_PORT),
		func(s *grpc.Server) {
			logger.Info("Registering OrganizationServiceServer")
			organizationpb.RegisterOrganizationServiceServer(s, organizationHandler)
			organizationpb.RegisterOrganizationMemberServiceServer(s, organizationMemberHandler)
			organizationpb.RegisterOrganizationRoleServiceServer(s, organizationRoleHandler)
			organizationpb.RegisterOrganizationPermissionServiceServer(s, organizationPermissionHandler)
			organizationpb.RegisterOrganizationBranchServiceServer(s, organizationBranchHandler)

			organizationpb.RegisterGroupDocumentServiceServer(s, groupDocumentHandler)

			organizationpb.RegisterGroupServiceServer(s, groupHandler)
			organizationpb.RegisterGroupDealServiceServer(s, groupDealHandler)
			organizationpb.RegisterGroupLogActivityServiceServer(s, groupLogActivityHandler)
			organizationpb.RegisterGroupNotificationServiceServer(s, groupNotificationHandler)
			organizationpb.RegisterGroupSettingServiceServer(s, groupSettingHandler)
			organizationpb.RegisterGroupMemberServiceServer(s, groupMemberHandler)

			organizationpb.RegisterOrganizationLogActivityServiceServer(s, organizationLogActivityHandler)
		},
		unaryInts,
		streamInts,
	)
	if err != nil {
		return nil, nil
	}

	return start, stop
}

func buildSignalHandler(logger *flogging.FabricLogger) (func() error, func(error)) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	return func() error {
			sig := <-sigs
			return fmt.Errorf("signal received: %v", sig)
		}, func(err error) {
			close(sigs)
		}
}

func BuildGRPCServer(
	addr string,
	registerFunc func(*grpc.Server),
	unaryInterceptors []grpc.UnaryServerInterceptor,
	streamInterceptors []grpc.StreamServerInterceptor,
) (
	server *grpc.Server,
	start func() error,
	stop func(err error),
	err error,
) {
	var opts []grpc.ServerOption

	if len(unaryInterceptors) > 0 {
		opts = append(opts, grpc.ChainUnaryInterceptor(unaryInterceptors...))
	}

	if len(streamInterceptors) > 0 {
		opts = append(opts, grpc.ChainStreamInterceptor(streamInterceptors...))
	}

	server = grpc.NewServer(opts...)
	registerFunc(server)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to listen: %w", err)
	}

	start = func() error {
		fmt.Println("Starting GRPC server at port ", addr)
		return server.Serve(lis)
	}

	stop = func(err error) {
		server.GracefulStop()
	}

	return server, start, stop, nil
}
