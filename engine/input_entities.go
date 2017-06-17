package engine

import "time"

// MessageEvent struct for incoming event
type MessageEvent struct {
	Message string
	Data    interface{}
}

// Company type for parse
type Company struct {
	ID         string
	IRI        string
	Name       string
	Categories []string
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
