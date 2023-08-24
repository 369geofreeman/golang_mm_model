package optimization

import (
	"fmt"
	"testing"

	"github.com/lanl/clp"
)

var testCases = []struct {
	currentPrice   float64
	inventory      *Inventory
	volatility     float64
	liquidity      float64
	orderBookDepth float64
	expectedBid    float64
	expectedAsk    float64
	expectedPrimal clp.SimplexStatus
}{
	{
		26080.15, &Inventory{1000, 500, 0.02}, 2.04828141212099, 0.59, 3,
		26080.050000000003, 26080.15, 0, // Expected values
	},
	{
		26080.15000, &Inventory{1000, 500, 0.02}, 2.04828, 0.59, 3,
		26080.04999, 26080.15000, 0,
	},
	{
		21045.67890, &Inventory{900, 400, 0.01}, 1.23456, 0.45, 2,
		21045.57890, 21045.67890, 0,
	},
	{
		29500.12345, &Inventory{2000, 100, 0.03}, 1.54321, 0.25, 1,
		29499.92345, 29500.12345, 0,
	},
	{
		25500.54321, &Inventory{500, 500, 0.02}, 0.76543, 0.50, 0,
		25499.94321, 25500.54321, 0,
	},
	{
		27000.98765, &Inventory{1500, 150, 0.01}, 0.87654, 0.60, 4,
		26999.88765, 27000.98765, 0,
	},
	{
		20500.13579, &Inventory{700, 700, 0.03}, 0.98765, 0.35, 2,
		20499.83579, 20500.13579, 0,
	},
	{
		23000.24680, &Inventory{1000, 400, 0.02}, 0.24680, 0.55, 3,
		22999.94680, 23000.24680, 0,
	},
	{
		24000.86420, &Inventory{1200, 300, 0.02}, 1.12345, 0.40, 1,
		23999.96420, 24000.86420, 0,
	},
	{
		27500.97531, &Inventory{1300, 200, 0.03}, 0.97531, 0.45, 2,
		27499.87531, 27500.97531, 0,
	},
	{
		29000.13579, &Inventory{900, 400, 0.01}, 0.86420, 0.50, 0,
		28999.93579, 29000.13579, 0,
	},
	{
		20000.86420, &Inventory{1100, 500, 0.02}, 1.75309, 0.35, 4,
		19999.96420, 20000.86420, 0,
	},
	{
		26500.97531, &Inventory{1400, 300, 0.03}, 0.97531, 0.60, 2,
		26499.87531, 26500.97531, 0,
	},
	{
		23500.64209, &Inventory{600, 600, 0.02}, 1.64209, 0.40, 3,
		23499.94209, 23500.64209, 0,
	},
	{
		24500.75309, &Inventory{1000, 200, 0.01}, 1.75309, 0.55, 1,
		24499.85309, 24500.75309, 0,
	},
	{
		20500.86420, &Inventory{1300, 350, 0.03}, 0.86420, 0.50, 2,
		20499.76420, 20500.86420, 0,
	},
	{
		25500.97531, &Inventory{700, 450, 0.02}, 0.97531, 0.35, 0,
		25499.87531, 25500.97531, 0,
	},
	{
		22500.08642, &Inventory{900, 550, 0.02}, 1.08642, 0.60, 4,
		22499.98642, 22500.08642, 0,
	},
	{
		21500.19753, &Inventory{1100, 300, 0.03}, 1.19753, 0.40, 2,
		21499.89753, 21500.19753, 0,
	},
	{
		28500.30864, &Inventory{1000, 400, 0.01}, 1.30864, 0.50, 3,
		28499.90864, 28500.30864, 0,
	},
	{
		20500.41975, &Inventory{800, 400, 0.02}, 1.41975, 0.35, 1,
		20499.91975, 20500.41975, 0,
	},	
}



func TestOptimizeSpread(t *testing.T) {
	for _, tt := range testCases {
		t.Run("", func(t *testing.T) {
			optimalBid, optimalAsk, primalStatus := OptimizeSpread(
				tt.currentPrice,
				tt.inventory,
				tt.volatility,
				tt.liquidity,
				tt.orderBookDepth,
			)

			if optimalBid != tt.expectedBid {
				t.Errorf("Expected Optimal Bid %f, got %f", tt.expectedBid, optimalBid)
			}

			if optimalAsk != tt.expectedAsk {
				t.Errorf("Expected Optimal Ask %f, got %f", tt.expectedAsk, optimalAsk)
			}

			if primalStatus != tt.expectedPrimal {
				t.Errorf("Expected Primal Status %d, got %d", tt.expectedPrimal, primalStatus)
			}

			fmt.Printf("Optimal Bid: %f, Optimal Ask: %f, Price: %f\n", optimalBid, optimalAsk, tt.currentPrice)
		})
	}
}
