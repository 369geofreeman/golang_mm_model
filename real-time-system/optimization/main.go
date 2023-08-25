package optimization

import (
	"errors"
	"log"
	"math"

	"github.com/lanl/clp"
)

// Parameters for our optimization model
var (
	alpha = 0.05
	beta  = 1.0
	gamma = 1.0
	delta = 0.5
	zeta  = 0.1

	// Variables to store the last successful bid and ask
	lastSuccessfulBid float64
	lastSuccessfulAsk float64

	isEmaInitialized bool    = false
	emaVolatility    float64 = 0
	emaFactor        float64 = 0.1  // Defines the sensitivity of the EMA. Value between 0 and 1.
	lowVolThreshold          = 0.5  // Example value, adjust based on observed volatility
	highVolThreshold         = 2.0  // Example value, adjust based on observed volatility
	minEmaFactor             = 0.05 // Minimum responsiveness
	maxEmaFactor             = 0.5  // Maximum responsiveness
)

// SetParameters allows updating of optimization parameters dynamically
func SetParameters(newAlpha, newBeta, newGamma, newDelta, newZeta float64) {
	alpha = newAlpha
	beta = newBeta
	gamma = newGamma
	delta = newDelta
	zeta = newZeta
}

// GetParameters returns the current values of the optimization parameters
func GetParameters() (float64, float64, float64, float64, float64) {
	return alpha, beta, gamma, delta, zeta
}

// validateParameters checks if the provided parameters are within expected ranges.
func validateParameters(currentPrice, maxInventory, volatility, liquidity, orderBookDepth float64) error {
	if currentPrice <= 0 {
		return errors.New("currentPrice should be greater than 0")
	}
	if maxInventory <= 0 {
		return errors.New("maxInventory should be greater than 0")
	}
	if volatility < 0 {
		return errors.New("volatility should not be negative")
	}
	if liquidity <= 0 {
		return errors.New("liquidity should be greater than 0")
	}
	if orderBookDepth < 0 {
		return errors.New("orderBookDepth should not be negative")
	}
	return nil
}

// OptimizeSpread calculates the optimal bid and ask prices based on market conditions
func OptimizeSpread(currentPrice float64, inventory *Inventory, volatility, liquidity, orderBookDepth float64) (float64, float64, clp.SimplexStatus) {
	// Fetch the current cash balance from the inventory
	maxInventory, _ := inventory.GetBalances()
	// Validate parameters before proceeding
	if err := validateParameters(currentPrice, maxInventory, volatility, liquidity, orderBookDepth); err != nil {
		log.Println("Warning:", err)
		// Return a default spread around the current price as a fallback
		defaultSpread := 0.5 // Default to a spread of 0.5, this can be adjusted
		return currentPrice + defaultSpread, currentPrice - defaultSpread, 0
	}

	// Initialize EMA with the first volatility value received
	if !isEmaInitialized {
		emaVolatility = volatility
		isEmaInitialized = true
	} else {
		emaVolatility = (1-emaFactor)*emaVolatility + emaFactor*volatility
	}

	// Using CLP to solve the LP problem
	simp := clp.NewSimplex()

	// Objective coefficients (minimize the sum of cost functions for bid and ask)
	// Weights to prioritize alignment with current price
	// Fetch the current cash balance from the inventory
	cash, assets := inventory.GetBalances()
	totalValue := cash + assets*currentPrice
	assetRatio := assets * currentPrice / totalValue

	// Define percentage-based deviations for bid and ask
	bidDeviationPercentage := 0.01 // 0.05 e.g., 5% below the current price
	askDeviationPercentage := 0.01 // 0.05 e.g., 5% above the current price

	// Calculate absolute deviations
	bidDeviation := currentPrice * bidDeviationPercentage
	askDeviation := currentPrice * askDeviationPercentage

	// Initial guesses for bid and ask
	// initialBid := currentPrice - bidDeviation
	// initialAsk := currentPrice + askDeviation

	// Dynamic weights to prioritize alignment with current price
	c := objectiveFunction(currentPrice, volatility, liquidity, orderBookDepth, assetRatio)


	// Define variable bounds
	minSpread := currentPrice * 0.01                   // 0.01 Minimum allowable spread between bid and ask
	varBounds := [][2]float64{
		{currentPrice - bidDeviation, currentPrice}, // Bounds for Bid Price
		{currentPrice, currentPrice + askDeviation}, // Bounds for Ask Price
	}

	// Constraints
	ineqs := [][]float64{
		{0, 1, -1, -minSpread}, // Ensure bid is less than ask by at least minSpread
	}
	
	simp.EasyLoadDenseProblem(c, varBounds, ineqs)
	simp.SetOptimizationDirection(clp.Minimize)
	test_primal := simp.Primal(clp.NoValuesPass, clp.NoStartFinishOptions)

	// Fetching the solution
	optX := simp.PrimalColumnSolution()

	// Extract the optimized bid and ask prices
	optimalBid := optX[0]
	optimalAsk := optX[1]

	// Print intermediate values for debugging
	log.Println("currentPrice:", currentPrice)
	log.Println("inventory:", inventory)
	log.Println("volatility:", volatility)
	log.Println("liquidity:", liquidity)
	log.Println("orderBookDepth:", orderBookDepth)
	log.Println("assetRatio:", assetRatio)
	log.Println("cash:", cash)
	log.Println("assets:", assets)
	log.Println("totalValue:", totalValue)
	log.Println("Variable Bounds:", varBounds)
	log.Println("minSpread:", minSpread)
	log.Println("Test primal", test_primal)
	log.Println("optimalBid:", optimalBid)
	log.Println("optimalAsk", optimalAsk)
	log.Println("optX", optX)

	return optimalBid, optimalAsk, test_primal
}

