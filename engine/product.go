package engine

type ProductOfCompany struct {
	Name             string
	IRI              string
	PreviewImageLink string
	Language         string
	Price            PriceOfProduct
	Company          CompanyData
	Category         CityData
}

// TODO
func (product *ProductOfCompany) UpdateInStorage() error {

	return nil
}
