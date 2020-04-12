package promotion

import (
	"database/sql"
	"errors"
	"promotion-management-api/internal/db"
	"promotion-management-api/internal/store"
)

func list() ([]Promotion, error) {
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

func create(promotion PromotionExtra) (int64, error) {
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

func read(id int64) (*PromotionExtra, error) {

	promotion := PromotionExtra{}
	promotion.Promotion = new(Promotion)
	promotion.Store = new([]store.Store)
	promotion.Time = new([]Time)

	db := db.GetConnection()

	results := db.QueryRow("SELECT `ID`, `Name`, `StartDate`, `EndDate`, `MainGoal`, `ApplyingType`, `ApplyingForm`, `ApplyingValue` FROM `promotion` WHERE `id` = ? ", id)
	err := results.Scan(&promotion.Id, &promotion.Name, &promotion.StartDate, &promotion.EndDate, &promotion.MainGoal, &promotion.ApplyingType, &promotion.ApplyingForm, &promotion.ApplyingValue)
	if err == sql.ErrNoRows {
		return nil, errors.New("Invalid id.")
	} else if err != nil {
		return nil, err
	}

	results2, err := db.Query("SELECT `ID`, `Name`, `ManagerID`, `Address` FROM `store` WHERE `ID` IN (SELECT `StoreID` FROM `store_constraint` WHERE `PromotionID` = ?)", id)
	if err != nil {
		return nil, err
	}
	defer results2.Close()

	stores := make([]store.Store, 0)

	for results2.Next() {
		var store store.Store

		err = results2.Scan(&store.Id, &store.Name, &store.ManagerID, &store.Address)
		if err != nil {
			return nil, err
		}

		stores = append(stores, store)
	}

	promotion.Store = &stores

	results2, err = db.Query("SELECT `Type` FROM `payment_constraint` WHERE `PromotionID` = ?", id)
	if err != nil {
		return nil, err
	}
	defer results2.Close()

	payments := make([]string, 0)

	for results2.Next() {
		var payment string

		err = results2.Scan(&payment)
		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	promotion.Payment = &payments

	results2, err = db.Query("SELECT `Type`, `StartTime`, `EndTime` FROM `time_constraint` WHERE `PromotionID` = ?", id)
	if err != nil {
		return nil, err
	}
	defer results2.Close()

	times := make([]Time, 0)

	for results2.Next() {
		var time Time

		err = results2.Scan(&time.Type, &time.StartTime, &time.EndTime)
		if err != nil {
			return nil, err
		}

		times = append(times, time)
	}

	promotion.Time = &times

	results2, err = db.Query("SELECT `Type` FROM `order_type_constraint` WHERE `PromotionID` = ?", id)
	if err != nil {
		return nil, err
	}
	defer results2.Close()

	orderTypes := make([]string, 0)

	for results2.Next() {
		var orderType string

		err = results2.Scan(&orderType)
		if err != nil {
			return nil, err
		}

		orderTypes = append(orderTypes, orderType)
	}

	promotion.OrderType = &orderTypes

	results2, err = db.Query("SELECT `Type` FROM `membership_constraint` WHERE `PromotionID` = ?", id)
	if err != nil {
		return nil, err
	}
	defer results2.Close()

	memberships := make([]string, 0)

	for results2.Next() {
		var membership string

		err = results2.Scan(&membership)
		if err != nil {
			return nil, err
		}

		memberships = append(memberships, membership)
	}

	promotion.Membership = &memberships

	return &promotion, err
}

func update(updatedPromo PromotionExtra) (*PromotionExtra, error) {
	db := db.GetConnection()

	_, err := db.Exec("UPDATE `promotion` SET `Name` = ?, `StartDate` = ?, `EndDate` = ?, `MainGoal` = ?, `ApplyingType` = ?, `ApplyingForm` = ?, `ApplyingValue` = ? WHERE `ID` = ?", updatedPromo.Name, updatedPromo.StartDate, updatedPromo.EndDate, updatedPromo.MainGoal, updatedPromo.ApplyingType, updatedPromo.ApplyingForm, updatedPromo.ApplyingValue, updatedPromo.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("DELETE FROM `store_constraint` WHERE `PromotionID` = ?;", updatedPromo.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("DELETE FROM `time_constraint` WHERE `PromotionID` = ?;", updatedPromo.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("DELETE FROM `payment_constraint` WHERE `PromotionID` = ?;", updatedPromo.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("DELETE FROM `membership_constraint` WHERE `PromotionID` = ?;", updatedPromo.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("DELETE FROM `order_type_constraint` WHERE `PromotionID` = ?;", updatedPromo.Id)
	if err != nil {
		return nil, err
	}

	if updatedPromo.Store != nil {
		for _, i := range *updatedPromo.Store {
			_, err := db.Exec("INSERT INTO `store_constraint` (`PromotionID`, `StoreID`) VALUES (?,?)", updatedPromo.Id, i.Id)
			if err != nil {
				return nil, err
			}
		}
	}

	if updatedPromo.Payment != nil {
		for _, i := range *updatedPromo.Payment {
			_, err := db.Exec("INSERT INTO `payment_constraint` (`PromotionID`, `Type`) VALUES (?,?)", updatedPromo.Id, i)
			if err != nil {
				return nil, err
			}
		}
	}

	if updatedPromo.Membership != nil {
		for _, i := range *updatedPromo.Membership {
			_, err := db.Exec("INSERT INTO `membership_constraint` (`PromotionID`, `Type`) VALUES (?,?)", updatedPromo.Id, i)
			if err != nil {
				return nil, err
			}
		}
	}

	if updatedPromo.OrderType != nil {
		for _, i := range *updatedPromo.OrderType {
			_, err := db.Exec("INSERT INTO `order_type_constraint` (`PromotionID`, `Type`) VALUES (?,?)", updatedPromo.Id, i)
			if err != nil {
				return nil, err
			}
		}
	}

	if updatedPromo.Time != nil {
		for _, i := range *updatedPromo.Time {
			_, err := db.Exec("INSERT INTO `time_constraint` (`PromotionID`, `Type`, `StartTime`, `EndTime`) VALUES (?,?,?,?)", updatedPromo.Id, i.Type, i.StartTime, i.EndTime)
			if err != nil {
				return nil, err
			}
		}
	}

	return &updatedPromo, nil
}

func delete(id int64) error {
	db := db.GetConnection()

	_, err := db.Exec("DELETE FROM `time_constraint` WHERE `PromotionID` = ?", id)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM `store_constraint` WHERE `PromotionID` = ?", id)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM `payment_constraint` WHERE `PromotionID` = ?", id)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM `membership_constraint` WHERE `PromotionID` = ?", id)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM `order_type_constraint` WHERE `PromotionID` = ?", id)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM `promotion` WHERE `ID` = ?", id)
	return err
}
