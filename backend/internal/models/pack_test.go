package models

import (
	"testing"
)

func TestPackSize(t *testing.T) {
	packSize := PackSize{Size: 250}
	if packSize.Size != 250 {
		t.Errorf("Expected pack size to be 250, got %d", packSize.Size)
	}
}

func TestPackCalculation(t *testing.T) {
	calculation := PackCalculation{
		OrderQuantity: 100,
		Packs: map[int]int{
			250: 1,
		},
		TotalItems: 250,
	}

	if calculation.OrderQuantity != 100 {
		t.Errorf("Expected order quantity to be 100, got %d", calculation.OrderQuantity)
	}

	if len(calculation.Packs) != 1 {
		t.Errorf("Expected 1 pack size, got %d", len(calculation.Packs))
	}

	if calculation.Packs[250] != 1 {
		t.Errorf("Expected 1 pack of size 250, got %d", calculation.Packs[250])
	}

	if calculation.TotalItems != 250 {
		t.Errorf("Expected total items to be 250, got %d", calculation.TotalItems)
	}
} 