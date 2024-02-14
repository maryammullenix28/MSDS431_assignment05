package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gonum/stat/combin"
	"github.com/sajari/regression"
)

type Result struct {
	Combination []string
	Formula     string
	MSE         float64
	AIC         float64
	BIC         float64
}

func main() {
	// Open the CSV file from the disk
	f, err := os.Open("../boston.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a new CSV reader specifying the number of columns it has
	housingData := csv.NewReader(f)

	// Read all the records
	records, err := housingData.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	//Data exploration
	fmt.Println("-----------------------------------------------")
	fmt.Println("-------------- DATA EXPLORATION ---------------")
	fmt.Println("-----------------------------------------------")
	fmt.Printf("The data set contains the following columns:%s\n", records[0])

	//Calculate average of mv
	totalPrice := 0.0
	n := len(records) - 1
	for i, record := range records {
		// Skip the header.
		if i == 0 {
			continue
		}

		// Parse the house price, "mv".
		price, err := strconv.ParseFloat(records[i][len(record)-1], 64) // mv is the last column
		if err != nil {
			log.Fatal(err)
		}

		// Accumulate total price
		totalPrice += price
	}

	avgPrice := totalPrice / float64(n)
	fmt.Printf("Average House Price (mv): %.2f\n\n", avgPrice)

	// Define the variables for linear regression
	variables := []string{"neighborhood", "mv", "nox", "crim", "zn", "indus", "chas", "rooms", "age", "dis", "rad", "tax", "ptratio", "lstat"}

	// Generate combinations of four or more explanatory variables
	explanatoryCombinations := generateCombinations(variables[2:len(variables)-2], 4)
	fmt.Printf("Number of combinations: %d\n\n", len(explanatoryCombinations))

	fmt.Println("-----------------------------------------------")
	fmt.Println("----------------- COMBINATIONS ----------------")
	fmt.Println("-----------------------------------------------")

	results := make(chan Result)
	var wg sync.WaitGroup

	// Function to handle processing for each combination
	processCombination := func(combination []string) {
		defer wg.Done()

		var r regression.Regression
		r.SetObserved("mv")
		//fmt.Printf("We're looking at the following variables: %s\n\n", combination)
		index := 0

		// Set variables for the current combination
		for _, varName := range combination {
			r.SetVar(index, varName)
			index++
		}

		// Loop over records in the CSV, adding the training data to the regression value.
		for i, record := range records {
			// Skip the header
			if i == 0 {
				continue
			}

			// Parse mv
			price, err := strconv.ParseFloat(records[i][len(variables)-1], 64) // mv is the last column
			if err != nil {
				log.Fatal(err)
			}

			// Parse the explanatory variable values for the current combination.
			var values []float64
			for _, varName := range combination {
				value, err := strconv.ParseFloat(record[getIndex(variables, varName)], 64)
				if err != nil {
					log.Fatal(err)
				}
				values = append(values, value)
			}

			// Add these points to the regression value.
			r.Train(regression.DataPoint(price, values))
		}

		// Train/fit the regression model.
		r.Run()

		// Compute Mean Square Error (MSE)
		mse := computeMSE(&r, records, variables, combination)

		// Compute Information Criterion (e.g., AIC or BIC)
		aic := computeAIC(&r, mse, len(combination)+1, len(records)-1) // +1 for intercept
		bic := computeBIC(&r, mse, len(combination)+1, len(records)-1) // +1 for intercept

		// Send results to the channel
		results <- Result{Combination: combination, Formula: r.Formula, MSE: mse, AIC: aic, BIC: bic}
	}

	// Start goroutines to process each combination
	for _, combination := range explanatoryCombinations {
		wg.Add(1)
		go processCombination(combination)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Receive and print results
	for result := range results {
		fmt.Printf("\nCombination: %v\nFormula: %s\nMSE: %f\nAIC: %f\nBIC: %f\n", result.Combination, result.Formula, result.MSE, result.AIC, result.BIC)
	}
}

// getIndex returns the index of a given string in a slice of strings
func getIndex(slice []string, str string) int {
	for i, v := range slice {
		if strings.ToLower(v) == strings.ToLower(str) {
			return i
		}
	}
	return -1
}

// generateCombinations generates combinations of variables
func generateCombinations(variables []string, r int) [][]string {
	combinations := combin.Combinations(len(variables), r)
	result := make([][]string, 0, combin.Binomial(len(variables), r))
	for _, combo := range combinations {
		var combination []string
		for _, idx := range combo {
			combination = append(combination, variables[idx])
		}
		result = append(result, combination)
	}
	return result
}

// computeMSE computes the Mean Square Error (MSE) for a given regression model
func computeMSE(r *regression.Regression, records [][]string, variables []string, combination []string) float64 {
	var mse float64
	n := len(records) - 1 // Number of observations

	for i, record := range records {
		// Skip the header and the neighborhood column.
		if i == 0 {
			continue
		}

		// Parse the house price, "y".
		price, err := strconv.ParseFloat(records[i][len(variables)-1], 64) // mv is the last column
		if err != nil {
			log.Fatal(err)
		}

		// Parse the explanatory variable values for the current combination.
		var values []float64
		for _, varName := range combination {
			value, err := strconv.ParseFloat(record[getIndex(variables, varName)], 64)
			if err != nil {
				log.Fatal(err)
			}
			values = append(values, value)
		}

		// Predict using the regression model
		prediction, _ := r.Predict(values)

		// Compute the squared error
		squaredError := math.Pow(price-prediction, 2)
		mse += squaredError
	}

	// Calculate the mean square error
	mse /= float64(n)
	return mse
}

// computeAIC computes the Akaike Information Criterion (AIC)
func computeAIC(r *regression.Regression, mse float64, k int, n int) float64 {
	aic := float64(n)*math.Log(mse) + 2*float64(k)
	return aic
}

// computeBIC computes the Bayesian Information Criterion (BIC)
func computeBIC(r *regression.Regression, mse float64, k int, n int) float64 {
	bic := float64(n)*math.Log(mse) + float64(k)*math.Log(float64(n))
	return bic
}
