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
	businessDomainHandler         *handler.BusinessDomainHandler

	groupDocumentHandler *handler.GroupDocumentHandler

	groupHandler             *handler.GroupHandler
	groupDealHandler         *handler.DealHandler
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
	colorHandler                   *handler.ColorHandler
	internalNoteHandler            *handler.InternalNoteHandler
	logger                         *flogging.FabricLogger
}

func NewGRPCServer(
	organizationHandler *handler.OrganizationHandler,
	organizationMemberHandler *handler.OrganizationMemberHandler,
	organizationRoleHandler *handler.OrganizationRoleHandler,
	organizationPermissionHandler *handler.OrganizationPermissionHandler,
	organizationBranchHandler *handler.OrganizationBranchHandler,
	businessDomainHandler *handler.BusinessDomainHandler,

	groupDocumentHandler *handler.GroupDocumentHandler,
	groupHandler *handler.GroupHandler,
	groupDealHandler *handler.DealHandler,
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
	colorHandler *handler.ColorHandler,
	internalNoteHandler *handler.InternalNoteHandler,
) *GRPCServer {
	return &GRPCServer{
		organizationHandler:           organizationHandler,
		organizationMemberHandler:     organizationMemberHandler,
		organizationRoleHandler:       organizationRoleHandler,
		organizationPermissionHandler: organizationPermissionHandler,
		organizationBranchHandler:     organizationBranchHandler,
		businessDomainHandler:         businessDomainHandler,

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
		colorHandler:                   colorHandler,
		logger:                         flogging.MustGetLogger(env.SERVICE_NAME),
		internalNoteHandler:            internalNoteHandler,
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
		s.businessDomainHandler,
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
		s.colorHandler,
		s.internalNoteHandler,
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
	businessDomainHandler *handler.BusinessDomainHandler,

	groupDocumentHandler *handler.GroupDocumentHandler,
	groupHandler *handler.GroupHandler,
	dealHandler *handler.DealHandler,
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
	colorHandler *handler.ColorHandler,
	internalNoteHandler *handler.InternalNoteHandler,
) (func() error, func(error)) {
	unaryInts := []grpc.UnaryServerInterceptor{
		interceptors.UnaryLoggerInterceptor(logger),
		// interceptors.UnaryAuthInterceptor(logger),
		interceptors.UnaryRecoveryInterceptor(logger),
		_middleware.ParseGrpcMetadataContextMiddleware,
	}
	streamInts := []grpc.StreamServerInterceptor{
		_middleware.ParseGrpcMetadataContextStreamMiddleware,
	}
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
			organizationpb.RegisterBusinessDomainServiceServer(s, businessDomainHandler)

			organizationpb.RegisterGroupDocumentServiceServer(s, groupDocumentHandler)
			organizationpb.RegisterGroupServiceServer(s, groupHandler)
			organizationpb.RegisterDealMilestoneServiceServer(s, dealMilestoneHandler)
			organizationpb.RegisterGroupLogActivityServiceServer(s, groupLogActivityHandler)
			organizationpb.RegisterGroupNotificationServiceServer(s, groupNotificationHandler)
			organizationpb.RegisterGroupSettingServiceServer(s, groupSettingHandler)
			organizationpb.RegisterGroupMemberServiceServer(s, groupMemberHandler)
			organizationpb.RegisterOrganizationLogActivityServiceServer(s, organizationLogActivityHandler)
			organizationpb.RegisterInvestmentServiceServer(s, investmentHandler)
			organizationpb.RegisterInternalOrganizationServiceServer(s, internalOrganizationHandler)

			organizationpb.RegisterDealServiceServer(s, dealHandler)
			organizationpb.RegisterDealMemberServiceServer(s, dealInvitationHandler)
			organizationpb.RegisterDealCommissionServiceServer(s, dealCommissionHandler)
			organizationpb.RegisterColorServiceServer(s, colorHandler)
			organizationpb.RegisterInternalNoteServiceServer(s, internalNoteHandler)
			// organizationpb.RegisterDealCostServiceServer(s, dealCostHandler)
			// organizationpb.RegisterDealContractServiceServer(s, dealContractHandler)
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
