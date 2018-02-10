package engine

import "time"

type ProductOfCompany struct {
	Name             string
	IRI              string
	PreviewImageLink string
	Language         string
	Price            PriceOfProduct
	Company          CompanyData
	Category         CityData
}

type PriceOfProduct struct {
	Value    string
	DateTime time.Time
}

// TODO
func (product *ProductOfCompany) UpdateInStorage() error {

	return nil
}
