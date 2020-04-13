package promotion

import (
	"database/sql"
	"errors"
	"promotion-management-api/internal/db"
	"promotion-management-api/internal/store"
	"promotion-management-api/pkg/utils"
	"time"
)

func List() ([]Promotion, error) {
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

func Create(promotion PromotionExtra) (int64, error) {
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

func Read(id int64) (*PromotionExtra, error) {

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

func Applicable(storeId int64, memberId int64, paymentType string, orderType string) ([]Promotion, error) {
	promotions := make([]Promotion, 0)

	db := db.GetConnection()

	// Check Store -----------------------------------------------------------------------------------------------------
	promoIds := make([]interface{}, 0)

	results, err := db.Query("SELECT DISTINCT `PromotionID` FROM `store_constraint` WHERE `StoreID` = ? OR `StoreID` = \"0\"", storeId)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	for results.Next() {
		var promoId string
		err = results.Scan(&promoId)
		if err != nil {
			return nil, err
		}
		promoIds = append(promoIds, promoId)
	}

	utils.Logg("After store")
	utils.Logg(promoIds)

	if len(promoIds) == 0 {
		return promotions, nil
	}

	// Check Payment ---------------------------------------------------------------------------------------------------
	var queryString string
	var stuffs []interface{}

	stuffs = append(stuffs, promoIds...)
	queryString = "SELECT DISTINCT `PromotionID` FROM `payment_constraint` WHERE `PromotionID` IN (?"

	for i := 1; i <= len(promoIds)-1; i++ {
		queryString += ",?"
	}

	queryString += ") AND (`Type` = ? OR `Type` = \"0\")"

	stuffs = append(stuffs, paymentType)

	results, err = db.Query(queryString, stuffs...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	promoIds = make([]interface{}, 0)

	for results.Next() {
		var promoId string
		err = results.Scan(&promoId)
		if err != nil {
			return nil, err
		}
		promoIds = append(promoIds, promoId)
	}

	utils.Logg("After payment")
	utils.Logg(promoIds)
	if len(promoIds) == 0 {
		return promotions, nil
	}

	// Check Order Type ------------------------------------------------------------------------------------------------

	queryString = ""
	stuffs = make([]interface{}, 0)

	stuffs = append(stuffs, promoIds...)
	queryString = "SELECT DISTINCT `PromotionID` FROM `order_type_constraint` WHERE `PromotionID` IN (?"

	for i := 1; i <= len(promoIds)-1; i++ {
		queryString += ",?"
	}

	queryString += ") AND (`Type` = ? OR `Type` = \"0\")"

	stuffs = append(stuffs, orderType)

	results, err = db.Query(queryString, stuffs...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	promoIds = make([]interface{}, 0)

	for results.Next() {
		var promoId string
		err = results.Scan(&promoId)
		if err != nil {
			return nil, err
		}
		promoIds = append(promoIds, promoId)
	}

	utils.Logg("After order type")
	utils.Logg(promoIds)
	if len(promoIds) == 0 {
		return promotions, nil
	}

	// Check time------------------------------------------------------------------------------------===----------------

	now := time.Now()
	hour := now.Hour()
	dayOfWeek := int(now.Weekday())

	queryString = ""
	stuffs = make([]interface{}, 0)

	stuffs = append(stuffs, promoIds...)
	queryString = "SELECT DISTINCT `PromotionID` FROM `time_constraint` WHERE `PromotionID` IN (?"

	for i := 1; i <= len(promoIds)-1; i++ {
		queryString += ",?"
	}

	queryString += ") AND ((`Type` = \"Hour\" AND `StartTime` <= ? AND `EndTime` >= ?) OR (`Type` = \"DayOfWeek\" AND `StartTime` <= ? AND `EndTime` >= ?) OR `Type` = \"0\")"

	stuffs = append(stuffs, hour, hour, dayOfWeek, dayOfWeek)

	results, err = db.Query(queryString, stuffs...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	promoIds = make([]interface{}, 0)

	for results.Next() {
		var promoId string
		err = results.Scan(&promoId)
		if err != nil {
			return nil, err
		}
		promoIds = append(promoIds, promoId)
	}

	utils.Logg("After time")
	utils.Logg(promoIds)
	if len(promoIds) == 0 {
		return promotions, nil
	}

	// Check membership -----------------------------------------------------------------------------===----------------

	var birthdate time.Time
	var point int64

	results2 := db.QueryRow("SELECT `Birthdate`, `Point` FROM `member` WHERE `ID` = ?", memberId)
	err = results2.Scan(&birthdate, &point)
	if err == sql.ErrNoRows {
		return nil, errors.New("Invalid member id.")
	} else if err != nil {
		return nil, err
	}

	bd := false
	pt := false

	if birthdate.Day() == now.Day() && birthdate.Month() == now.Month() {
		bd = true
	}

	if point >= 10 {
		pt = true
	}

	queryString = ""
	stuffs = make([]interface{}, 0)

	stuffs = append(stuffs, promoIds...)
	queryString = "SELECT DISTINCT `PromotionID` FROM `membership_constraint` WHERE `PromotionID` IN (?"

	for i := 1; i <= len(promoIds)-1; i++ {
		queryString += ",?"
	}

	queryString += ")"

	if bd && pt {
		queryString += " AND `Type` IN (\"Birthday\", \"Point\", \"0\")"
	} else if bd {
		queryString += " AND `Type` IN (\"Birthday\", \"0\")"
	} else if pt {
		queryString += " AND `Type` IN (\"Point\", \"0\")"
	}

	results, err = db.Query(queryString, stuffs...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	promoIds = make([]interface{}, 0)

	for results.Next() {
		var promoId string
		err = results.Scan(&promoId)
		if err != nil {
			return nil, err
		}
		promoIds = append(promoIds, promoId)
	}

	utils.Logg("After membership")
	utils.Logg(promoIds)
	if len(promoIds) == 0 {
		return promotions, nil
	}

	// Get promotion -----------------------------------------------------------------------------===-------------------

	queryString = "SELECT `ID`, `Name`, `StartDate`, `EndDate`, `MainGoal`, `ApplyingType`, `ApplyingForm`, `ApplyingValue` FROM `promotion`  WHERE `StartDate` <= CURDATE() AND `EndDate` >= CURDATE() AND `ID` IN (?"

	for i := 1; i <= len(promoIds)-1; i++ {
		queryString += ",?"
	}

	queryString += ")"

	results, err = db.Query(queryString, promoIds...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

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

func Update(updatedPromo PromotionExtra) (*PromotionExtra, error) {
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

func Delete(id int64) error {
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
