package handlers

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nehjoshi/fetch-challenge/types"
)

// In-memory storage solution to store points for each receipt
var Scores = make(map[string]int)

// Function to generate a random UUID for each valid receipt
func generateReceiptId() string {
	id := uuid.New()
	return id.String()
}

// Function to calculate retailer points
func calculateRetailerScore(name string) int {
	points := 0
	//One point for every alphanumeric character in the retailer name.
	for _, c := range name {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			points++
		}
	}
	return points
}

// Function to calculate price points
func calculateDollarScore(dollarsRaw string) (int, error) {
	points := 0
	dollars, err := strconv.ParseFloat(dollarsRaw, 32)
	if err != nil {
		return 0, errors.New("Error parsing dollar amount to float. Make sure you don't include special symbols.")
	}

	//50 points if the total is a round dollar amount with no cents.
	floored := math.Floor(dollars)
	if floored == dollars {
		points += 50
	}
	//25 points if the total is a multiple of 0.25.
	if math.Mod(dollars, 0.25) == 0 {
		points += 25
	}
	return points, nil
}

// Function to calculate pairwise item points
func calculatePairwiseScore(items []types.Item) int {
	points := 0
	numberOfItems := len(items)
	pairs := numberOfItems / 2

	//5 points for every two items on the receipt.
	points += pairs * 5
	return points
}

// Function to calculate trimmed description points
func calculateDescriptionScore(items []types.Item) (int, error) {
	points := 0
	//If the trimmed length of the item description is a multiple of 3,
	//multiply the price by 0.2 and round up to the nearest integer.
	for _, v := range items {
		trimmedDesc := strings.TrimSpace(v.ShortDescription)
		trimmedLength := len(trimmedDesc)
		if trimmedLength%3 == 0 {
			priceFloat, err := strconv.ParseFloat(v.Price, 32)
			if err == nil {
				points += int(math.Ceil(priceFloat * 0.2))
			} else {
				return 0, errors.New("Error converting item price to float. Make sure the price doesn't include any special symbols.")
			}
		}
	}
	return points, nil
}

// Function to calculate date points
func calculateDateScore(date string) (int, error) {

	if len(date) > 10 || len(date) == 0 {
		return 0, errors.New("Error formatting date. Date must be of format YYYY-MM-DD")
	}

	points := 0
	day := date[len(date)-1:]
	intDay, err := strconv.Atoi(day)
	if err != nil {
		return 0, errors.New("Error formatting date. Date must be of format YYYY-MM-DD")
	}
	if err == nil && intDay%2 != 0 {
		points += 6
	}
	return points, nil
}

// Function to calculate time points
func calculateTimeScore(time string) (int, error) {

	if len(time) > 5 || len(time) == 0 {
		return 0, errors.New("Error formatting time. Time must be of format HH:MM")
	}

	//10 points if the time of purchase is after 2:00pm and before 4:00pm.
	hours, errHours := strconv.Atoi(time[0:2])
	minutes, errMinutes := strconv.Atoi(time[3:])

	if errHours != nil || errMinutes != nil {
		return 0, errors.New("Error formatting time. Time must be of format HH:MM")
	}
	points := 0
	if (hours > 14 || hours == 14 && minutes > 0) && hours < 16 {
		points += 10
	}
	return points, nil
}

// Function to get receipt score or points
func getReceiptScore(rct types.Receipt) (int, error) {
	var points int = 0

	//Get retailer points
	points += calculateRetailerScore(rct.Retailer)

	//Get dollar points
	scoreDollar, errDollar := calculateDollarScore(rct.Total)
	if errDollar != nil {
		return 0, errDollar
	}
	points += scoreDollar

	//Get pairwise item points
	points += calculatePairwiseScore(rct.Items)

	//Get trimmed description points
	scoreDesc, errDesc := calculateDescriptionScore(rct.Items)
	if errDesc != nil {
		return 0, errDesc
	}
	points += scoreDesc

	//Get date points
	scoreDate, errDate := calculateDateScore(rct.PurchaseDate)
	if errDate != nil {
		return 0, errDate
	}
	points += scoreDate

	//Get time points
	scoreTime, errTime := calculateTimeScore(rct.PurchaseTime)
	if errTime != nil {
		return 0, errTime
	}
	points += scoreTime

	return points, nil
}

func ProcessReceipt(c *gin.Context) {
	var receipt types.Receipt
	if err := c.BindJSON(&receipt); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"description": "The receipt is invalid"})
		return
	}
	id := generateReceiptId()
	score, err := getReceiptScore(receipt)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	Scores[id] = score
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}
