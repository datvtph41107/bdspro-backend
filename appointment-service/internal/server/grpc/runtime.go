package grpc

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "pb/types/appointment"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/oklog/run"
	"github.com/tuandq2112/appointment-service/appconfig"
	"github.com/tuandq2112/appointment-service/internal/interface/handler"
	appointmentreminder "github.com/tuandq2112/appointment-service/internal/job/appointment_reminder"
	"github.com/tuandq2112/appointment-service/internal/server/grpc/interceptors"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	appointmentHandler     *handler.AppointmentHandler
	logger                 *flogging.FabricLogger
	appointmentReminderJob *appointmentreminder.AppointmentReminderJob
	// transactionHandler *handler.TransactionHandler
}

func NewGRPCServer(
	appointmentHandler *handler.AppointmentHandler,
	appointmentReminderJob *appointmentreminder.AppointmentReminderJob,
	// transactionHandler *handler.TransactionHandler,
) *GRPCServer {
	return &GRPCServer{
		appointmentHandler:     appointmentHandler,
		logger:                 flogging.MustGetLogger(appconfig.SERVICE_NAME),
		appointmentReminderJob: appointmentReminderJob,
		// transactionHandler: transactionHandler,
	}
}

func (s *GRPCServer) Start() {
	g := &run.Group{}

	start, stop := buildGRPCServer(s.logger, s.appointmentHandler)
	g.Add(start, stop)

	signalStart, signalStop := buildSignalHandler(s.logger)
	g.Add(signalStart, signalStop)

	go s.appointmentReminderJob.Run()
	if err := g.Run(); err != nil {
		s.logger.Errorf("Server exited: %v", err)
	}
}

func buildGRPCServer(logger *flogging.FabricLogger,
	appointmentHandler *handler.AppointmentHandler,
	// transactionHandler *handler.TransactionHandler,
) (func() error, func(error)) {
	unaryInts := []grpc.UnaryServerInterceptor{
		interceptors.UnaryLoggerInterceptor(logger),
		interceptors.UnaryAuthInterceptor(logger),
		interceptors.UnaryRecoveryInterceptor(logger),
	}
	streamInts := []grpc.StreamServerInterceptor{}
	logger.Info(fmt.Sprintf("Starting GRPC server at port %s", appconfig.GRPC_PORT))
	_, start, stop, err := BuildGRPCServer(
		// fmt.Sprintf("localhost:%s", appconfig.GRPC_PORT),
		fmt.Sprintf(":%s", appconfig.GRPC_PORT),
		func(s *grpc.Server) {
			logger.Info("Registering AppointmentServiceServer")
			pb.RegisterAppointmentServiceServer(s, appointmentHandler)
			// pbTransaction.RegisterTransactionServiceServer(s, transactionHandler)
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
		return server.Serve(lis)
	}

	stop = func(err error) {
		server.GracefulStop()
	}

	return server, start, stop, nil
}
