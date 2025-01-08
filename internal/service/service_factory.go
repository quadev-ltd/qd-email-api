package service

import (
	commonConfig "github.com/quadev-ltd/qd-common/pkg/config"

	"qd-email-api/internal/config"
)

// Factoryer is a factory for creating a service
type Factoryer interface {
	CreateService(config *config.Config, centralConfig *commonConfig.Config) (EmailServicer, error)
}

// Factory is the implementation of the service factory
type Factory struct{}

var _ Factoryer = &Factory{}

// CreateService creates a service
func (serviceFactory *Factory) CreateService(
	config *config.Config,
	centralConfig *commonConfig.Config,
) (EmailServicer, error) {

	emailServiceConfig := EmailServiceConfig{
		From:     config.SMTP.From,
		Domain:   config.SMTP.Domain,
		AppName:  centralConfig.AppName,
		Username: config.SMTP.Username,
		Password: config.SMTP.Password,
		Host:     config.SMTP.Host,
		Port:     config.SMTP.Port,
	}
	return NewEmailService(emailServiceConfig, &SMTPService{}), nil
}
