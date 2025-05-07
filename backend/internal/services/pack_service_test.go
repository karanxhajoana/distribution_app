package services

import (
	"testing"
)

func TestNewPackSizesManager(t *testing.T) {
	manager := NewPackSizesManager()
	expectedSizes := []int{250, 500, 1000, 2000, 5000}
	
	sizes := manager.GetSizes()
	if len(sizes) != len(expectedSizes) {
		t.Errorf("Expected %d pack sizes, got %d", len(expectedSizes), len(sizes))
	}
	
	for i, size := range sizes {
		if size != expectedSizes[i] {
			t.Errorf("Expected pack size %d to be %d, got %d", i, expectedSizes[i], size)
		}
	}
}

func TestPackSizesManager_AddSize(t *testing.T) {
	manager := NewPackSizesManager()
	
	// Test adding a new size
	manager.AddSize(300)
	sizes := manager.GetSizes()
	found := false
	for _, size := range sizes {
		if size == 300 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find size 300 in pack sizes")
	}
	
	// Test adding a duplicate size
	originalLength := len(sizes)
	manager.AddSize(300)
	sizes = manager.GetSizes()
	if len(sizes) != originalLength {
		t.Errorf("Expected pack sizes length to remain %d, got %d", originalLength, len(sizes))
	}
}

func TestPackSizesManager_RemoveSize(t *testing.T) {
	manager := NewPackSizesManager()
	
	// Test removing an existing size
	manager.RemoveSize(250)
	sizes := manager.GetSizes()
	for _, size := range sizes {
		if size == 250 {
			t.Error("Expected size 250 to be removed from pack sizes")
		}
	}
	
	// Test removing a non-existent size
	originalLength := len(sizes)
	manager.RemoveSize(999)
	sizes = manager.GetSizes()
	if len(sizes) != originalLength {
		t.Errorf("Expected pack sizes length to remain %d, got %d", originalLength, len(sizes))
	}
}

func TestPackSizesManager_UpdateSize(t *testing.T) {
	tests := []struct {
		name        string
		oldSize     int
		newSize     int
		expectError bool
		wantSizes   []int
	}{
		{
			name:        "Valid update",
			oldSize:     250,
			newSize:     300,
			expectError: false,
			wantSizes:   []int{300, 500, 1000, 2000, 5000},
		},
		{
			name:        "Non-existent old size",
			oldSize:     999,
			newSize:     300,
			expectError: true,
			wantSizes:   []int{250, 500, 1000, 2000, 5000},
		},
		{
			name:        "New size already exists",
			oldSize:     250,
			newSize:     500,
			expectError: true,
			wantSizes:   []int{250, 500, 1000, 2000, 5000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh manager for each test case
			manager := NewPackSizesManager()
			
			err := manager.UpdateSize(tt.oldSize, tt.newSize)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			sizes := manager.GetSizes()
			if len(sizes) != len(tt.wantSizes) {
				t.Errorf("Expected %d sizes, got %d", len(tt.wantSizes), len(sizes))
			}

			for i, size := range sizes {
				if size != tt.wantSizes[i] {
					t.Errorf("Expected size at index %d to be %d, got %d", i, tt.wantSizes[i], size)
				}
			}
		})
	}
}

func TestCalculatePacks(t *testing.T) {
	tests := []struct {
		name          string
		orderQuantity int
		packSizes     []int
		expectedPacks map[int]int
		expectedTotal int
	}{
		{
			name:          "Zero order quantity",
			orderQuantity: 0,
			packSizes:     []int{250, 500, 1000},
			expectedPacks: map[int]int{},
			expectedTotal: 0,
		},
		{
			name:          "Order smaller than smallest pack",
			orderQuantity: 100,
			packSizes:     []int{250, 500, 1000},
			expectedPacks: map[int]int{250: 1},
			expectedTotal: 250,
		},
		{
			name:          "Exact pack size match",
			orderQuantity: 500,
			packSizes:     []int{250, 500, 1000},
			expectedPacks: map[int]int{500: 1},
			expectedTotal: 500,
		},
		{
			name:          "Multiple packs needed",
			orderQuantity: 750,
			packSizes:     []int{250, 500, 1000},
			expectedPacks: map[int]int{500: 1, 250: 1},
			expectedTotal: 750,
		},
		{
			name:          "Large order requiring multiple packs",
			orderQuantity: 1200,
			packSizes:     []int{250, 500, 1000},
			expectedPacks: map[int]int{1000: 1, 250: 1},
			expectedTotal: 1250,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePacks(tt.orderQuantity, tt.packSizes)
			
			if result.OrderQuantity != tt.orderQuantity {
				t.Errorf("Expected order quantity %d, got %d", tt.orderQuantity, result.OrderQuantity)
			}
			
			if result.TotalItems != tt.expectedTotal {
				t.Errorf("Expected total items %d, got %d", tt.expectedTotal, result.TotalItems)
			}
			
			if len(result.Packs) != len(tt.expectedPacks) {
				t.Errorf("Expected %d pack sizes, got %d", len(tt.expectedPacks), len(result.Packs))
			}
			
			for size, count := range tt.expectedPacks {
				if result.Packs[size] != count {
					t.Errorf("Expected %d packs of size %d, got %d", count, size, result.Packs[size])
				}
			}
		})
	}
} 