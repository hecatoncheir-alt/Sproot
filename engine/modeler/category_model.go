package modeler

import (
	"github.com/hecatoncheir/Sproot/engine/storage"
	"log"
)

func setCategoryModel(store *storage.Storage) {
	categoryForCreate := storage.Category{
		Name: "Смартфоны"}

	_, err := store.Categories.CreateCategory(categoryForCreate, "ru")
	if err != nil {
		log.Fatal(err)
	}
}
