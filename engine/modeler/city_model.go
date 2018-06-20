package modeler

import (
	"github.com/hecatoncheir/Sproot/engine/storage"
	"log"
)

func setCityModel(store *storage.Storage) {
	cityForCreate := storage.City{
		Name: "Москва"}

	_, err := store.Cities.CreateCity(cityForCreate, "ru")
	if err != nil {
		log.Fatal(err)
	}
}
