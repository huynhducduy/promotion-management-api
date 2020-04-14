package employee

import (
	"net/http"
	"promotion-management-api/pkg/utils"
)

func RouterList(w http.ResponseWriter, r *http.Request) {
	data, err := List()
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.Response(w, 200, data)
}

func RouterMe(w http.ResponseWriter, r *http.Request) {
	employee := r.Context().Value("employee").(*Employee)
	utils.Response(w, 200, employee)
}