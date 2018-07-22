package engine

import (
	"time"

	"github.com/hecatoncheir/Sproot/engine/storage"
)

// ProductOfCompany object for create new price or new product in storage
type ProductOfCompany struct {
	Name             string
	IRI              string
	PreviewImageLink string
	Language         string
	Price            PriceOfProduct
	Company          CompanyData
	Category         CategoryData
}

// PriceOfProduct part of structure of ProductOfCompany
type PriceOfProduct struct {
	Value    float64
	DateTime time.Time
	City     CityData
}

// UpdateInStorage method for create product if it needed or add price to product
func (product *ProductOfCompany) UpdateInStorage(store *storage.Storage) (storage.Product, error) {
	products, err := store.Products.ReadProductsByName(product.Name, product.Language)

	productFromStorage := storage.Product{
		Name:             product.Name,
		IRI:              product.IRI,
		PreviewImageLink: product.PreviewImageLink}

	if len(products) != 0 {
		for _, productByName := range products {
			for _, category := range productByName.Categories {
				if category.ID == product.Category.ID {
					productFromStorage = productByName
					break
				}
			}
			//break
		}
	}

	if err != nil && err != storage.ErrProductsByNameNotFound {
		return productFromStorage, err
	}

	if err == storage.ErrProductsByNameNotFound || products == nil {
		productForStorage := storage.Product{
			Name:             product.Name,
			IRI:              product.IRI,
			PreviewImageLink: product.PreviewImageLink}

		productInStorage, err := store.Products.CreateProduct(productForStorage, product.Language)
		if err != nil {
			return productInStorage, err
		}

		productFromStorage.ID = productInStorage.ID
		productFromStorage.IsActive = productInStorage.IsActive

		err = store.Products.AddCategoryToProduct(productInStorage.ID, product.Category.ID)
		if err != nil {
			return productInStorage, err
		}

		err = store.Products.AddCompanyToProduct(productInStorage.ID, product.Company.ID)
		if err != nil {
			return productInStorage, err
		}
	}

	// priceValue, err := strconv.ParseFloat(product.Price.Value, 64)
	// if err != nil {
	// 	return productFromStorage, err
	// }

	priceForStorage := storage.Price{Value: product.Price.Value, DateTime: product.Price.DateTime}

	priceFromStorage, err := store.Prices.CreatePrice(priceForStorage)
	if err != nil {
		return productFromStorage, err
	}

	err = store.Prices.AddCompanyToPrice(priceFromStorage.ID, product.Company.ID)
	if err != nil {
		return productFromStorage, err
	}

	err = store.Prices.AddProductToPrice(priceFromStorage.ID, productFromStorage.ID)
	if err != nil {
		return productFromStorage, err
	}

	err = store.Prices.AddCityToPrice(priceFromStorage.ID, product.Price.City.ID)
	if err != nil {
		return productFromStorage, err
	}

	productFromStorage, err = store.Products.ReadProductByID(productFromStorage.ID, product.Language)
	if err != nil {
		return productFromStorage, err
	}

	return productFromStorage, nil
}
