package configuration

import (
	"github.com/prometheus/common/log"
	"os"
	"strconv"
)

type Configuration struct {
	Production struct {
		Database struct {
			Host string
			Port int
		}
	}

	Development struct {
		Database struct {
			Host string
			Port int
		}
	}
}

func GetConfiguration() (Configuration, error) {
	configuration := Configuration{}

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
