package engine

import (
	"encoding/json"
	"testing"

	"github.com/hecatoncheir/Sproot/configuration"
	"github.com/hecatoncheir/Sproot/engine/broker"
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

func TestIntegrationCompanyCanSendInstructionsToBroker(test *testing.T) {
	instructionsJSON := `[
	   {  
		  "language":"en",
		  "company":{  
			 "id":"0x2786",
			 "name":"Company test name",
			 "iri":""
		  },
		  "category":{  
			 "id":"",
			 "name":""
		  },
		  "city":{  
			 "id":"0x2788",
			 "name":"Test city"
		  },
		  "page":{  
			 "uid":"0x2789",
			 "path":"smartfony-i-svyaz/smartfony-205",
			 "pageInPaginationSelector":".pagination-list .pagination-item",
			 "previewImageOfSelector":"",
			 "pageParamPath":"/f/page=",
			 "cityParamPath":"?cityId=",
			 "cityParam":"CityCZ_975",
			 "itemSelector":".grid-view .product-tile",
			 "nameOfItemSelector":".product-tile-title",
			 "cityInCookieKey":"",
			 "cityIdForCookie":"",
			 "priceOfItemSelector":".product-price-current"
		  }
	   }
	]`

	type instructionsOfCompanyForParseAllProductsFromCategory struct {
		instructions []InstructionOfCompany
	}

	var decodedInstructions instructionsOfCompanyForParseAllProductsFromCategory

	json.Unmarshal([]byte(instructionsJSON), &decodedInstructions)

	config, err := configuration.GetConfiguration()
	if err != nil {
		test.Error(err)
	}

	bro := broker.New()
	err = bro.Connect(config.Development.Broker.Host, config.Development.Broker.Port)
	if err != nil {
		test.Error(err)
	}

	company := Company{Broker: bro, ApiVersion: config.ApiVersion}

	items, err := bro.ListenTopic(config.ApiVersion, "Crawler")
	if err != nil {
		test.Error(err)
	}

	go company.ParseAll(decodedInstructions.instructions)

	for item := range items {
		data := map[string]string{}
		json.Unmarshal(item, &data)
		if data["test key"] == "test value" {
			break
		}
	}
}
