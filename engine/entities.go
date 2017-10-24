package engine

import "time"

// MessageEvent struct for incoming event
type MessageEvent struct {
	Message string
	Data    interface{}
}

// Category is a structure of Category in database
type Category struct {
	ID   string
	Name string
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

// GraphDataResponseField is a structure for field in response from DGraph database operation
type GraphDataResponseField struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Uids    map[string]string `json:"uids"`
}

// GraphResponse is a response structure from DGraph database operation
type GraphResponse struct {
	Data       GraphDataResponseField `json:"data"`
	Extensions map[string]string      `json:"extensions"`
}
