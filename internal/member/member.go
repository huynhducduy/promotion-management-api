package member

type Member struct {
	Id         *int64  `json:"id,omitempty"`
	Name       *string `json:"name"`
	Phone      *string `json:"phone"`
	Address    *string `json:"address"`
	Birthday   *string `json:"birthday"`
	Point      *int64  `json:"point"`
	JoinedDate *string `json:"joined_date"`
}
