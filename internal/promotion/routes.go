package promotion

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"promotion-management-api/pkg/utils"
	"time"
)

func List(w http.ResponseWriter, r *http.Request) {
	data, err := getAll()
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

	var newPromo PromotionWithConstrains
	json.Unmarshal(reqBody, &newPromo)

	t, _ := time.Parse("02/01/2006", *newPromo.StartDate)
	*newPromo.StartDate = t.Format("2006-01-02")

	t, _ = time.Parse("02/01/2006", *newPromo.EndDate)
	*newPromo.EndDate = t.Format("2006-01-02")

	id, err := insert(newPromo)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.ResponseCreated(w, id)
}

func Read(w http.ResponseWriter, r *http.Request) {

}
