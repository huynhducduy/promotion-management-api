package order

import (
	"promotion-management-api/internal/employee"
	"promotion-management-api/internal/member"
	"promotion-management-api/internal/product"
	"promotion-management-api/internal/promotion"
	"promotion-management-api/internal/store"
)

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
	Discount    *string `json:"discount"`
}

type ProductWithQuantity struct {
	*product.Product
	Quantity *int64 `json:"quantity"`
}

type OrderExtra struct {
	*Order
	Member    *member.Member         `json:"member"`
	Cashier   *employee.Employee     `json:"cashier"`
	Store     *store.Store           `json:"store"`
	Promotion *promotion.Promotion   `json:"promotion"`
	Product   *[]ProductWithQuantity `json:"product"`
}
