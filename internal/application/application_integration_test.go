package application

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/mhale/smtpd"
	commonPB "github.com/quadev-ltd/qd-common/pb/gen/go/pb_email"
	commonConfig "github.com/quadev-ltd/qd-common/pkg/config"
	commonLogger "github.com/quadev-ltd/qd-common/pkg/log"
	commonTLS "github.com/quadev-ltd/qd-common/pkg/tls"
	commonUtil "github.com/quadev-ltd/qd-common/pkg/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"qd-email-api/internal/config"
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

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	os.Setenv(commonConfig.AppEnvironmentKey, "test")

	// Save current working directory and change it
	originalWD, err := commonUtil.ChangeCurrentWorkingDirectory("../..")
	if err != nil {
		t.Fatalf("Failed to change working directory: %s", err)
	}
	// Defer the reset of the working directory
	defer os.Chdir(*originalWD)

	var config config.Config
	config.Load("internal/config")

	centralConfig := commonConfig.Config{
		TLSEnabled:                true,
		EmailVerificationEndpoint: "http://localhost:2222/",
		EmailService: commonConfig.Address{
			Host: "qd.email.api",
			Port: "1111",
		},
		AuthenticationService: commonConfig.Address{
			Host: "qd.authentication.api",
			Port: "3333",
		},
	}

	smtpServer := startMockSMTPServer(config.SMTP.Host, config.SMTP.Port)
	defer smtpServer.Close()

	application := NewApplication(&config, &centralConfig)
	go func() {
		application.StartServer()
	}()
	defer application.Close()

	waitForServerUp(application)

	t.Run("SendEmail_Success", func(t *testing.T) {
		connection, err := commonTLS.CreateGRPCConnection(application.GetGRPCServerAddress(), centralConfig.TLSEnabled)
		assert.NoError(t, err)

		client := commonPB.NewEmailServiceClient(connection)
		ctx := context.Background()
		sendEmailResponse, err := client.SendEmail(
			commonLogger.AddCorrelationIDToOutgoingContext(ctx, correlationID),
			&commonPB.SendEmailRequest{
				To:      email,
				Subject: subject,
				Body:    body,
			})

		assert.NoError(t, err)
		assert.Equal(t, "Email sent", sendEmailResponse.Message)
		assert.True(t, sendEmailResponse.Success)
	})

	t.Run("SendEmail_Email_Not_Sent_Error", func(t *testing.T) {
		connection, err := commonTLS.CreateGRPCConnection(application.GetGRPCServerAddress(), centralConfig.TLSEnabled)
		assert.NoError(t, err)

		client := commonPB.NewEmailServiceClient(connection)
		ctx := context.Background()
		registerResponse, err := client.SendEmail(
			commonLogger.AddCorrelationIDToOutgoingContext(ctx, correlationID),
			&commonPB.SendEmailRequest{
				To:      wrongEmail,
				Subject: subject,
				Body:    body,
			})

		assert.Error(t, err)
		assert.Equal(t, "rpc error: code = Internal desc = Error sending email", err.Error())
		assert.Nil(t, registerResponse)
	})

	t.Run("SendEmail_Email_Error_Missing_Correlation_ID", func(t *testing.T) {
		connection, err := commonTLS.CreateGRPCConnection(application.GetGRPCServerAddress(), centralConfig.TLSEnabled)
		assert.NoError(t, err)

		client := commonPB.NewEmailServiceClient(connection)
		ctx := context.Background()
		registerResponse, err := client.SendEmail(
			ctx,
			&commonPB.SendEmailRequest{
				To:      wrongEmail,
				Subject: subject,
				Body:    body,
			})

		assert.Error(t, err)
		assert.Equal(t, "rpc error: code = Unknown desc = Correlation ID not found in metadata", err.Error())
		assert.Nil(t, registerResponse)
	})
}
