package service

import (
	"context"
	"errors"
	"qd-email-api/internal/service/mock"
	"qd-email-api/pb/gen/go/pb_email"
	"testing"

	"github.com/gustavo-m-franco/qd-common/pkg/log"

	loggerMock "github.com/gustavo-m-franco/qd-common/pkg/log/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEmailServiceServer(test *testing.T) {
	sendEmailRequest := &pb_email.SendEmailRequest{
		To:      "test@test.com",
		Subject: "Test subject",
		Body:    "Test body",
	}
	test.Run("Send_Email_Error_No_Logger", func(test *testing.T) {
		controller := gomock.NewController(test)
		defer controller.Finish()

		emailServiceMock := mock.NewMockEmailServicer(controller)
		// loggerMock := loggerMock.NewMockLoggerer(controller)
		// ctx := context.WithValue(context.Background(), log.LoggerKey, loggerMock)

		server := NewEmailServiceServer(emailServiceMock)

		response, returnedError := server.SendEmail(context.Background(), sendEmailRequest)

		assert.Error(test, returnedError)
		assert.Equal(test, "rpc error: code = Internal desc = Internal server error. No logger in context", returnedError.Error())
		assert.Nil(test, response)
	})

	test.Run("Send_Email_Error_Send_Email", func(test *testing.T) {
		controller := gomock.NewController(test)
		defer controller.Finish()

		emailServiceMock := mock.NewMockEmailServicer(controller)
		loggerMock := loggerMock.NewMockLoggerer(controller)
		ctx := context.WithValue(context.Background(), log.LoggerKey, loggerMock)

		server := NewEmailServiceServer(emailServiceMock)

		const expectedError = "Error sending email"
		loggerMock.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)
		emailServiceMock.EXPECT().SendEmail(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Times(1).Return(errors.New(expectedError))

		response, returnedError := server.SendEmail(ctx, sendEmailRequest)

		assert.Error(test, returnedError)
		assert.Equal(test, "rpc error: code = Internal desc = Error sending email", returnedError.Error())
		assert.Nil(test, response)
	})

	test.Run("Send_Email_Success", func(test *testing.T) {
		controller := gomock.NewController(test)
		defer controller.Finish()

		emailServiceMock := mock.NewMockEmailServicer(controller)
		loggerMock := loggerMock.NewMockLoggerer(controller)
		ctx := context.WithValue(context.Background(), log.LoggerKey, loggerMock)

		server := NewEmailServiceServer(emailServiceMock)

		loggerMock.EXPECT().Info(gomock.Any()).Times(1)
		emailServiceMock.EXPECT().SendEmail(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Times(1).Return(nil)

		response, returnedError := server.SendEmail(ctx, sendEmailRequest)

		assert.NoError(test, returnedError)
		assert.True(test, response.Success)
		assert.Equal(test, "Email sent", response.Message)
	})

}
