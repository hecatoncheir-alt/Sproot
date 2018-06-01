package engine

import (
	"testing"

	"github.com/hecatoncheir/Configuration"
)

func TestIntegrationEngineCanBeSetUp(test *testing.T) {
	config := configuration.New()

	engine := New(config)
	err := engine.SetUpStorage(config.Development.Database.Host, config.Development.Database.Port)
	if err != nil {
		test.Error(err)
	}
}
