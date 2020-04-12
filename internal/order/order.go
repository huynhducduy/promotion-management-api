package order

type Order struct {
	Id          *int64  `json:"id"`
	StoreId     *int64  `json:"store_id"`
	CashierId   *int64  `json:"cashier_id"`
	MemberId    *int64  `json:"member_id"`
	PromotionId *int64  `json:"promotion_id"`
	CreatedTime *string `json:"created_time"`
	Total       *int64  `json:"total"`
	PaymentType *string `json:"payment_type"`
	Type        *string `json:"type"`
}
