package store

type Store struct {
	Id        *int64  `json:"id"`
	Name      *string `json:"name"`
	ManagerID *int64  `json:"manager_id"`
	Address   *string `json:"address"`
}
