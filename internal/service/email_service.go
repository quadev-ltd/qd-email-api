package service

import (
	"context"
	"fmt"
	"net/smtp"
)

// EmailServiceConfig constains the configuration for the email service
type EmailServiceConfig struct {
	From     string
	Password string
	Host     string
	Port     string
}

// EmailServicer is the interface for the email service
type EmailServicer interface {
	SendEmail(ctx context.Context, dest, subject, body string) error
}

// EmailService is the implementation of the email service
type EmailService struct {
	config EmailServiceConfig
	sender SMTPServicer
}

var _ EmailServicer = &EmailService{}

// NewEmailService creates a new email service
func NewEmailService(config EmailServiceConfig, sender SMTPServicer) *EmailService {
	return &EmailService{
		config: config,
		sender: &SMTPService{},
	}
}

// SendEmail sends an email to a single destination
func (service *EmailService) SendEmail(_ context.Context, dest, subject, body string) error {
	message := "From: " + service.config.From + "\n" +
		"To: " + dest + "\n" +
		"Subject: " + subject + "\n\n" +
		body
	config := service.config
	auth := smtp.PlainAuth("", config.From, config.Password, config.Host)
	resultError := smtp.SendMail(
		fmt.Sprintf("%s:%s", config.Host, config.Port),
		auth,
		config.From, []string{dest}, []byte(message))
	return resultError
}
