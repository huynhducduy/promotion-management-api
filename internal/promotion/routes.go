package promotion

import (
	"net/http"
	"promotion-management-api/pkg/utils"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := getAll()
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.Response(w, 200, data)
}
