package engine

import (
	"encoding/json"
	"github.com/hecatoncheir/Sproot/engine/storage"
	"testing"
)

func TestIntegrationCompanyCanGetInstructions(test *testing.T) {
	store := storage.New(databaseHost, databasePort)
	store.SetUp()

	createdCompany, _ := store.Companies.CreateCompany(storage.Company{Name: "Test company"}, "en")
	defer store.Companies.DeleteCompany(createdCompany)

	instruction, _ := store.Instructions.CreateInstructionForCompany(createdCompany.ID, "en")
	defer store.Instructions.DeleteInstruction(instruction)

	city, _ := store.Cities.CreateCity(storage.City{Name: "Test city"}, "en")
	defer store.Cities.DeleteCity(city)

	store.Instructions.AddCityToInstruction(instruction.ID, city.ID)

	mVideoPageInstruction := storage.PageInstruction{
		Path: "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector: ".pagination-list .pagination-item",
		PageParamPath:            "/f/page=",
		CityParamPath:            "?cityId=",
		CityParam:                "CityCZ_975",
		ItemSelector:             ".grid-view .product-tile",
		NameOfItemSelector:       ".product-tile-title",
		PriceOfItemSelector:      ".product-price-current"}

	page, _ := store.Instructions.CreatePageInstruction(mVideoPageInstruction)
	defer store.Instructions.DeletePageInstruction(page)

	store.Instructions.AddPageInstructionToInstruction(instruction.ID, page.ID)

	company := Company{ID: createdCompany.ID, Storage: store}

	instructionsJSON, err := company.GetJSONInstructions()
	if err != nil {
		test.Error(err)
	}

	if instructionsJSON == "" {
		test.Fail()
	}

	var inst interface{}

	json.Unmarshal([]byte(instructionsJSON), &inst)

	if inst.(map[string]interface{})["instructions"].([]interface{})[0].(map[string]interface{})["language"] != "en" {
		test.Fail()
	}

	if inst.(map[string]interface{})["instructions"].([]interface{})[0].(map[string]interface{})["company"].(map[string]interface{})["name"] != "Test company" {
		test.Fail()
	}
}
