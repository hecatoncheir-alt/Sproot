package storage

import (
	"context"
	"log"

	dataBaseAPI "github.com/dgraph-io/dgraph/protos/api"
)

// Product is a structure of products in database
type Product struct {
	ID               string      `json:"uid, omitempty"`
	Name             string      `json:"productName, omitempty"`
	IRI              string      `json:"productIri, omitempty"`
	PreviewImageLink string      `json:"previewImageLink, omitempty"`
	IsActive         bool        `json:"productsIsActive, omitempty"`
	Categories       []Category  `json:"belongs_to_category, omitempty"`
	Companies        []Companies `json:"belongs_to_company, omitempty"`
}

// Products is resource os storage for CRUD operations
type Products struct {
	storage *Storage
}

// SetUp is a method of Products resource for prepare database client and schema.
func (products *Products) SetUp() (err error) {
	schema := `
		productName: string @index(exact, term) .
		productIri: string @index(exact, term) .
		productImageLink: string @index(exact, term) .
		productIsActive: bool @index(bool) .
		belongs_to_category: uid .
		belongs_to_company: uid .
	`
	operation := &dataBaseAPI.Operation{Schema: schema}

	err = products.storage.Client.Alter(context.Background(), operation)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
