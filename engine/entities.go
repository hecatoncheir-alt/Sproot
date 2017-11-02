package engine

import "time"

// MessageEvent struct for incoming event
type MessageEvent struct {
	Message string
	Data    interface{}
}

// Category is a structure of Category in database
type Category struct {
	ID   uint64 `json:"_uid_,omitempty"`
	Name string `json:"name,omitempty"`
}

// Company type for parse
type Company struct {
	ID         uint64 `json:"_uid_,omitempty"`
	IRI        string `json:"iri,omitempty"`
	Name       string `json:"name,omitempty"`
	Categories []Category
}

// Price structure
type Price struct {
	Value    string
	City     string
	DateTime time.Time
}

// Item is a structure of one product from one page
type Item struct {
	Name             string
	Link             string
	PreviewImageLink string
	Price            Price
	Company          Company
}
