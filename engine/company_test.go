package engine

import (
	"encoding/json"
	"testing"

	"github.com/hecatoncheir/Configuration"
	"github.com/hecatoncheir/Sproot/engine/storage"
)

func TestIntegrationCompanyCanGetInstructions(test *testing.T) {
	config := configuration.New()

	store := storage.New(config.Development.Database.Host, config.Development.Database.Port)
	err := store.SetUp()
	if err != nil {
		test.Fatalf(err.Error())
	}

	createdCompany, err := store.Companies.CreateCompany(storage.Company{Name: "Company test name"}, "en")
	defer func() {
		_, err := store.Companies.DeleteCompany(createdCompany)
		if err != nil {
			test.Error(err)
		}
	}()

	instruction, err := store.Instructions.CreateInstructionForCompany(createdCompany.ID, "en")
	defer func() {
		_, err := store.Instructions.DeleteInstruction(instruction)
		if err != nil {
			test.Error(err)
		}
	}()

	if err != nil {
		test.Error(err)
	}

	city, err := store.Cities.CreateCity(storage.City{Name: "Test city"}, "en")
	defer func() {
		_, err := store.Cities.DeleteCity(city)
		if err != nil {
			test.Error(err)
		}
	}()

	if err != nil {
		test.Error(err)
	}

	err = store.Instructions.AddCityToInstruction(instruction.ID, city.ID)
	if err != nil {
		test.Error(err)
	}

	mVideoPageInstruction := storage.PageInstruction{
		Path: "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector: ".pagination-list .pagination-item",
		PageParamPath:            "/f/page=",
		CityParamPath:            "?cityId=",
		ItemSelector:             ".grid-view .product-tile",
		NameOfItemSelector:       ".product-tile-title",
		PriceOfItemSelector:      ".product-price-current"}

	page, err := store.Instructions.CreatePageInstruction(mVideoPageInstruction)
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := store.Instructions.DeletePageInstruction(page)
		if err != nil {
			test.Error(err)
		}
	}()

	err = store.Instructions.AddPageInstructionToInstruction(instruction.ID, page.ID)
	if err != nil {
		test.Error(err)
	}

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
	err = json.Unmarshal(instructionsJSON, &inst)
	if err != nil {
		test.Error(err)
	}

	if inst.([]interface{})[0].(map[string]interface{})["Language"] != "en" {
		test.Fail()
	}

	if inst.([]interface{})[0].(map[string]interface{})["Company"].(map[string]interface{})["Name"] != "Company test name" {
		test.Fail()
	}
}
