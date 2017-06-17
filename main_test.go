package main

import (
	"testing"

	"github.com/hecatoncheir/Sproot/engine"
)

func TestSprootCanSaveGetAndDeleteData(test *testing.T) {
	item, err := engine.SavePriceForProductOfCompany()
	if err != nil {
		test.Error(err)
	}
}
