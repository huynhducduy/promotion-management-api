package employee

import (
	"database/sql"
	"errors"
	"promotion-management-api/internal/db"
)

func Read(id int64) (*Employee, error) {
	var emp Employee
	db := db.GetConnection()

	results := db.QueryRow("SELECT `ID`, `Name`, `Phone`, `Address`, `JoinedDate`, `Username` FROM `employee` WHERE `id` = ? ", id)
	err := results.Scan(&emp.Id, &emp.Name, &emp.Phone, &emp.Address, &emp.JoinedDate, &emp.Username)
	if err == sql.ErrNoRows {
		return nil, errors.New("Invalid id.")
	} else if err != nil {
		return nil, err
	}

	return &emp, nil
}
