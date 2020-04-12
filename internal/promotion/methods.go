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

func insert(promotion PromotionWithConstrains) (int64, error) {
	db := db.GetConnection()

	results, err := db.Exec("INSERT INTO `promotion` (`Name`, `StartDate`, `EndDate`, `MainGoal`, `ApplyingType`, `ApplyingForm`, `ApplyingValue`) VALUES (?,?,?,?,?,?,?)", promotion.Name, promotion.StartDate, promotion.EndDate, promotion.MainGoal, promotion.ApplyingType, promotion.ApplyingForm, promotion.ApplyingValue)
	if err != nil {
		return 0, err
	}

	lid, err := results.LastInsertId()
	if err != nil {
		return 0, err
	}

	if promotion.Store != nil {
		for _, i := range *promotion.Store {
			_, err := db.Exec("INSERT INTO `store_constraint` (`PromotionID`, `StoreID`) VALUES (?,?)", lid, i.Id)
			if err != nil {
				return 0, err
			}
		}
	}

	if promotion.Payment != nil {
		for _, i := range *promotion.Payment {
			_, err := db.Exec("INSERT INTO `payment_constraint` (`PromotionID`, `Type`) VALUES (?,?)", lid, i)
			if err != nil {
				return 0, err
			}
		}
	}

	if promotion.Membership != nil {
		for _, i := range *promotion.Membership {
			_, err := db.Exec("INSERT INTO `membership_constraint` (`PromotionID`, `Type`) VALUES (?,?)", lid, i)
			if err != nil {
				return 0, err
			}
		}
	}

	if promotion.OrderType != nil {
		for _, i := range *promotion.OrderType {
			_, err := db.Exec("INSERT INTO `order_type_constraint` (`PromotionID`, `Type`) VALUES (?,?)", lid, i)
			if err != nil {
				return 0, err
			}
		}
	}

	if promotion.Time != nil {
		for _, i := range *promotion.Time {
			_, err := db.Exec("INSERT INTO `time_constraint` (`PromotionID`, `Type`, `StartTime`, `EndTime`) VALUES (?,?,?,?)", lid, i.Type, i.StartTime, i.EndTime)
			if err != nil {
				return 0, err
			}
		}
	}

	return lid, nil
}

func getOne(id int) (PromotionWithConstrains, error) {

}
