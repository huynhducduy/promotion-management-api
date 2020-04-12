package product

type Product struct {
	Id          *int64  `json:"id,omitempty"`
	Name        *string `json:"name"`
	Price       *int64  `json:"price"`
	Type        *string `json:"type"`
	Description *string `json:"description"`
}
