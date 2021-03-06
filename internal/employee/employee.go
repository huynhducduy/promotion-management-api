package employee

type Employee struct {
	Id         *int64  `json:"id,omitempty"`
	Username   *string `json:"username"`
	Password   *string `json:"password,omitempty"`
	Name       *string `json:"name"`
	Phone      *string `json:"phone"`
	Address    *string `json:"address"`
	JoinedDate *string `json:"joined_date"`
	Email	   *string `json:"email"`
}
