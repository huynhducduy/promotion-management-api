package promotion

type Promotion struct {
	Id            *int64   `json:"id"`
	Name          *int64   `json:"name"`
	StartDate     *string  `json:"start_date"`
	EndDate       *string  `json:"end_date"`
	MainGoal      *string  `json:"main_goal"`
	ApplyingType  *string  `json:"main_goal"`
	ApplyingForm  *string  `json:"applying_form"`
	ApplyingValue *float64 `json:"applying_value"`
}
