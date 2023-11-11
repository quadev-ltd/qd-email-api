package application

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	pkgConfig "github.com/gustavo-m-franco/qd-common/pkg/config"
	pkgLogger "github.com/gustavo-m-franco/qd-common/pkg/log"
	"github.com/mhale/smtpd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"qd-email-api/internal/config"
	pb_email "qd-email-api/pb/gen/go/pb_email"
)

func isServerUp(addr string) bool {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func waitForServerUp(application Applicationer) {
	maxWaitTime := 10 * time.Second
	startTime := time.Now()

	for {
		if time.Since(startTime) > maxWaitTime {
			log.Error().Msg("Server didn't start within the specified time")
		}

		if isServerUp(application.GetGRPCServerAddress()) {
			log.Error().Msg("Server is up")
			break
		}

		time.Sleep(1 * time.Second)
	}
}

const wrongEmail = "wrong@email.com"

func startMockSMTPServer(mockSMTPServerHost string, mockSMTPServerPort string) *smtpd.Server {
	authMechanisms := map[string]bool{
		"PLAIN": true,
	}
	smtpServer := smtpd.Server{
		Addr:     fmt.Sprintf("%s:%s", mockSMTPServerHost, mockSMTPServerPort),
		Appname:  "Mock SMTP Server",
		Hostname: mockSMTPServerPort,
		Handler: func(remoteAddress net.Addr, from string, to []string, data []byte) error {
			if to[0] == wrongEmail {
				return fmt.Errorf("Invalid email address")
			}
			return nil
		},
		AuthHandler: func(remoteAddress net.Addr, mechanism string, username []byte, password []byte, shared []byte) (bool, error) {
			return true, nil
		},
		AuthRequired: true,
		AuthMechs:    authMechanisms,
	}

	go func() {
		log.Info().Msg(fmt.Sprintf("Starting mock SMTP server %s... ", fmt.Sprintf("%s:%s", mockSMTPServerHost, mockSMTPServerPort)))
		err := smtpServer.ListenAndServe()
		if err != nil {
			log.Err(err)
		}
	}()
	return &smtpServer
}

func TestEmailMicroService(t *testing.T) {
	email := "test@test.com"
	subject := "Test Subject"
	body := "Test Body"
	correlationID := "1234567890"

	zerolog.SetGlobalLevel(zerolog.Disabled)
	os.Setenv(pkgConfig.AppEnvironmentKey, "test")

	var config config.Config
	config.Load("../../internal/config")
	config.SMTP.Port = "9999"

	smtpServer := startMockSMTPServer(config.SMTP.Host, config.SMTP.Port)
	defer smtpServer.Close()

	application := NewApplication(&config)
	go func() {
		application.StartServer()
	}()
	defer application.Close()

	waitForServerUp(application)

	t.Run("SendEmail_Success", func(t *testing.T) {
		connection, err := grpc.Dial(application.GetGRPCServerAddress(), grpc.WithInsecure())
		assert.NoError(t, err)

		client := pb_email.NewEmailServiceClient(connection)
		ctx := context.Background()
		registerResponse, err := client.SendEmail(
			pkgLogger.AddCorrelationIDToContext(ctx, correlationID),
			&pb_email.SendEmailRequest{
				To:      email,
				Subject: subject,
				Body:    body,
			})

		assert.NoError(t, err)
		assert.Equal(t, "Email sent", registerResponse.Message)
		assert.True(t, registerResponse.Success)
	})

	t.Run("SendEmail_Email_Not_Sent_Error", func(t *testing.T) {
		connection, err := grpc.Dial(application.GetGRPCServerAddress(), grpc.WithInsecure())
		assert.NoError(t, err)

		client := pb_email.NewEmailServiceClient(connection)
		ctx := context.Background()
		registerResponse, err := client.SendEmail(
			pkgLogger.AddCorrelationIDToContext(ctx, correlationID),
			&pb_email.SendEmailRequest{
				To:      wrongEmail,
				Subject: subject,
				Body:    body,
			})

		assert.Error(t, err)
		assert.Equal(t, "rpc error: code = Internal desc = Error sending email", err.Error())
		assert.Nil(t, registerResponse)
	})

	t.Run("SendEmail_Email_Error_Missing_Correlation_ID", func(t *testing.T) {
		connection, err := grpc.Dial(application.GetGRPCServerAddress(), grpc.WithInsecure())
		assert.NoError(t, err)

		client := pb_email.NewEmailServiceClient(connection)
		ctx := context.Background()
		registerResponse, err := client.SendEmail(
			ctx,
			&pb_email.SendEmailRequest{
				To:      wrongEmail,
				Subject: subject,
				Body:    body,
			})

		assert.Error(t, err)
		assert.Equal(t, "rpc error: code = Internal desc = Internal server error. Dubious request", err.Error())
		assert.Nil(t, registerResponse)
	})
}
