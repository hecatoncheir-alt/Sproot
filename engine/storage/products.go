package storage

import (
	"context"
	"log"

	"encoding/json"
	"errors"
	"fmt"

	"bytes"
	"text/template"

	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
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

	// ErrProductByIDCanNotBeFound means that the product can't be found in database
	ErrProductByIDCanNotBeFound = errors.New("product by id can not be found")

	// ErrProductDoesNotExist means than the product does not exist in database
	ErrProductDoesNotExist = errors.New("product by id not found")

	// ErrCompanyCanNotBeAddedToProduct means that the company can't be added to product
	ErrCompanyCanNotBeAddedToProduct = errors.New("company can not be added to product")

	// ErrCategoryCanNotBeAddedToProduct means that the category can't be added to product
	ErrCategoryCanNotBeAddedToProduct = errors.New("category can not be added to product")

	// ErrPriceCanNotBeAddedToProduct means that the price can't be added to product
	ErrPriceCanNotBeAddedToProduct = errors.New("price can not be added to product")
)

// Product is a structure of products in database
type Product struct {
	ID               string     `json:"uid,omitempty"`
	Name             string     `json:"productName,omitempty"`
	IRI              string     `json:"productIri,omitempty"`
	PreviewImageLink string     `json:"previewImageLink,omitempty"`
	IsActive         bool       `json:"productIsActive"`
	Categories       []Category `json:"belongs_to_category,omitempty"`
	Companies        []Company  `json:"belongs_to_company,omitempty"`
	Prices           []Price    `json:"has_price,omitempty"`
}

// Products is resource of storage for CRUD operations
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
		productName: string @lang @index(term, trigram) .
		productIri: string @index(term) .
		productImageLink: string @index(term) .
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

func (products *Products) ReadTotalCountOfProductsByName(productName, language string) (int, error) {
	type Variables struct {
		ProductName, Language     string
		CurrentPage, ItemsPerPage int
	}

	variables := Variables{
		ProductName: productName,
		Language:    language}

	totalQueryTemplate, err := template.New("totalQuery").Parse(`{
				total(func: regexp(productName@{{.Language}}, /{{.ProductName}}/))
				@filter(eq(productIsActive, true) AND has(productName)) {
					count: count(uid)
				}
			}`)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	totalQueryBuf := bytes.Buffer{}

	err = totalQueryTemplate.Execute(&totalQueryBuf, variables)
	if err != nil {
		log.Println(err)
		return 0, ErrProductsByNameCanNotBeFound
	}

	transaction := products.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), totalQueryBuf.String())
	if err != nil {
		log.Println(err)
		return 0, ErrProductsByNameCanNotBeFound
	}

	type productsCountInStorage struct {
		Total []map[string]int `json:"total"`
	}

	var foundedProducts productsCountInStorage
	err = json.Unmarshal(response.GetJson(), &foundedProducts)
	if err != nil {
		log.Println(err)
		return 0, ErrProductsByNameCanNotBeFound
	}

	return foundedProducts.Total[0]["count"], nil
}

type ProductsByNameForPage struct {
	Products                []Product
	CurrentPage             int
	TotalProductsForOnePage int
	TotalProductsFound      int
	SearchedName            string
	Language                string
}

