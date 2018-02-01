package engine

import (
	"time"

	"encoding/json"
	"github.com/hecatoncheir/Sproot/configuration"
	"github.com/hecatoncheir/Sproot/engine/broker"
	"github.com/hecatoncheir/Sproot/engine/storage"
)

type Company struct {
	ID      string
	Storage *storage.Storage
	Broker  *broker.Broker

	Configuration configuration.Configuration
}

type CompanyData struct {
	ID   string
	Name string
	IRI  string
}

type CategoryData struct {
	ID   string
	Name string
}

type CityData struct {
	ID   string
	Name string
}

type InstructionOfCompany struct {
	Language string
	Company  CompanyData
	Category CategoryData
	City     CityData
	Page     storage.PageInstruction
}

type PriceOfProduct struct {
	Value    string
	DateTime time.Time
}

type ProductOfCompany struct {
	Language string
	Name     string
	Price    PriceOfProduct
	Company  CompanyData
	Category CityData
}

// TODO
func (entity *Company) ParseAll(instructions []InstructionOfCompany) error {
	products, err := entity.Broker.ListenTopic(entity.Configuration.ApiVersion, entity.Configuration.Production.ParserChannel)
	if err != nil {
		return err
	}

	for _, instruction := range instructions {
		message := map[string]interface{}{"Message": "Parse products of company", "Data": instruction}
		go entity.Broker.WriteToTopic(entity.Configuration.ApiVersion, message)
	}

	go func() {
		for companyProduct := range products {
			var product ProductOfCompany
			json.Unmarshal(companyProduct, &product)

			// TODO break by timeout maybe
			break
		}
	}()

	return nil
}

func (entity *Company) GetInstructions() ([]InstructionOfCompany, error) {
	instructions, err := entity.Storage.Instructions.ReadAllInstructionsForCompany(entity.ID, ".")
	if err != nil {
		return nil, err
	}

	type instructionsOfCompany struct {
		Instructions []InstructionOfCompany `json:"instructions"`
	}

	instOfCompany := instructionsOfCompany{}

	for _, instruction := range instructions {

		inst := InstructionOfCompany{
			Language: instruction.Language,
			Company: CompanyData{
				ID:   instruction.Companies[0].ID,
				IRI:  instruction.Companies[0].IRI,
				Name: instruction.Companies[0].Name}}

		if len(instruction.Categories) > 0 {
			inst.Category = CategoryData{
				ID:   instruction.Categories[0].ID,
				Name: instruction.Categories[0].Name}
		}

		if len(instruction.Cities) > 0 {
			inst.City = CityData{
				ID:   instruction.Cities[0].ID,
				Name: instruction.Cities[0].Name}
		}

		if len(instruction.Pages) > 0 {
			inst.Page = instruction.Pages[0]
		}

		instOfCompany.Instructions = append(instOfCompany.Instructions, inst)
	}

	return instOfCompany.Instructions, nil
}
