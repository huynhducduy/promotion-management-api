package employee

import (
	"database/sql"
	"errors"
	"promotion-management-api/internal/db"
)

func List() ([]Employee, error) {
	db := db.GetConnection()

	results, err := db.Query("SELECT `ID`, `Name`, `Phone`, `Address`, `JoinedDate`, `Username`, `Email` FROM `employee`")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	emps := make([]Employee, 0)

	for results.Next() {
		var emp Employee

		err = results.Scan(&emp.Id, &emp.Name, &emp.Phone, &emp.Address, &emp.JoinedDate, &emp.Username, &emp.Email)
		if err != nil {
			return nil, err
		}

		emps = append(emps, emp)

	}

	return emps, nil
}

func Read(id int64) (*Employee, error) {
	var emp Employee
	db := db.GetConnection()

	results := db.QueryRow("SELECT `ID`, `Name`, `Phone`, `Address`, `JoinedDate`, `Username`, `Email` FROM `employee` WHERE `id` = ? ", id)
	err := results.Scan(&emp.Id, &emp.Name, &emp.Phone, &emp.Address, &emp.JoinedDate, &emp.Username, &emp.Email)
	if err == sql.ErrNoRows {
		return nil, errors.New("Invalid id.")
	} else if err != nil {
		return nil, err
	}

	return &emp, nil
}
