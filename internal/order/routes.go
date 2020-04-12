package order

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

	newOrder := OrderExtra{}
	newOrder.Order = new(Order)
	json.Unmarshal(reqBody, &newOrder)

	now := time.Now().Format("2006-01-02 15:04:05")
	newOrder.CreatedTime = &now

	id, err := Create(newOrder)
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

	order, err := Read(id)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.Response(w, http.StatusOK, order)
}
