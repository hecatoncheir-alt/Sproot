package modeler

import (
	"github.com/hecatoncheir/Sproot/engine/storage"
	"log"
)

type Modeler struct {
	Storage *storage.Storage
}

func New(storage *storage.Storage) *Modeler {
	return &Modeler{Storage: storage}
}

func (modeler *Modeler) SetUpAll() {
	companyForCreate := storage.Company{
		IRI:  "http://www.mvideo.ru/",
		Name: "М.Видео"}

	createdCompany, err := modeler.Storage.Companies.CreateCompany(companyForCreate, "ru")
	if err != nil && err != storage.ErrCompanyAlreadyExist {
		log.Fatal(err)
	}

	categoryForCreate := storage.Category{
		Name: "Смартфоны"}

	createdCategory, err := modeler.Storage.Categories.CreateCategory(categoryForCreate, "ru")
	if err != nil && err != storage.ErrCategoryAlreadyExist {
		log.Fatal(err)
	}

	err = modeler.Storage.Categories.AddCompanyToCategory(createdCategory.ID, createdCompany.ID)
	if err != nil {
		log.Fatal(err)
	}

	createdInstruction, err := modeler.Storage.Instructions.CreateInstructionForCompany(createdCompany.ID, "ru")
	if err != nil {
		log.Println(err)
	}

	err = modeler.Storage.Instructions.AddCategoryToInstruction(createdInstruction.ID, createdCategory.ID)
	if err != nil {
		log.Println(err)
	}

	cityForCreate := storage.City{
		Name: "Москва"}

	createdCity, err := modeler.Storage.Cities.CreateCity(cityForCreate, "ru")
	if err != nil && err != storage.ErrCityAlreadyExist {
		log.Fatal(err)
	}

	err = modeler.Storage.Instructions.AddCityToInstruction(createdInstruction.ID, createdCity.ID)
	if err != nil {
		log.Println(err)
	}

	pageInstructionForCreate := storage.PageInstruction{
		Path: "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector:   ".c-pagination > .c-pagination__num",
		PageParamPath:              "/f/page=",
		CityParamPath:              "?cityId=",
		ItemSelector:               ".c-product-tile",
		PreviewImageOfItemSelector: ".c-product-tile-picture__link .lazy-load-image-holder img",
		NameOfItemSelector:         ".c-product-tile__description .sel-product-tile-title",
		LinkOfItemSelector:         ".c-product-tile__description .sel-product-tile-title",
		PriceOfItemSelector:        ".c-product-tile__checkout-section .c-pdp-price__current"}

	createdPageInstruction, err := modeler.Storage.Instructions.CreatePageInstruction(pageInstructionForCreate)
	if err != nil {
		log.Println(err)
	}

	err = modeler.Storage.Instructions.AddPageInstructionToInstruction(createdInstruction.ID, createdPageInstruction.ID)
	if err != nil {
		log.Println(err)
	}

}
