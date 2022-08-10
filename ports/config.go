package ports

import (
	"github.com/mikew79/port-domain-service/application/entities"
)

// ConfigurationPort - Port for the domain service, provides access to configure the service
type ConfigurationPort interface {
	GetConfiguration() (entities.Configuration, error)
}
