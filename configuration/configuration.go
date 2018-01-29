package configuration

import (
	"github.com/prometheus/common/log"
	"os"
	"strconv"
)

type Configuration struct {
	ApiVersion string

	Production struct {
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

func GetConfiguration() (Configuration, error) {
	configuration := Configuration{}

	apiVersion := os.Getenv("Api-Version")
	if apiVersion == "" {
		configuration.ApiVersion = "v1"
	} else {
		configuration.ApiVersion = apiVersion
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
