package engine

import (
	"github.com/hecatoncheir/Sproot/engine/broker"
	"github.com/hecatoncheir/Sproot/engine/storage"
)

type Company struct {
	ID         string
	Storage    *storage.Storage
	Broker     *broker.Broker
	ApiVersion string
}

type CompanyData struct {
	ID   string `json:"id, omitempty"`
	Name string `json:"name, omitempty"`
	IRI  string `json:"iri, omitempty"`
}

type CategoryData struct {
	ID   string `json:"id, omitempty"`
	Name string `json:"name, omitempty"`
}

type CityData struct {
	ID   string `json:"id, omitempty"`
	Name string `json:"name, omitempty"`
}

type InstructionOfCompany struct {
	Language string                  `json:"language, omitempty"`
	Company  CompanyData             `json:"company, omitempty"`
	Category CategoryData            `json:"category, omitempty"`
	City     CityData                `json:"city, omitempty"`
	Page     storage.PageInstruction `json:"page, omitempty"`
}

func (entity *Company) ParseAll([]InstructionOfCompany) error {
	message := map[string]string{"test key": "test value"}
	go entity.Broker.WriteToTopic(entity.ApiVersion,message)
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
