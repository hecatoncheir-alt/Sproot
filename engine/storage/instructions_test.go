package storage

import (
	"testing"
)

func TestIntegrationInstructionCanBeCreated(test *testing.T) {
	once.Do(prepareStorage)
}
