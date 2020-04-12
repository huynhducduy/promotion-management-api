package member

import "promotion-management-api/internal/db"

func List() ([]Member, error) {
	db := db.GetConnection()

	results, err := db.Query("SELECT `ID`, `Name`, `Phone`, `Address`, `Birthdate`, `Point`, `JoinedDate` FROM `member`")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	mems := make([]Member, 0)

	for results.Next() {
		var mem Member

		err = results.Scan(&mem.Id, &mem.Name, &mem.Phone, &mem.Address, &mem.Birthday, &mem.Point, &mem.JoinedDate)
		if err != nil {
			return nil, err
		}

		mems = append(mems, mem)

	}

	return mems, nil
}
