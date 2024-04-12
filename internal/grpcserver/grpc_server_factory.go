package grpcserver

import (
	commonPB "github.com/quadev-ltd/qd-common/pb/gen/go/pb_email"
	"github.com/quadev-ltd/qd-common/pkg/grpcserver"
	"github.com/quadev-ltd/qd-common/pkg/log"
	commonTLS "github.com/quadev-ltd/qd-common/pkg/tls"
	"google.golang.org/grpc"

	"qd-email-api/internal/service"
)

// Factoryer is the interfact for creating a gRPC server
type Factoryer interface {
	Create(
		grpcServerAddress string,
		authenticationService service.EmailServicer,
		logFactory log.Factoryer,
		tlsEnabled bool,
	) (grpcserver.GRPCServicer, error)
}

// Factory is the implementation of the gRPC server factory
type Factory struct{}

var _ Factoryer = &Factory{}

// Create creates a gRPC server
func (grpcServerFactory *Factory) Create(
	grpcServerAddress string,
	emailService service.EmailServicer,
	logFactory log.Factoryer,
	tlsEnabled bool,
) (grpcserver.GRPCServicer, error) {
	// TODO: Set domain info in the config file
	const certFilePath = "certs/qd.email.api.crt"
	const keyFilePath = "certs/qd.email.api.key"
	// Create a listener for the gRPC server which eventually will start accepting connections when server is served
	grpcListener, err := commonTLS.CreateTLSListener(
		grpcServerAddress,
		certFilePath,
		keyFilePath,
		tlsEnabled,
	)
	if err != nil {
		return nil, err
	}

	// Create a gRPC server with a registered email service
	emailServiceGRPCServer := service.NewEmailServiceServer(emailService)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(log.CreateLoggerInterceptor(logFactory)),
	)
	commonPB.RegisterEmailServiceServer(grpcServer, emailServiceGRPCServer)

	return grpcserver.NewGRPCService(grpcServer, grpcListener), nil
}
