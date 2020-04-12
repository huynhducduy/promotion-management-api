package promotion

import (
	"promotion-management-api/internal/store"
)

type Promotion struct {
	Id            *int64   `json:"id"`
	Name          *string  `json:"name"`
	StartDate     *string  `json:"start_date"`
	EndDate       *string  `json:"end_date"`
	MainGoal      *string  `json:"main_goal"`
	ApplyingType  *string  `json:"applying_type"`
	ApplyingForm  *string  `json:"applying_form"`
	ApplyingValue *float64 `json:"applying_value"`
}

type Time struct {
	Type      *string `json:"type"`
	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`
}

type PromotionWithConstrains struct {
	*Promotion
	Store      *[]store.Store `json:"store_constraint"`
	Payment    *[]string      `json:"payment_constraint"`
	Membership *[]string      `json:"membership_constraint"`
	OrderType  *[]string      `json:"order_type_constraint"`
	Time       *[]Time        `json:"time_constraint"`
}
