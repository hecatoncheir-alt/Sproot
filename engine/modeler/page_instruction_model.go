package modeler

import (
	"github.com/hecatoncheir/Sproot/engine/storage"
	"log"
)

func setPageInstructionModel(store *storage.Storage) {
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

	_, err := store.Instructions.CreatePageInstruction(pageInstructionForCreate)
	if err != nil {
		log.Fatal(err)
	}
}
