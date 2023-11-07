package application

import (
	"fmt"
	"qd_email_api/internal/config"
	grpcFactory "qd_email_api/internal/grpcserver"
	"qd_email_api/internal/service"
	"qd_email_api/pkg/grpcserver"

	"qd_email_api/pkg/log"
)

// Applicationer provides the main functions to start the application
type Applicationer interface {
	StartServer()
	Close()
	GetGRPCServerAddress() string
}

// Application is the main application
type Application struct {
	logger            log.Loggerer
	grpcServiceServer grpcserver.GRPCServicer
	grpcServerAddress string
	service           service.EmailServicer
}

// NewApplication creates a new application
func NewApplication(config *config.Config) Applicationer {
	logFactory := log.NewLogFactory(config.Environment)
	logger := logFactory.NewLogger()

	emailService, err := (&service.Factory{}).CreateService(config)
	if err != nil {
		logger.Error(err, "Failed to create email service")
	}

	grpcServerAddress := fmt.Sprintf("%s:%s", config.GRPC.Host, config.GRPC.Port)
	grpcServiceServer, err := (&grpcFactory.Factory{}).Create(
		grpcServerAddress,
		emailService,
		logFactory,
	)
	if err != nil {
		logger.Error(err, "Failed to create grpc server: %v")
	}

	return New(grpcServiceServer, grpcServerAddress, emailService, logger)
}

func New(
	grpcServiceServer grpcserver.GRPCServicer,
	grpcServerAddress string,
	service service.EmailServicer,
	logger log.Loggerer,
) Applicationer {
	return &Application{
		grpcServiceServer: grpcServiceServer,
		grpcServerAddress: grpcServerAddress,
		service:           service,
		logger:            logger,
	}
}

// StartServer starts the gRPC server
func (application *Application) StartServer() {
	application.logger.Info(fmt.Sprintf("Starting gRPC server on %s:...", application.grpcServerAddress))
	err := application.grpcServiceServer.Serve()
	if err != nil {
		application.logger.Error(err, "Failed to serve grpc server")
		return
	}
}

// Close closes the gRPC server and services used by the application
func (application *Application) Close() {
	switch {
	case application.service == nil:
		application.logger.Error(nil, "Service is not created")
		return
	case application.grpcServiceServer == nil:
		application.logger.Error(nil, "gRPC server is not created")
		return
	}
	application.grpcServiceServer.Close()
	application.logger.Info("gRPC server closed")
}

// GetGRPCServerAddress returns the gRPC server address
func (application *Application) GetGRPCServerAddress() string {
	return application.grpcServerAddress
}
