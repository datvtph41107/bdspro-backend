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
	"organization/infrastructure/handler"
	"organization/infrastructure/server/grpc/interceptors"

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
	dealMilestoneHandler     *handler.DealMilestoneHandler
	groupLogActivityHandler  *handler.GroupLogActivityHandler
	groupNotificationHandler *handler.GroupNotificationHandler
	groupSettingHandler      *handler.GroupSettingHandler
	groupMemberHandler       *handler.GroupMemberHandler
	investmentHandler        *handler.InvestmentHandler

	organizationLogActivityHandler *handler.OrganizationLogActivityHandler
	internalOrganizationHandler    *handler.InternalOrganizationHandler
	dealInvitationHandler          *handler.DealInvitationHandler
	dealCommissionHandler          *handler.DealCommissionHandler
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
	dealMilestoneHandler *handler.DealMilestoneHandler,
	groupLogActivityHandler *handler.GroupLogActivityHandler,
	groupNotificationHandler *handler.GroupNotificationHandler,
	groupSettingHandler *handler.GroupSettingHandler,
	groupMemberHandler *handler.GroupMemberHandler,
	investmentHandler *handler.InvestmentHandler,

	organizationLogActivityHandler *handler.OrganizationLogActivityHandler,
	internalOrganizationHandler *handler.InternalOrganizationHandler,
	dealInvitationHandler *handler.DealInvitationHandler,
	dealCommissionHandler *handler.DealCommissionHandler,
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
		dealMilestoneHandler:     dealMilestoneHandler,
		groupLogActivityHandler:  groupLogActivityHandler,
		groupNotificationHandler: groupNotificationHandler,
		groupSettingHandler:      groupSettingHandler,
		groupMemberHandler:       groupMemberHandler,
		investmentHandler:        investmentHandler,

		organizationLogActivityHandler: organizationLogActivityHandler,
		internalOrganizationHandler:    internalOrganizationHandler,
		dealInvitationHandler:          dealInvitationHandler,
		dealCommissionHandler:          dealCommissionHandler,
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
		s.dealMilestoneHandler,
		s.groupLogActivityHandler,
		s.groupNotificationHandler,
		s.groupSettingHandler,
		s.groupMemberHandler,
		s.investmentHandler,
		s.organizationLogActivityHandler,
		s.internalOrganizationHandler,
		s.dealInvitationHandler,
		s.dealCommissionHandler,
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
	dealMilestoneHandler *handler.DealMilestoneHandler,
	groupLogActivityHandler *handler.GroupLogActivityHandler,
	groupNotificationHandler *handler.GroupNotificationHandler,
	groupSettingHandler *handler.GroupSettingHandler,
	groupMemberHandler *handler.GroupMemberHandler,
	investmentHandler *handler.InvestmentHandler,

	organizationLogActivityHandler *handler.OrganizationLogActivityHandler,
	internalOrganizationHandler *handler.InternalOrganizationHandler,
	dealInvitationHandler *handler.DealInvitationHandler,
	dealCommissionHandler *handler.DealCommissionHandler,
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
			fmt.Println("Go1")
			organizationpb.RegisterOrganizationMemberServiceServer(s, organizationMemberHandler)
			fmt.Println("Go2")
			organizationpb.RegisterOrganizationRoleServiceServer(s, organizationRoleHandler)
			fmt.Println("Go3")
			organizationpb.RegisterOrganizationPermissionServiceServer(s, organizationPermissionHandler)
			fmt.Println("Go4")
			organizationpb.RegisterOrganizationBranchServiceServer(s, organizationBranchHandler)

			organizationpb.RegisterGroupDocumentServiceServer(s, groupDocumentHandler)
			fmt.Println("Go5")
			organizationpb.RegisterGroupServiceServer(s, groupHandler)
			fmt.Println("Go6")
			organizationpb.RegisterGroupDealServiceServer(s, groupDealHandler)
			fmt.Println("Go7")
			organizationpb.RegisterDealMilestoneServiceServer(s, dealMilestoneHandler)
			organizationpb.RegisterGroupLogActivityServiceServer(s, groupLogActivityHandler)
			organizationpb.RegisterGroupNotificationServiceServer(s, groupNotificationHandler)
			organizationpb.RegisterGroupSettingServiceServer(s, groupSettingHandler)
			organizationpb.RegisterGroupMemberServiceServer(s, groupMemberHandler)
			fmt.Println("Go8")
			organizationpb.RegisterOrganizationLogActivityServiceServer(s, organizationLogActivityHandler)
			fmt.Println("Go9")
			organizationpb.RegisterInvestmentServiceServer(s, investmentHandler)
			fmt.Println("Go10")
			organizationpb.RegisterInternalOrganizationServiceServer(s, internalOrganizationHandler)
			organizationpb.RegisterDealInvitationServiceServer(s, dealInvitationHandler)
			organizationpb.RegisterDealCommissionServiceServer(s, dealCommissionHandler)
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
