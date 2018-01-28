package engine

import "testing"

func TestIntegrationEngineCanBeSetUp(test *testing.T) {
	engine := New()

	err := engine.SetUpStorage(databaseHost, databasePort)
	if err != nil {
		test.Error(err)
	}
}
