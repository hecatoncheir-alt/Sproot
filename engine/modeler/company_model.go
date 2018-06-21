package modeler

import (
	"github.com/hecatoncheir/Sproot/engine/storage"
	"log"
)

func setCompanyModel(store *storage.Storage) {
	companyForCreate := storage.Company{
		IRI:  "http://www.mvideo.ru/",
		Name: "М.Видео"}

	_, err := store.Companies.CreateCompany(companyForCreate, "ru")
	if err != nil && err != storage.ErrCompanyAlreadyExist {
		log.Fatal(err)
	}
}
