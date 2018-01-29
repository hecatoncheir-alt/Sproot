package engine

import (
	"github.com/hecatoncheir/Sproot/engine/broker"
	"github.com/hecatoncheir/Sproot/engine/storage"
)

type Company struct {
	ID      string
	Storage *storage.Storage
	Broker  *broker.Broker
}

type company struct {
	ID   string `json:"id, omitempty"`
	Name string `json:"name, omitempty"`
	IRI  string `json:"iri, omitempty"`
}

type category struct {
	ID   string `json:"id, omitempty"`
	Name string `json:"name, omitempty"`
}

type city struct {
	ID   string `json:"id, omitempty"`
	Name string `json:"name, omitempty"`
}

type instructionOfCompany struct {
	Language string                  `json:"language, omitempty"`
	Company  company                 `json:"company, omitempty"`
	Category category                `json:"category, omitempty"`
	City     city                    `json:"city, omitempty"`
	Page     storage.PageInstruction `json:"page, omitempty"`
}

type instructionsOfCompany struct {
	Instructions []instructionOfCompany `json:"instructions"`
}

// TODO
func (entity *Company) ParseAll() error {
	//instructions, err := entity.GetInstructions()
	//if err != nil {
	//	return err
	//}

	//inst, err := json.Marshal(instructions)
	//if err != nil {
	//	return "", err
	//}
	return nil
}

func (entity *Company) GetInstructions() ([]instructionOfCompany, error) {
	instructions, err := entity.Storage.Instructions.ReadAllInstructionsForCompany(entity.ID, ".")
	if err != nil {
		return nil, err
	}

	instOfCompany := instructionsOfCompany{}

	for _, instruction := range instructions {

		inst := instructionOfCompany{
			Language: instruction.Language,
			Company: company{
				ID:   instruction.Companies[0].ID,
				IRI:  instruction.Companies[0].IRI,
				Name: instruction.Companies[0].Name}}

		if len(instruction.Categories) > 0 {
			inst.Category = category{
				ID:   instruction.Categories[0].ID,
				Name: instruction.Categories[0].Name}
		}

		if len(instruction.Cities) > 0 {
			inst.City = city{
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
