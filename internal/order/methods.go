package order

import "promotion-management-api/internal/db"

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
