package order

import (
	"database/sql"
	"errors"
	"promotion-management-api/internal/db"
	"promotion-management-api/internal/employee"
	"promotion-management-api/internal/member"
	"promotion-management-api/internal/product"
	"promotion-management-api/internal/promotion"
	"promotion-management-api/internal/store"
)

func List() ([]Order, error) {
	db := db.GetConnection()

	results, err := db.Query("SELECT `ID`, `StoreID`, `CashierID`, `MemberID`, `PromotionID`, `PaymentType`, `Type`, `CreatedTime`, `Discount`, `Total` FROM `order`")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	orders := make([]Order, 0)

	for results.Next() {
		var order Order

		err = results.Scan(&order.Id, &order.StoreId, &order.CashierId, &order.MemberId, &order.PromotionId, &order.PaymentType, &order.Type, &order.CreatedTime, &order.Discount, &order.Total)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)

	}

	return orders, nil
}

func Create(order OrderExtra) (int64, error) {
	db := db.GetConnection()

	results, err := db.Exec("INSERT INTO `order` (`StoreID`, `CashierID`, `MemberID`, `PromotionID`, `PaymentType`, `Type`, `CreatedTime`, `Discount`, `Total`) VALUES (?,?,?,?,?,?,?,?,?)", order.StoreId, order.CashierId, order.MemberId, order.PromotionId, order.PaymentType, order.Type, order.CreatedTime, order.Discount, order.Total)
	if err != nil {
		return 0, err
	}

	lid, err := results.LastInsertId()
	if err != nil {
		return 0, err
	}

	if order.Product != nil {
		for _, i := range *order.Product {
			_, err := db.Exec("INSERT INTO `product_of_order` (`OrderID`, `ProductID`, `Quantity`) VALUES (?,?,?)", lid, i.Id, i.Quantity)
			if err != nil {
				return 0, err
			}
		}
	}

	if order.MemberId != nil {
		point := *order.Total / 10000

		if order.PromotionId != nil {
			var applyingValue float64
			results := db.QueryRow("SELECT `ApplyingValue` FROM `promotion` WHERE `ID` = ? AND `ApplyingForm` = \"Customer point\"", order.PromotionId)
			err := results.Scan(&applyingValue)
			if err != sql.ErrNoRows && err != nil {
				return 0, err
			}

			if applyingValue != 0 {
				point -= *order.Discount / int64(applyingValue*1000)
			}
		}

		_, err = db.Exec("UPDATE `member` SET `point` = `point` + ?", point)
		if err != nil {
			return 0, err
		}
	}

	return lid, nil
}

func Read(id int64) (*OrderExtra, error) {

	order := OrderExtra{}
	order.Order = new(Order)
	order.Product = new([]ProductWithQuantity)

	db := db.GetConnection()

	results := db.QueryRow("SELECT `StoreID`, `CashierID`, `MemberID`, `PromotionID`, `PaymentType`, `Type`, `CreatedTime`, `Discount`, `Total` FROM `order` WHERE `id` = ? ", id)
	err := results.Scan(&order.StoreId, &order.CashierId, &order.MemberId, &order.PromotionId, &order.PaymentType, &order.Type, &order.CreatedTime, &order.Discount, &order.Total)
	if err == sql.ErrNoRows {
		return nil, errors.New("Invalid id.")
	} else if err != nil {
		return nil, err
	}

	if order.PromotionId != nil {
		order.Promotion = new(promotion.Promotion)
		results := db.QueryRow("SELECT `ID`, `Name`, `StartDate`, `EndDate`, `MainGoal`, `ApplyingType`, `ApplyingForm`, `ApplyingValue` FROM `promotion` WHERE `id` = ? ", order.PromotionId)
		err := results.Scan(&order.Promotion.Id, &order.Promotion.Name, &order.Promotion.StartDate, &order.Promotion.EndDate, &order.Promotion.MainGoal, &order.Promotion.ApplyingType, &order.Promotion.ApplyingForm, &order.Promotion.ApplyingValue)
		if err == sql.ErrNoRows {
			return nil, errors.New("Invalid promotion id.")
		} else if err != nil {
			return nil, err
		}
	}

	if order.StoreId != nil {
		order.Store = new(store.Store)
		results := db.QueryRow("SELECT `ID`, `Name`, `ManagerID`, `Address` FROM `store` WHERE `id` = ? ", order.StoreId)
		err := results.Scan(&order.Store.Id, &order.Store.Name, &order.Store.ManagerID, &order.Store.Address)
		if err == sql.ErrNoRows {
			return nil, errors.New("Invalid store id.")
		} else if err != nil {
			return nil, err
		}
	}

	if order.MemberId != nil {
		order.Member = new(member.Member)
		results := db.QueryRow("SELECT `ID`, `Name`, `Phone`, `Address`, `Birthdate`, `Point`, `JoinedDate` FROM `member` WHERE `id` = ? ", order.MemberId)
		err := results.Scan(&order.Member.Id, &order.Member.Name, &order.Member.Phone, &order.Member.Address, &order.Member.Birthday, &order.Member.Point, &order.Member.JoinedDate)
		if err == sql.ErrNoRows {
			return nil, errors.New("Invalid member id.")
		} else if err != nil {
			return nil, err
		}
	}

	if order.CashierId != nil {
		order.Cashier = new(employee.Employee)
		results := db.QueryRow("SELECT `ID`, `Name`, `Phone`, `Address`, `JoinedDate`, `Username` FROM `employee` WHERE `id` = ?", order.CashierId)
		err := results.Scan(&order.Cashier.Id, &order.Cashier.Name, &order.Cashier.Phone, &order.Cashier.Address, &order.Cashier.JoinedDate, &order.Cashier.Username)
		if err == sql.ErrNoRows {
			return nil, errors.New("Invalid cashier id.")
		} else if err != nil {
			return nil, err
		}
	}

	results2, err := db.Query("SELECT `ID`, `Name`, `Price`, `Type`, `Description`, `Quantity` FROM `product`, `product_of_order` WHERE `product_of_order`.`OrderID` = ? AND `product`.`ID` = `product_of_order`.`ProductID`", id)
	if err != nil {
		return nil, err
	}
	defer results2.Close()

	products := make([]ProductWithQuantity, 0)

	for results2.Next() {
		pro := ProductWithQuantity{new(product.Product), new(int64)}

		err = results2.Scan(&pro.Id, &pro.Name, &pro.Price, &pro.Type, &pro.Description, &pro.Quantity)
		if err != nil {
			return nil, err
		}

		products = append(products, pro)
	}

	order.Product = &products

	return &order, err
}
