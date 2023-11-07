package grpcserver

import (
	"net"
	"qd_email_api/internal/service"
	"qd_email_api/pb/gen/go/pb_email"
	"qd_email_api/pkg/grpcserver"
	"qd_email_api/pkg/log"

	"google.golang.org/grpc"
)

// Factoryer is the interfact for creating a gRPC server
type Factoryer interface {
	Create(
		grpcServerAddress string,
		authenticationService service.EmailServicer,
		logFactory log.LogFactoryer,
	) (grpcserver.GRPCServicer, error)
}

// Factory is the implementation of the gRPC server factory
type Factory struct{}

var _ Factoryer = &Factory{}

// Create creates a gRPC server
func (grpcServerFactory *Factory) Create(
	grpcServerAddress string,
	emailService service.EmailServicer,
	logFactory log.LogFactoryer,
) (grpcserver.GRPCServicer, error) {
	// Create a gRPC server with a registered email service
	emailServiceGRPCServer := service.NewEmailServiceServer(emailService)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(log.CreateLoggerInterceptor(logFactory)),
	)
	pb_email.RegisterEmailServiceServer(grpcServer, emailServiceGRPCServer)
	// Create a listener for the gRPC server which eventually will start accepting connections when server is served
	grpcListener, err := net.Listen("tcp", grpcServerAddress)
	if err != nil {
		return nil, err
	}
	return grpcserver.NewGRPCService(grpcServer, grpcListener), nil
}
