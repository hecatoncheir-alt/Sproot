package modeler

import "github.com/hecatoncheir/Sproot/engine/storage"

type Modeler struct {
	Storage *storage.Storage
}

func New(storage *storage.Storage) *Modeler {
	return &Modeler{Storage: storage}
}

func (modeler *Modeler) SetUpAll() {
	setCompanyModel(modeler.Storage)
}