func (products *Products) ReadProductsByNameWithPagination(productName, language string, currentPage, itemsPerPage int) (*ProductsByNameForPage, error) {

	type Variables struct {
		ProductName, Language             string
		CurrentPage, ItemsPerPage, Offset int
	}

	productsByPageTemplate, err := template.New("productsByPage").Parse(`{
				all as counters(func: regexp(productName@{{.Language}}, /{{.ProductName}}/i))
				@filter(eq(productIsActive, true) AND has(productName)){
					total: count(uid)
				}

				products(func: uid(all), first: {{.ItemsPerPage}}, offset: {{.Offset}})
				@filter(eq(productIsActive, true) AND has(productName)) {
					uid
					productName: productName@{{.Language}}
					productIri
					previewImageLink
					productIsActive
					belongs_to_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName: categoryName@{{.Language}}
						categoryIsActive
						belongs_to_company @filter(eq(companyIsActive, true)) {
							uid
							companyName: companyName@{{.Language}}
							has_category @filter(eq(categoryIsActive, true)) {
								uid
								categoryName: categoryName@{{.Language}}
								categoryIsActive
								belong_to_company @filter(eq(companyIsActive, true)){
									uid
									companyName: companyName@{{.Language}}
									companyIsActive
								}
							}
						}
					}
					belongs_to_company @filter(eq(companyIsActive, true)) {
						uid
						companyName: companyName@{{.Language}}
						has_category @filter(eq(categoryIsActive, true)) {
							uid
							categoryName: categoryName@{{.Language}}
							categoryIsActive
							belong_to_company @filter(eq(companyIsActive, true)){
								uid
								companyName: companyName@{{.Language}}
								companyIsActive
							}
						}
					}
					has_price @filter(eq(priceIsActive, true)) {
						uid
						priceValue
						priceDateTime
						priceCity
						priceIsActive
						belongs_to_product @filter(eq(productIsActive, true)) {
							uid
							productName: productName@{{.Language}}
							productIri
							previewImageLink
							productIsActive
							has_price @filter(eq(priceIsActive, true)) {
								uid
								priceValue
								priceDateTime
								priceCity
								priceIsActive
							}
						}
						belongs_to_city @filter(eq(cityIsActive, true)) {
							uid
							cityName: cityName@{{.Language}}
							cityIsActive
						}
						belongs_to_company @filter(eq(companyIsActive, true)) {
							uid
							companyName: companyName@{{.Language}}
							companyIri
							companyIsActive
						}
					}
				}
			}`)

	if err != nil {
		log.Println(err)
		return nil, ErrProductsByNameCanNotBeFound
	}

	variables := Variables{
		ProductName:  productName,
		ItemsPerPage: itemsPerPage,
		CurrentPage:  currentPage,
		Offset:       currentPage*itemsPerPage - itemsPerPage,
		Language:     language}

	productsByPageBuf := bytes.Buffer{}

	err = productsByPageTemplate.Execute(&productsByPageBuf, variables)
	if err != nil {
		log.Println(err)
		return nil, ErrProductsByNameCanNotBeFound
	}

	transaction := products.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), productsByPageBuf.String())
	if err != nil {
		log.Println(err)
		return nil, ErrProductsByNameCanNotBeFound
	}

	type productsInStorage struct {
		Total                    []map[string]int `json:"counters"`
		AllProductsFoundedByName []Product        `json:"products"`
	}

	var foundedProducts productsInStorage
	err = json.Unmarshal(response.GetJson(), &foundedProducts)
	if err != nil {
		log.Println(err)
		return nil, ErrProductsByNameCanNotBeFound
	}

	foundedProductsByNameForPage := ProductsByNameForPage{
		Products:                foundedProducts.AllProductsFoundedByName,
		CurrentPage:             currentPage,
		TotalProductsForOnePage: itemsPerPage,
		SearchedName:            productName,
		TotalProductsFound:      foundedProducts.Total[0]["total"],
		Language:                language}

	if len(foundedProducts.AllProductsFoundedByName) == 0 {
		return &foundedProductsByNameForPage, ErrProductsByNameNotFound
	}

	return &foundedProductsByNameForPage, nil
}

// ReadProductsByName is a method for get all nodes by product name
func (products *Products) ReadProductsByName(productName, language string) ([]Product, error) {

	variables := struct {
		ProductName, Language string
	}{
		ProductName: productName,
		Language:    language}

	productsByNameTemplate, err := template.New("productsByName").Parse(`{
				products(func: regexp(productName@{{.Language}}, /{{.ProductName}}/)) 
				@filter(eq(productIsActive, true) AND has(productName)) {
					uid
					productName: productName@{{.Language}}
					productIri
					previewImageLink
					productIsActive
					belongs_to_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName: categoryName@{{.Language}}
						categoryIsActive
						belongs_to_company @filter(eq(companyIsActive, true)) {
							uid
							companyName: companyName@{{.Language}}
							has_category @filter(eq(categoryIsActive, true)) {
								uid
								categoryName: categoryName@{{.Language}}
								categoryIsActive
								belong_to_company @filter(eq(companyIsActive, true)){
									uid
									companyName: companyName@{{.Language}}
									companyIsActive
								}
							}
						}
					}
					belongs_to_company @filter(eq(companyIsActive, true)) {
						uid
						companyName: companyName@{{.Language}}
						companyIri
						companyIsActive
						has_category @filter(eq(categoryIsActive, true)) {
							uid
							categoryName: categoryName@{{.Language}}
							categoryIsActive
							belong_to_company @filter(eq(companyIsActive, true)){
								uid
								companyName: companyName@{{.Language}}
								companyIsActive
							}
						}
					}
					has_price @filter(eq(priceIsActive, true)) {
						uid
						priceValue
						priceDateTime
						priceIsActive
						belongs_to_company @filter(eq(companyIsActive, true)) {
							uid
							companyName: companyName@{{.Language}}
							companyIri
							companyIsActive
						}
						belongs_to_product @filter(eq(productIsActive, true)) {
							uid
							productName: productName@{{.Language}}
							productIri
							previewImageLink
							productIsActive
							has_price @filter(eq(priceIsActive, true)) {
								uid
								priceValue
								priceDateTime
								priceIsActive
								belongs_to_company @filter(eq(companyIsActive, true)) {
									uid
									companyName: companyName@{{.Language}}
									companyIri
									companyIsActive
								}
							}
						}
						belongs_to_city @filter(eq(cityIsActive, true)) {
							uid
							cityName: cityName@{{.Language}}
							cityIsActive
						}
					}
				}
			}`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	productsByNameQueryBuf := bytes.Buffer{}

	err = productsByNameTemplate.Execute(&productsByNameQueryBuf, variables)
	if err != nil {
		log.Println(err)
		return nil, ErrProductsByNameCanNotBeFound
	}

	transaction := products.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), productsByNameQueryBuf.String())
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