func GetOptimizationFrequency() int {
	if emaVolatility > 1.5 { // Thresholds can be adjusted based on your needs
		return 1 // Optimize every minute
	} else if emaVolatility > 1.0 {
		return 5 // Optimize every 5 minutes
	} else {
		return 10 // Optimize every 10 minutes
	}
}

// Helper functions to get the minimum and maximum of two floats
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Objective function to define the coefficients for the optimization problem
func objectiveFunction(currentPrice, volatility, liquidity, orderBookDepth float64, assetRatio float64) []float64 {
	volatilityPenalty := alpha * volatility
	liquidityPenalty := beta / (liquidity + 1) // Assuming liquidity is always positive
	orderBookPenalty := gamma * math.Log(1 + orderBookDepth)

	return []float64{
		-1 + volatilityPenalty + liquidityPenalty + orderBookPenalty,
		1 - volatilityPenalty - liquidityPenalty - orderBookPenalty,
	}
}


func costFunction(vA, vB, a, b, currentPrice, volatility float64) float64 {
	inventoryRisk := alpha * math.Pow(vA-vB, 2) // Penalize imbalances between vA and vB
	priceRisk := beta * (math.Pow(a-currentPrice, 2) + math.Pow(b-currentPrice, 2))
	return inventoryRisk + priceRisk
}

func baseSpreadFunction(volatility, liquidity, orderBookDepth float64) float64 {
	// Adjusted the coefficients to make the spread more sensitive to market conditions
	return 2*gamma*volatility + delta/(liquidity+1) + 2*zeta*math.Log(1+orderBookDepth)
}

func AdjustEmaFactorBasedOnVolatility(volatility float64) {
	if volatility <= lowVolThreshold {
		emaFactor = minEmaFactor
	} else if volatility >= highVolThreshold {
		emaFactor = maxEmaFactor
	} else {
		// Linear scaling between the thresholds
		scale := (volatility - lowVolThreshold) / (highVolThreshold - lowVolThreshold)
		emaFactor = minEmaFactor + scale*(maxEmaFactor-minEmaFactor)
	}
}
