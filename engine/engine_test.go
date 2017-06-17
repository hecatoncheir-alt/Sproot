package engine

import (
	"testing"
)

// {Смартфон Samsung Galaxy S8+ 64Gb Мистический аметист http://www.mvideo.ru//products/smartfon-samsung-galaxy-s8-64gb-misticheskii-ametist-30027823 img.mvideo.ru/Pdb/30027823m.jpg {59990  2017-06-17 16:07:13.888498569 +0000 UTC} { http://www.mvideo.ru/ M.Video [Телефоны]}}
func TestSprootCanSaveGetAndDeleteData(test *testing.T) {
	incomingItem := Item{}

	item, err := SavePriceForProductOfCompany(incomingItem)
	if err != nil {
		test.Error(err)
	}

	if item.Name != "test item name" {
		test.Fail()
	}
}
