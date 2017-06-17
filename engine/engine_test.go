package engine

import (
	"testing"
)

func TestSprootCanSaveGetAndDeleteData(test *testing.T) {

	item, err := SavePriceForProductOfCompany()

	if err != nil {
		test.Error(err)
	}
}
