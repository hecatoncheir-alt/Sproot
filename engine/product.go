package engine

import (
	"strconv"
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
	Value    string
	DateTime time.Time
	City     CityData
}

// UpdateInStorage method for create product if it needed or add price to product
func (product *ProductOfCompany) UpdateInStorage(store *storage.Storage) (storage.Product, error) {
	products, err := store.Products.ReadProductsByName(product.Name, product.Language)

	var productFromStorage = storage.Product{}

	if products != nil {
		for _, productByName := range products {
			for _, category := range productByName.Categories {
				if category.ID == product.Category.ID {
					productFromStorage = productByName
					break
				}
			}
			break
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

		productFromStorage, err = store.Products.CreateProduct(productForStorage, product.Language)
		if err != nil {
			return productFromStorage, err
		}

		err = store.Products.AddCategoryToProduct(productFromStorage.ID, product.Category.ID)
		if err != nil {
			return productFromStorage, err
		}

		err = store.Products.AddCompanyToProduct(productFromStorage.ID, product.Company.ID)
		if err != nil {
			return productFromStorage, err
		}
	}

	priceValue, err := strconv.ParseFloat(product.Price.Value, 64)
	if err != nil {
		return productFromStorage, err
	}

	priceForStorage := storage.Price{Value: priceValue, DateTime: product.Price.DateTime}

	priceFromStorage, err := store.Prices.CreatePrice(priceForStorage)
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
