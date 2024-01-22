package service

import (
	"context"

	"github.com/quadev-ltd/qd-common/pkg/log"
	"golang.org/x/time/rate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"qd-email-api/pb/gen/go/pb_email"
)

// EmailServiceServer is the implementation of the authentication service
type EmailServiceServer struct {
	emailService EmailServicer
	limitter     *rate.Limiter
	pb_email.UnimplementedEmailServiceServer
}

var _ pb_email.EmailServiceServer = &EmailServiceServer{}

// NewEmailServiceServer creates a new authentication service
func NewEmailServiceServer(emailService EmailServicer) *EmailServiceServer {
	return &EmailServiceServer{
		emailService: emailService,
	}
}

// SendEmail sends an email
func (server *EmailServiceServer) SendEmail(ctx context.Context, request *pb_email.SendEmailRequest) (*pb_email.SendEmailResponse, error) {
	logger := log.GetLoggerFromContext(ctx)
	if logger == nil {
		return nil, status.Errorf(codes.Internal, "No logger in context")
	}

	// Check the rate limit
	if !limiter.Allow() {
		logger.Error(nil, "Too many requests")
		return nil, status.Errorf(codes.ResourceExhausted, "Too many requests")
	}

	// Send the email
	err := server.emailService.SendEmail(ctx, request.To, request.Subject, request.Body)
	if err != nil {
		logger.Error(err, "Error sending email")
		return nil, status.Errorf(codes.Internal, "Error sending email")
	}

	logger.Info("Email sent")
	return &pb_email.SendEmailResponse{
		Success: true,
		Message: "Email sent",
	}, nil
}

// TODO: inject this into the struct (dependency injection) and add the unit tests
var limiter = rate.NewLimiter(rate.Limit(1), 5)
