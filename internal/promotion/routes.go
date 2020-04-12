package promotion

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
	"promotion-management-api/pkg/utils"
	"strconv"
	"time"
)

func PromotionContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			utils.ResponseMessage(w, http.StatusBadRequest, "Id must be an integer!")
			return
		}

		promotion, err := read(id)
		if err != nil {
			utils.ResponseInternalError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), "promotion", promotion)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func List(w http.ResponseWriter, r *http.Request) {
	data, err := list()
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.Response(w, 200, data)
}

func Create(w http.ResponseWriter, r *http.Request) {

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	var newPromo PromotionExtra
	json.Unmarshal(reqBody, &newPromo)

	t, _ := time.Parse("02/01/2006", *newPromo.StartDate)
	*newPromo.StartDate = t.Format("2006-01-02")

	t, _ = time.Parse("02/01/2006", *newPromo.EndDate)
	*newPromo.EndDate = t.Format("2006-01-02")

	id, err := create(newPromo)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.ResponseCreated(w, id)
}

func Read(w http.ResponseWriter, r *http.Request) {
	promotion := r.Context().Value("promotion")
	utils.Response(w, http.StatusOK, promotion)
}

func Update(w http.ResponseWriter, r *http.Request) {

}

func Delete(w http.ResponseWriter, r *http.Request) {

}
