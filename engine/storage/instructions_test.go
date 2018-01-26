package storage

import (
	"testing"
)

func TestIntegrationPageInstructionCanBeCreated(test *testing.T) {
	once.Do(prepareStorage)

	mVideoPageInstruction := PageInstruction{
		Path: "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector: ".pagination-list .pagination-item",
		PageParamPath:            "/f/page=",
		CityParamPath:            "?cityId=",
		CityParam:                "CityCZ_975",
		ItemSelector:             ".grid-view .product-tile",
		NameOfItemSelector:       ".product-tile-title",
		PriceOfItemSelector:      ".product-price-current"}

	createdPageInstruction, err := storage.Instructions.CreatePageInstruction(mVideoPageInstruction)
	if err != nil {
		test.Fail()
	}

	defer storage.Instructions.DeletePageInstruction(createdPageInstruction)

	if createdPageInstruction.ID == "" {
		test.Fail()
	}
}

func TestIntegrationPageInstructionCanBeReadById(test *testing.T) {
	once.Do(prepareStorage)
	mVideoPageInstruction := PageInstruction{
		Path: "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector: ".pagination-list .pagination-item",
		PageParamPath:            "/f/page=",
		CityParamPath:            "?cityId=",
		CityParam:                "CityCZ_975",
		ItemSelector:             ".grid-view .product-tile",
		NameOfItemSelector:       ".product-tile-title",
		PriceOfItemSelector:      ".product-price-current"}

	createdPageInstruction, err := storage.Instructions.CreatePageInstruction(mVideoPageInstruction)
	if err != nil {
		test.Fail()
	}

	defer storage.Instructions.DeletePageInstruction(createdPageInstruction)

	pageInstructionFromStore, err := storage.Instructions.ReadPageInstructionByID(createdPageInstruction.ID)
	if err != nil {
		test.Fail()
	}

	if createdPageInstruction.ID != pageInstructionFromStore.ID {
		test.Fail()
	}
}

func TestIntegrationPageInstructionCanBeDeleted(test *testing.T) {
	once.Do(prepareStorage)

	mVideoPageInstruction := PageInstruction{
		Path: "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector: ".pagination-list .pagination-item",
		PageParamPath:            "/f/page=",
		CityParamPath:            "?cityId=",
		CityParam:                "CityCZ_975",
		ItemSelector:             ".grid-view .product-tile",
		NameOfItemSelector:       ".product-tile-title",
		PriceOfItemSelector:      ".product-price-current"}

	createdPageInstruction, err := storage.Instructions.CreatePageInstruction(mVideoPageInstruction)
	if err != nil {
		test.Fail()
	}

	deletedPageInstructionID, err := storage.Instructions.DeletePageInstruction(createdPageInstruction)
	if err != nil {
		test.Error(err)
	}

	if deletedPageInstructionID != deletedPageInstructionID {
		test.Fail()
	}

	_, err = storage.Instructions.ReadPageInstructionByID(createdPageInstruction.ID)
	if err != ErrPageInstructionDoesNotExist {
		test.Error(err)
	}
}
