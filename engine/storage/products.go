package storage

import (
	"context"
	"log"

	"encoding/json"
	"errors"
	"fmt"

	dataBaseAPI "github.com/dgraph-io/dgraph/protos/api"
)

var (
	// ErrProductCanNotBeWithoutID means that product can't be found in storage for make some operation
	ErrProductCanNotBeWithoutID = errors.New("product can not be without id")

	// ErrProductAlreadyExist means that the product is in the database already
	ErrProductAlreadyExist = errors.New("product already exist")

	// ErrProductsByNameNotFound means than the products does not exist in database
	ErrProductsByNameNotFound = errors.New("products by name not found")

	// ErrProductsByNameCanNotBeFound means that the products can't be found in database
	ErrProductsByNameCanNotBeFound = errors.New("products by name can not be found")

	// ErrLanguageOfProductNameCanNotBeAdded means that language of productName edge can not be added
	ErrLanguageOfProductNameCanNotBeAdded = errors.New("language of productName can not be added")

	// ErrProductCanNotBeCreated means that the product can't be added to database
	ErrProductCanNotBeCreated = errors.New("product can't be created")

	// ErrProductCanNotBeDeleted means that the product can't be removed from database
	ErrProductCanNotBeDeleted = errors.New("product can't be deleted")
)

// Product is a structure of products in database
type Product struct {
	ID               string      `json:"uid, omitempty"`
	Name             string      `json:"productName, omitempty"`
	IRI              string      `json:"productIri, omitempty"`
	PreviewImageLink string      `json:"previewImageLink, omitempty"`
	IsActive         bool        `json:"productIsActive, omitempty"`
	Categories       []Category  `json:"belongs_to_category, omitempty"`
	Companies        []Companies `json:"belongs_to_company, omitempty"`
}

// Products is resource os storage for CRUD operations
type Products struct {
	storage *Storage
}

// NewProductsResourceForStorage is a constructor of Products resource
func NewProductsResourceForStorage(storage *Storage) *Products {
	return &Products{storage: storage}
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

// ReadProductsByName is a method for get all nodes by categories name
func (products *Products) ReadProductsByName(productName, language string) ([]Product, error) {
	query := fmt.Sprintf(`{
				products(func: eq(productName@%v, "%v")) @filter(eq(productIsActive, true)) {
					uid
					productName: productName@%v
					productIri
					previewImageLink
					productIsActive
					belongs_to_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName: categoryName@%v
						categoryIsActive
						belongs_to_company @filter(eq(companyIsActive, true)) {
							uid
							companyName: companyName@%v
							has_category @filter(eq(categoryIsActive, true)) {
								uid
								categoryName: categoryName@%v
								categoryIsActive
								belong_to_company @filter(eq(companyIsActive, true)){
									uid
									companyName: companyName@%v
									companyIsActive
								}
							}
						}
					}
					belongs_to_company @filter(eq(companyIsActive, true)) {
						uid
						companyName: companyName@%v
						has_category @filter(eq(categoryIsActive, true)) {
							uid
							categoryName: categoryName@%v
							categoryIsActive
							belong_to_company @filter(eq(companyIsActive, true)){
								uid
								companyName: companyName@%v
								companyIsActive
							}
						}
					}
				}
			}`, language, productName, language, language, language, language, language, language, language, language)

	transaction := products.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), query)
	if err != nil {
		log.Println(err)
		return nil, ErrProductsByNameCanNotBeFound
	}

	type productsInStorage struct {
		AllProductsFoundedByName []Product `json:"products"`
	}

	var foundedProducts productsInStorage
	err = json.Unmarshal(response.GetJson(), &foundedProducts)
	if err != nil {
		log.Println(err)
		return nil, ErrProductsByNameCanNotBeFound
	}

	if len(foundedProducts.AllProductsFoundedByName) == 0 {
		return nil, ErrProductsByNameNotFound
	}

	return foundedProducts.AllProductsFoundedByName, nil
}

// AddLanguageOfProductName is a method for add predicate "categoryName" for companyName value with new language
func (products *Products) AddLanguageOfProductName(productID, name, language string) error {
	forProductNamePredicate := fmt.Sprintf(`<%s> <productName> %s .`, productID, "\""+name+"\""+"@"+language)

	mutation := dataBaseAPI.Mutation{
		SetNquads: []byte(forProductNamePredicate),
		CommitNow: true}

	transaction := products.storage.Client.NewTxn()
	_, err := transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrLanguageOfProductNameCanNotBeAdded
	}

	return nil
}

// CreateProduct make product and save it to storage
func (products *Products) CreateProduct(product Product, language string) (Product, error) {
	existsProducts, err := products.ReadProductsByName(product.Name, language)
	if err != nil && err != ErrProductsByNameNotFound {
		log.Println(err)
		return product, ErrProductCanNotBeCreated
	}
	if existsProducts != nil {
		return existsProducts[0], ErrProductAlreadyExist
	}

	transaction := products.storage.Client.NewTxn()

	product.IsActive = true
	encodedProduct, err := json.Marshal(product)
	if err != nil {
		log.Println(err)
		return product, ErrProductCanNotBeCreated
	}

	mutation := &dataBaseAPI.Mutation{
		SetJson:   encodedProduct,
		CommitNow: true}

	assigned, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		log.Println(err)
		return product, ErrProductCanNotBeCreated
	}

	product.ID = assigned.Uids["blank-0"]
	if product.ID == "" {
		return product, ErrProductCanNotBeCreated
	}

	err = products.AddLanguageOfProductName(product.ID, product.Name, language)
	if err != nil {
		return product, err
	}

	return product, nil
}

// DeleteProduct method for remove product from database
func (products *Products) DeleteProduct(product Product) (string, error) {
	if product.ID == "" {
		return "", ErrProductCanNotBeWithoutID
	}

	deleteProductData, _ := json.Marshal(map[string]string{"uid": product.ID})

	mutation := dataBaseAPI.Mutation{
		DeleteJson: deleteProductData,
		CommitNow:  true}

	transaction := products.storage.Client.NewTxn()

	var err error
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		log.Println(err)
		return product.ID, ErrProductCanNotBeDeleted
	}

	return product.ID, nil
}
