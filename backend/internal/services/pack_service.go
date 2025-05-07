package services

import (
	"fmt"
	"sort"
	"sync"

	"pack-sizer/internal/models"
)

// PackSizesManager manages the pack sizes configuration
type PackSizesManager struct {
	sizes []int
	mu    sync.RWMutex
}

// NewPackSizesManager creates a new pack sizes manager with default sizes
func NewPackSizesManager() *PackSizesManager {
	return &PackSizesManager{
		sizes: []int{250, 500, 1000, 2000, 5000},
	}
}

// GetSizes returns the current pack sizes
func (p *PackSizesManager) GetSizes() []int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return append([]int{}, p.sizes...)
}

// AddSize adds a new pack size
func (p *PackSizesManager) AddSize(size int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// Check if size already exists
	for _, s := range p.sizes {
		if s == size {
			return
		}
	}
	
	p.sizes = append(p.sizes, size)
	sort.Ints(p.sizes)
}

// RemoveSize removes a pack size
func (p *PackSizesManager) RemoveSize(size int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	newSizes := make([]int, 0, len(p.sizes))
	for _, s := range p.sizes {
		if s != size {
			newSizes = append(newSizes, s)
		}
	}
	p.sizes = newSizes
}

// UpdateSize updates an existing pack size with a new one
func (p *PackSizesManager) UpdateSize(oldSize, newSize int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if old size exists
	oldSizeExists := false
	for _, s := range p.sizes {
		if s == oldSize {
			oldSizeExists = true
			break
		}
	}

	if !oldSizeExists {
		return fmt.Errorf("old size %d does not exist", oldSize)
	}

	// Check if new size already exists
	for _, s := range p.sizes {
		if s == newSize {
			return fmt.Errorf("new size %d already exists", newSize)
		}
	}

	// Remove old size
	newSizes := make([]int, 0, len(p.sizes))
	for _, s := range p.sizes {
		if s != oldSize {
			newSizes = append(newSizes, s)
		}
	}

	// Add new size
	newSizes = append(newSizes, newSize)
	sort.Ints(newSizes)
	p.sizes = newSizes

	return nil
}

// CalculatePacks calculates the optimal pack distribution for a given order quantity
func CalculatePacks(orderQuantity int, packSizes []int) models.PackCalculation {
	// Sort pack sizes in ascending order
	sort.Ints(packSizes)

	result := models.PackCalculation{
		OrderQuantity: orderQuantity,
		Packs:         make(map[int]int),
		TotalItems:    0,
	}

	// If order quantity is 0, return empty result
	if orderQuantity == 0 {
		return result
	}

	// Handle case where orderQuantity is less than smallest pack size
	smallestPack := packSizes[0]
	if orderQuantity < smallestPack {
		result.Packs[smallestPack] = 1
		result.TotalItems = smallestPack
		return result
	}

	// Define the maximum possible reasonable order we need to consider
	maxPossibleOrder := orderQuantity + packSizes[len(packSizes)-1]

	// Define a more efficient solution struct
	type solution struct {
		totalItems int
		packs      map[int]int
	}

	// Initialize our dp array (indexed by order quantity)
	dp := make([]solution, maxPossibleOrder+1)
	for i := range dp {
		dp[i].totalItems = maxPossibleOrder + 1 // "Infinity"
		dp[i].packs = make(map[int]int)
	}
	
	// Base case: 0 items needs 0 items/packs
	dp[0].totalItems = 0
	
	// Build up solutions for each quantity
	for quantity := 1; quantity <= maxPossibleOrder; quantity++ {
		// Try each pack size
		for _, packSize := range packSizes {
			if packSize > quantity {
				continue
			}
			
			prevSolution := dp[quantity-packSize]
			newTotalItems := prevSolution.totalItems + packSize
			
			// Calculate total packs for comparing solutions
			prevTotalPacks := 0
			for _, count := range prevSolution.packs {
				prevTotalPacks += count
			}
			
			currentTotalPacks := 0
			for _, count := range dp[quantity].packs {
				currentTotalPacks += count
			}
			
			// Better solution if:
			// 1. Total items is less, OR
			// 2. Total items equal but fewer total packs
			isBetter := newTotalItems < dp[quantity].totalItems || 
				(newTotalItems == dp[quantity].totalItems && prevTotalPacks+1 < currentTotalPacks)
			
			if isBetter {
				// Copy previous solution's packs
				newPacks := make(map[int]int)
				for size, count := range prevSolution.packs {
					newPacks[size] = count
				}
				
				// Add current pack
				newPacks[packSize]++
				
				dp[quantity] = solution{
					totalItems: newTotalItems,
					packs:      newPacks,
				}
			}
		}
	}
	
	// Find the smallest valid solution (one that meets or exceeds orderQuantity)
	var bestSolution solution
	bestTotalItems := maxPossibleOrder + 1
	
	for quantity := orderQuantity; quantity <= maxPossibleOrder; quantity++ {
		if dp[quantity].totalItems < bestTotalItems {
			bestTotalItems = dp[quantity].totalItems
			bestSolution = dp[quantity]
		}
	}
	
	result.Packs = bestSolution.packs
	result.TotalItems = bestSolution.totalItems
	
	return result
}