package store

import "promotion-management-api/internal/db"

func List() ([]Store, error) {
	db := db.GetConnection()

	results, err := db.Query("SELECT `ID`, `Name`, `ManagerID`, `Address` FROM `store`")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	stores := make([]Store, 0)

	for results.Next() {
		var store Store

		err = results.Scan(&store.Id, &store.Name, &store.ManagerID, &store.Address)
		if err != nil {
			return nil, err
		}

		stores = append(stores, store)

	}

	return stores, nil
}
