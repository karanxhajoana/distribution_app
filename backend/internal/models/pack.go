package models

// PackSize represents a pack size configuration
type PackSize struct {
	Size int `json:"size"`
}

// PackCalculation represents the result of a pack calculation
type PackCalculation struct {
	OrderQuantity int         `json:"orderQuantity"`
	Packs         map[int]int `json:"packs"`
	TotalItems    int         `json:"totalItems"`
} 