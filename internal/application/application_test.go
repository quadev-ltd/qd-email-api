package application

import (
	"errors"
	"qd_email_api/internal/service/mock"
	grpcserverMock "qd_email_api/pkg/grpcserver/mock"
	loggerMock "qd_email_api/pkg/log/mock"
	"testing"

	"github.com/golang/mock/gomock"
)

func setupApplication(t *testing.T, useEmailService, useGRPCServer bool) (Applicationer, *gomock.Controller, *grpcserverMock.MockGRPCServicer, *mock.MockEmailServicer, *loggerMock.MockLoggerer) {
	controller := gomock.NewController(t)
	var application Applicationer

	emailServiceMock := mock.NewMockEmailServicer(controller)
	grpcServiceServerMock := grpcserverMock.NewMockGRPCServicer(controller)
	loggerMock := loggerMock.NewMockLoggerer(controller)
	grpcAddres := "localhost:8080"

	switch {
	case useEmailService && useGRPCServer:
		application = New(grpcServiceServerMock, grpcAddres, emailServiceMock, loggerMock)
	case !useEmailService:
		application = New(grpcServiceServerMock, grpcAddres, nil, loggerMock)
	case !useGRPCServer:
		application = New(nil, grpcAddres, emailServiceMock, loggerMock)
	}

	return application, controller, grpcServiceServerMock, emailServiceMock, loggerMock
}

func TestApplication(t *testing.T) {

	t.Run("Serve_Error", func(t *testing.T) {
		application, controller, grpcServiceServerMock, _, loggerMock := setupApplication(t, true, true)
		defer controller.Finish()

		expectedError := errors.New("Error sending email")
		grpcServiceServerMock.EXPECT().Serve().Return(expectedError)
		loggerMock.EXPECT().Info("Starting gRPC server on localhost:8080:...").Times(1)
		loggerMock.EXPECT().Error(expectedError, "Failed to serve grpc server").Times(1)

		application.StartServer()
	})
	t.Run("Serve_Success", func(t *testing.T) {
		application, controller, grpcServiceServerMock, _, loggerMock := setupApplication(t, true, true)
		defer controller.Finish()

		grpcServiceServerMock.EXPECT().Serve().Times(1).Return(nil)
		loggerMock.EXPECT().Info("Starting gRPC server on localhost:8080:...").Times(1)

		application.StartServer()
	})

	t.Run("Close_No_Service_Error", func(t *testing.T) {
		application, controller, _, _, loggerMock := setupApplication(t, false, true)
		defer controller.Finish()

		loggerMock.EXPECT().Error(nil, "Service is not created").Times(1)

		application.Close()
	})

	t.Run("Close_No_GRPC_Server_Error", func(t *testing.T) {
		application, controller, _, _, loggerMock := setupApplication(t, true, false)
		defer controller.Finish()

		loggerMock.EXPECT().Error(nil, "gRPC server is not created").Times(1)

		application.Close()
	})

	t.Run("Close_Success", func(t *testing.T) {
		application, controller, grpcServiceServerMock, _, loggerMock := setupApplication(t, true, true)
		defer controller.Finish()

		grpcServiceServerMock.EXPECT().Close().Times(1)
		loggerMock.EXPECT().Info("gRPC server closed").Times(1)

		application.Close()
	})
}