// ReadProductByID is a method for get all nodes of products by ID
func (products *Products) ReadProductByID(productID, language string) (Product, error) {

	variables := struct {
		ProductID, Language string
	}{
		ProductID: productID,
		Language:  language}

	queryTemplate, err := template.New("ReadProductByID").Parse(`{
				products(func: uid("{{.ProductID}}")) @filter(has(productName)) {
					uid
					productName: productName@{{.Language}}
					productIri
					previewImageLink
					productIsActive
					belongs_to_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName: categoryName@{{.Language}}
						categoryIsActive
						belongs_to_company @filter(eq(companyIsActive, true)) {
							uid
							companyName: companyName@{{.Language}}
							has_category @filter(eq(categoryIsActive, true)) {
								uid
								categoryName: categoryName@{{.Language}}
								categoryIsActive
								belong_to_company @filter(eq(companyIsActive, true)){
									uid
									companyName: companyName@{{.Language}}
									companyIsActive
								}
							}
						}
					}
					belongs_to_company @filter(eq(companyIsActive, true)) {
						uid
						companyName: companyName@{{.Language}}
						has_category @filter(eq(categoryIsActive, true)) {
							uid
							categoryName: categoryName@{{.Language}}
							categoryIsActive
							belong_to_company @filter(eq(companyIsActive, true)){
								uid
								companyName: companyName@{{.Language}}
								companyIsActive
							}
						}
					}
					has_price @filter(eq(priceIsActive, true)) {
						uid
						priceValue
						priceDateTime
						priceCity
						priceIsActive
						belongs_to_product @filter(eq(productIsActive, true)) {
							uid
							productName: productName@{{.Language}}
							productIri
							previewImageLink
							productIsActive
							has_price @filter(eq(priceIsActive, true)) {
								uid
								priceValue
								priceDateTime
								priceCity
								priceIsActive
							}
						}
						belongs_to_city @filter(eq(cityIsActive, true)) {
							uid
							cityName: cityName@{{.Language}}
							cityIsActive
						}
					}
				}
			}`)

	product := Product{ID: productID}
	if err != nil {
		log.Println(err)
		return product, ErrProductByIDCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	execErr := queryTemplate.Execute(&queryBuf, variables)
	if  execErr!= nil {
		log.Println(execErr)
		return product,execErr
	}

	transaction := products.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		log.Println(err)
		return product, ErrProductByIDCanNotBeFound
	}

	type productsInStore struct {
		Products []Product `json:"products"`
	}

	var foundedProducts productsInStore

	err = json.Unmarshal(response.GetJson(), &foundedProducts)
	if err != nil {
		log.Println(err)
		return product, ErrProductByIDCanNotBeFound
	}

	if len(foundedProducts.Products) == 0 {
		return product, ErrProductDoesNotExist
	}

	return foundedProducts.Products[0], nil
}

// AddCategoryToProduct method for set quad of predicate about product and category
func (products *Products) AddCategoryToProduct(productID, categoryID string) error {
	var err error
	var mutation dataBaseAPI.Mutation

	forCategoryPredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, categoryID, "has_product", productID)

	mutation = dataBaseAPI.Mutation{
		SetNquads: []byte(forCategoryPredicate),
		CommitNow: true}

	transaction := products.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrProductCanNotBeAddedToCategory
	}

	forProductPredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, productID, "belongs_to_category", categoryID)
	mutation = dataBaseAPI.Mutation{
		SetNquads: []byte(forProductPredicate),
		CommitNow: true}

	transaction = products.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrCategoryCanNotBeAddedToProduct
	}

	return nil
}

// AddCompanyToProduct method for set quad of predicate about product and company
func (products *Products) AddCompanyToProduct(productID, companyID string) error {
	var err error
	var mutation dataBaseAPI.Mutation

	forProductPredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, productID, "belongs_to_company", companyID)
	mutation = dataBaseAPI.Mutation{
		SetNquads: []byte(forProductPredicate),
		CommitNow: true}

	transaction := products.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrCompanyCanNotBeAddedToProduct
	}

	return nil
}

// AddPriceToProduct method for set quad of predicate about product and price
func (products *Products) AddPriceToProduct(productID, priceID string) error {
	var err error
	var mutation dataBaseAPI.Mutation

	forPricePredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, priceID, "belongs_to_product", productID)
	mutation = dataBaseAPI.Mutation{
		SetNquads: []byte(forPricePredicate),
		CommitNow: true}

	transaction := products.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrPriceCanNotBeAddedToProduct
	}

	forProductPredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, productID, "has_price", priceID)
	mutation = dataBaseAPI.Mutation{
		SetNquads: []byte(forProductPredicate),
		CommitNow: true}

	transaction = products.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrPriceCanNotBeAddedToProduct
	}

	return nil
}
