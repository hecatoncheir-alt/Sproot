package engine

import (
	"encoding/json"
	"github.com/hecatoncheir/Sproot/configuration"
	"github.com/hecatoncheir/Sproot/engine/broker"
	"github.com/hecatoncheir/Sproot/engine/storage"
)

// Company is a part of configuration with methods for parse products
type Company struct {
	ID      string
	Storage *storage.Storage
	Broker  *broker.Broker

	Configuration configuration.Configuration
}

// CompanyData is a part of configuration for parse products
type CompanyData struct {
	ID   string
	Name string
	IRI  string
}

// CategoryData is a part of configuration for parse products
type CategoryData struct {
	ID   string
	Name string
}

// CityData is a part of configuration for parse products
type CityData struct {
	ID   string
	Name string
}

// InstructionOfCompany is a part of configuration for parse products
type InstructionOfCompany struct {
	Language        string
	Company         CompanyData
	Category        CategoryData
	City            CityData
	PageInstruction storage.PageInstruction
}

// ParseAll is a method of Company for get all instructions of company and send events for each
func (entity *Company) ParseAll(instructions []InstructionOfCompany) error {

	for _, instruction := range instructions {
		data, err := json.Marshal(instruction)
		if err != nil {
			return err
		}

		message := broker.EventData{
			Message: "Need products of category of company",
			Data:    string(data)}

		err = entity.Broker.WriteToTopic(entity.Configuration.APIVersion, message)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetInstructions is a method of Company for get all instruction of category for parse products
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

		if len(instruction.PagesInstruction) > 0 {
			inst.PageInstruction = instruction.PagesInstruction[0]
		}

		instOfCompany.Instructions = append(instOfCompany.Instructions, inst)
	}

	return instOfCompany.Instructions, nil
}
