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
		SprootTopic  string
		InitialTopic string

		Broker struct {
			Host string
			Port int
		}
		Database struct {
			Host string
			Port int
		}
	}

	Development struct {
		SprootTopic  string
		InitialTopic string

		Broker struct {
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
		configuration.APIVersion = "1.0.0"
	} else {
		configuration.APIVersion = apiVersion
	}

	productionSprootChannel := os.Getenv("Production-Sproot-Channel")
	if productionSprootChannel == "" {
		configuration.Production.SprootTopic = "Sproot"
	} else {
		configuration.Production.SprootTopic = productionSprootChannel
	}

	developmentSprootChannel := os.Getenv("Development-Sproot-Channel")
	if developmentSprootChannel == "" {
		configuration.Development.SprootTopic = "DevSproot"
	} else {
		configuration.Development.SprootTopic = developmentSprootChannel
	}

	productionInitialChannel := os.Getenv("Production-Initial-Channel")
	if productionSprootChannel == "" {
		configuration.Production.InitialTopic = "Initial"
	} else {
		configuration.Production.InitialTopic = productionInitialChannel
	}

	developmentInitialChannel := os.Getenv("Development-Initial-Channel")
	if developmentInitialChannel == "" {
		configuration.Development.InitialTopic = "DevInitial"
	} else {
		configuration.Development.InitialTopic = developmentInitialChannel
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
