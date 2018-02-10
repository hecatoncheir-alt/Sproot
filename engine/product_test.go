package engine

import (
	"encoding/json"
	"testing"
)

func TestIntegrationNewPriceWithProductCanBeCreated(test *testing.T) {

	//TODO create category, company, and city with one language

	productWithPriceJSON :=
		`{
			"Name":"Смартфон Samsung Galaxy S8 64Gb Черный бриллиант",
			"IRI":"http://www.mvideo.ru//products/smartfon-samsung-galaxy-s8-64gb-chernyi-brilliant-30027818",
			"PreviewImageLink":"img.mvideo.ru/Pdb/30027818m.jpg",
			"Language":"en",
			"Price":{
				"Value":"46990",
				"City":{
					"ID":"0x2788",
					"Name":"Москва"
				},
				"DateTime":"2018-02-10T08:34:35.6055814Z"
			},
			"Company":{
				"ID":"0x2786",
				"Name":"М.Видео",
				"IRI":"http://www.mvideo.ru/"
			},
			"Category":{
				"ID":"",
				"Name":"Тестовая категория"
			},
			"City":{
				"ID":"0x2788",
				"Name":"Москва"
			}
		}`

	product := ProductOfCompany{}
	json.Unmarshal([]byte(productWithPriceJSON), &product)

	err := product.UpdateInStorage()
	if err != nil {
		test.Fail()
	}

	// TODO get product with prices from storage

}
