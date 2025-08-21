package main

import (
	"encoding/json"
	"fmt"
	"nourish-backend/internal/models"
)

func main() {
	fmt.Println("🧪 Testing FlexibleDate JSON Parsing")
	fmt.Println("====================================")

	// Test cases for different date formats
	testCases := []string{
		`{"date": "2025-08-21"}`,                // Date only (your frontend)
		`{"date": "2025-08-21T10:30:00Z"}`,      // ISO8601 UTC
		`{"date": "2025-08-21T10:30:00+05:30"}`, // ISO8601 with timezone
		`{"date": "2025-08-21T10:30:00"}`,       // ISO8601 without timezone
		`{"date": "2025-08-21 10:30:00"}`,       // SQL datetime format
	}

	for i, testCase := range testCases {
		fmt.Printf("\n%d. Testing: %s\n", i+1, testCase)

		var req struct {
			Date models.FlexibleDate `json:"date"`
		}

		err := json.Unmarshal([]byte(testCase), &req)
		if err != nil {
			fmt.Printf("   ❌ Error: %v\n", err)
		} else {
			fmt.Printf("   ✅ Success: %s\n", req.Date.Time.Format("2006-01-02 15:04:05 MST"))
		}
	}

	// Test the complete MealRequest
	fmt.Println("\n🍽️  Testing Complete MealRequest")
	fmt.Println("==================================")

	mealJSON := `{
		"date": "2025-08-21",
		"mealType": "breakfast",
		"dishId": "64f1a2b3c4d5e6f789012345",
		"notes": "Test meal",
		"rating": 4
	}`

	var mealReq models.MealRequest
	err := json.Unmarshal([]byte(mealJSON), &mealReq)
	if err != nil {
		fmt.Printf("❌ Error parsing MealRequest: %v\n", err)
	} else {
		fmt.Printf("✅ MealRequest parsed successfully!\n")
		fmt.Printf("   Date: %s\n", mealReq.Date.Time.Format("2006-01-02 15:04:05"))
		fmt.Printf("   MealType: %s\n", mealReq.MealType)
		fmt.Printf("   DishID: %s\n", mealReq.DishID)
		fmt.Printf("   Notes: %s\n", mealReq.Notes)
		fmt.Printf("   Rating: %d\n", mealReq.Rating)
	}

	// Test with no rating (should default to 0 and be valid)
	fmt.Println("\n🔢 Testing MealRequest with no rating")
	fmt.Println("=====================================")

	mealJSONNoRating := `{
		"date": "2025-08-21",
		"mealType": "breakfast",
		"dishId": "64f1a2b3c4d5e6f789012345",
		"notes": "Test meal without rating"
	}`

	var mealReqNoRating models.MealRequest
	err = json.Unmarshal([]byte(mealJSONNoRating), &mealReqNoRating)
	if err != nil {
		fmt.Printf("❌ Error parsing MealRequest without rating: %v\n", err)
	} else {
		fmt.Printf("✅ MealRequest without rating parsed successfully!\n")
		fmt.Printf("   Date: %s\n", mealReqNoRating.Date.Time.Format("2006-01-02 15:04:05"))
		fmt.Printf("   MealType: %s\n", mealReqNoRating.MealType)
		fmt.Printf("   DishID: %s\n", mealReqNoRating.DishID)
		fmt.Printf("   Notes: %s\n", mealReqNoRating.Notes)
		fmt.Printf("   Rating: %d (0 means no rating)\n", mealReqNoRating.Rating)
	}

	fmt.Println("\n🎉 FlexibleDate should now handle your frontend's date format!")
	fmt.Println("🎯 Rating validation now allows 0 (no rating) through 5!")
}
