package configuration

import (
	"os"
	"strconv"

	"log"
)

// Configuration is structure of config data from environment with default values
type Configuration struct {
	APIVersion string

	Production struct {
		Channel string
		Broker  struct {
			Host string
			Port int
		}
		Database struct {
			Host string
			Port int
		}
	}

	Development struct {
		Channel string
		Broker  struct {
			Host string
			Port int
		}
		Database struct {
			Host string
			Port int
		}
	}
}

// GetConfiguration function check environment variables and return structure with values
func GetConfiguration() (Configuration, error) {
	configuration := Configuration{}

	apiVersion := os.Getenv("API-Version")
	if apiVersion == "" {
		configuration.APIVersion = "v1"
	} else {
		configuration.APIVersion = apiVersion
	}

	productionParserChannel := os.Getenv("Production-Channel")
	if productionParserChannel == "" {
		configuration.Production.Channel = "Sproot"
	} else {
		configuration.Production.Channel = productionParserChannel
	}

	developmentParserChannel := os.Getenv("Development-Channel")
	if developmentParserChannel == "" {
		configuration.Development.Channel = "test"
	} else {
		configuration.Development.Channel = developmentParserChannel
	}

	productionBrokerHostFromEnvironment := os.Getenv("Production-Broker-Host")
	if productionBrokerHostFromEnvironment == "" {
		configuration.Production.Broker.Host = "192.168.99.100"
	} else {
		configuration.Production.Broker.Host = productionBrokerHostFromEnvironment
	}

	productionBrokerPortFromEnvironment := os.Getenv("Production-Broker-Port")
	if productionBrokerPortFromEnvironment == "" {
		configuration.Production.Broker.Port = 4150
	} else {
		port, err := strconv.Atoi(productionBrokerPortFromEnvironment)
		if err != nil {
			log.Fatal(err)
		}

		configuration.Production.Broker.Port = port
	}

	developmentBrokerHostFromEnvironment := os.Getenv("Development-Broker-Host")
	if developmentBrokerHostFromEnvironment == "" {
		configuration.Development.Broker.Host = "192.168.99.100"
	} else {
		configuration.Development.Broker.Host = developmentBrokerHostFromEnvironment
	}

	developmentBrokerPortFromEnvironment := os.Getenv("Development-Broker-Port")
	if developmentBrokerPortFromEnvironment == "" {
		configuration.Development.Broker.Port = 4150
	} else {
		port, err := strconv.Atoi(developmentBrokerPortFromEnvironment)
		if err != nil {
			log.Fatal(err)
		}

		configuration.Development.Broker.Port = port
	}

	productionDatabaseHostFromEnvironment := os.Getenv("Production-Database-Host")
	if productionDatabaseHostFromEnvironment == "" {
		configuration.Production.Database.Host = "192.168.99.100"
	} else {
		configuration.Production.Database.Host = productionDatabaseHostFromEnvironment
	}

	productionDatabasePortFromEnvironment := os.Getenv("Production-Database-Port")
	if productionDatabasePortFromEnvironment == "" {
		configuration.Production.Database.Port = 9080
	} else {
		port, err := strconv.Atoi(productionDatabasePortFromEnvironment)
		if err != nil {
			log.Fatal(err)
		}

		configuration.Production.Database.Port = port
	}

	developmentDatabaseHostFromEnvironment := os.Getenv("Development-Database-Host")
	if developmentDatabaseHostFromEnvironment == "" {
		configuration.Development.Database.Host = "192.168.99.100"
	} else {
		configuration.Development.Database.Host = developmentDatabaseHostFromEnvironment
	}

	developmentDatabasePortFromEnvironment := os.Getenv("Development-Database-Port")
	if developmentDatabasePortFromEnvironment == "" {
		configuration.Development.Database.Port = 9080
	} else {
		port, err := strconv.Atoi(developmentDatabasePortFromEnvironment)
		if err != nil {
			log.Fatal(err)
		}

		configuration.Development.Database.Port = port
	}

	return configuration, nil
}
