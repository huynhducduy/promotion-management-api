package promotion

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
	"promotion-management-api/pkg/utils"
	"strconv"
	"time"
)

func RouterList(w http.ResponseWriter, r *http.Request) {
	data, err := List()
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.Response(w, 200, data)
}

func RouterCreate(w http.ResponseWriter, r *http.Request) {

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

	id, err := Create(newPromo)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.ResponseCreated(w, id)
}

func RouterRead(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ResponseMessage(w, http.StatusBadRequest, "Id must be an integer!")
		return
	}

	promotion, err := Read(id)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.Response(w, http.StatusOK, promotion)
}

func RouterUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ResponseMessage(w, http.StatusBadRequest, "Id must be an integer!")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	var updatedPromo PromotionExtra
	json.Unmarshal(reqBody, &updatedPromo)
	updatedPromo.Id = &id

	t, _ := time.Parse("02/01/2006", *updatedPromo.StartDate)
	*updatedPromo.StartDate = t.Format("2006-01-02")

	t, _ = time.Parse("02/01/2006", *updatedPromo.EndDate)
	*updatedPromo.EndDate = t.Format("2006-01-02")

	result, err := Update(updatedPromo)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.Response(w, http.StatusOK, result)
}

func RouterDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ResponseMessage(w, http.StatusBadRequest, "Id must be an integer!")
		return
	}

	err = Delete(id)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.ResponseMessage(w, http.StatusOK, "Delete succeed!")
}

func RouterApplicable(w http.ResponseWriter, r *http.Request) {
	storeId, err := strconv.ParseInt(r.URL.Query().Get("store_id"), 10, 64)
	if err != nil {
		utils.ResponseMessage(w, http.StatusBadRequest, "Store id must be an integer!")
		return
	}

	var memberId int64 = 0

	if r.URL.Query().Get("member_id") != "" {
		memberId, err = strconv.ParseInt(r.URL.Query().Get("member_id"), 10, 64)
		if err != nil {
			utils.ResponseMessage(w, http.StatusBadRequest, "Member id must be an integer!")
			return
		}
	}

	paymentType := r.URL.Query().Get("payment_type")
	orderType := r.URL.Query().Get("order_type")

	promotions, err := Applicable(storeId, memberId, paymentType, orderType)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.Response(w, http.StatusOK, promotions)
}
