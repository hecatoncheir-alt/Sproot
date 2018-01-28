package engine

import (
	"github.com/hecatoncheir/Sproot/configuration"
	"testing"
)

func TestIntegrationEngineCanBeSetUp(test *testing.T) {
	config, err := configuration.GetConfiguration()
	if err != nil {
		test.Error(err)
	}

	engine := New()
	err = engine.SetUpStorage(config.Development.Database.Host, config.Development.Database.Port)
	if err != nil {
		test.Error(err)
	}
}
