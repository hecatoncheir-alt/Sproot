package engine

// Product if a struct of item of company
type Product struct {
	Name    string
	Preview string
	Link    string
	Company string
}

// GetProductOfCompany is method for get product from database
func (engine *Engine) GetProductOfCompany(product *Product) {}

// SaveProductOfCompany is method for add product to database
func (engine *Engine) SaveProductOfCompany(product *Product) {
	// Проверить наличие продукта по имени
}
