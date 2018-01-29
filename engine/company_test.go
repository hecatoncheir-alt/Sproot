package engine

import (
	"encoding/json"
	"testing"

	"github.com/hecatoncheir/Sproot/configuration"
	"github.com/hecatoncheir/Sproot/engine/storage"
)

func TestIntegrationCompanyCanGetInstructions(test *testing.T) {
	config, err := configuration.GetConfiguration()
	if err != nil {
		test.Error(err)
	}

	store := storage.New(config.Development.Database.Host, config.Development.Database.Port)
	store.SetUp()

	createdCompany, err := store.Companies.CreateCompany(storage.Company{Name: "Company test name"}, "en")
	defer store.Companies.DeleteCompany(createdCompany)
	if err != nil {
		test.Error(err)
	}

	instruction, err := store.Instructions.CreateInstructionForCompany(createdCompany.ID, "en")
	defer store.Instructions.DeleteInstruction(instruction)
	if err != nil {
		test.Error(err)
	}

	city, err := store.Cities.CreateCity(storage.City{Name: "Test city"}, "en")
	defer store.Cities.DeleteCity(city)
	if err != nil {
		test.Error(err)
	}

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

	instructions, err := company.GetInstructions()
	if err != nil {
		test.Error(err)
	}

	instructionsJSON, err := json.Marshal(instructions)
	if err != nil {
		test.Error(err)
	}

	if string(instructionsJSON) == "" {
		test.Fail()
	}

	var inst interface{}
	json.Unmarshal(instructionsJSON, &inst)

	if inst.([]interface{})[0].(map[string]interface{})["language"] != "en" {
		test.Fail()
	}

	if inst.([]interface{})[0].(map[string]interface{})["company"].(map[string]interface{})["name"] != "Company test name" {
		test.Fail()
	}
}
