package promotion

import (
	"promotion-management-api/internal/db"
)

func getAll() ([]Promotion, error) {
	db := db.GetConnection()

	results, err := db.Query("SELECT `ID`, `Name`, `StartDate`, `EndDate`, `MainGoal`, `ApplyingType`, `ApplyingForm`, `ApplyingValue` FROM `promotion`")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	promotions := make([]Promotion, 0)

	for results.Next() {
		var promotion Promotion

		err = results.Scan(&promotion.Id, &promotion.Name, &promotion.StartDate, &promotion.EndDate, &promotion.MainGoal, &promotion.ApplyingType, &promotion.ApplyingForm, &promotion.ApplyingValue)
		if err != nil {
			return nil, err
		}

		promotions = append(promotions, promotion)

	}

	return promotions, nil
}
