package product

import "promotion-management-api/internal/db"

func List() ([]Product, error) {
	db := db.GetConnection()

	results, err := db.Query("SELECT `ID`, `Name`, `Price`, `Type`, `Description` FROM `product`")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	products := make([]Product, 0)

	for results.Next() {
		var product Product

		err = results.Scan(&product.Id, &product.Name, &product.Price, &product.Type, &product.Description)
		if err != nil {
			return nil, err
		}

		products = append(products, product)

	}

	return products, nil
}
