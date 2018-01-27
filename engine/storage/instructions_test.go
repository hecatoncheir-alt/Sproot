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

	if deletedPageInstructionID != createdPageInstruction.ID {
		test.Fail()
	}

	_, err = storage.Instructions.ReadPageInstructionByID(createdPageInstruction.ID)
	if err != ErrPageInstructionDoesNotExist {
		test.Error(err)
	}
}

func TestIntegrationInstructionCanBeCreated(test *testing.T) {
	once.Do(prepareStorage)

	company, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeleteCompany(company)

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	defer storage.Instructions.DeleteInstruction(instruction)
	if err != nil {
		test.Fail()
	}

	if instruction.ID == "" {
		test.Fail()
	}

	if instruction.IsActive != true {
		test.Fail()
	}
}

func TestIntegrationInstructionCanBeDeleted(test *testing.T) {
	once.Do(prepareStorage)

	company, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeleteCompany(company)

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")

	deletedInstructionID, err := storage.Instructions.DeleteInstruction(instruction)
	if err != nil {
		test.Error(err)
	}

	if deletedInstructionID != instruction.ID {
		test.Fail()
	}

	_, err = storage.Instructions.ReadInstructionByID(instruction.ID, "en")
	if err != ErrInstructionDoesNotExist {
		test.Error(err)
	}
}
