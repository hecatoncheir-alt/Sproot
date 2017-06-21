package engine

// PriceOfProduct is agregate of price data of product
type PriceOfProduct struct {
	Name string
}

// GetPricesOfProductsByName is method for get Price of product from database
func (engine *Engine) GetPricesOfProductsByName(productName string) ([]*PriceOfProduct, error) {
	return []*PriceOfProduct{&PriceOfProduct{Name: productName}}, nil
}

// SavePriceForProductOfCompany method for save subject, predicate and object in graph database
func (engine *Engine) SavePriceForProductOfCompany(item *Item) (*PriceOfProduct, error) {
	price := PriceOfProduct{Name: item.Name}
	return &price, nil
}
