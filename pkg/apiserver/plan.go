package apiserver

import (
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/models/plans"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *ApiServer) ListPlans(w http.ResponseWriter, r *http.Request) {
	var res []*plans.Plan
	var err error

	if plansType := r.URL.Query()["type"]; len(plansType) > 0 {
		if plansType[0] != string(plans.TypeContainer) && plansType[0] != string(plans.TypeImage) {
			response.ErrorBadRequest(w, "Invalid plan type.", nil)
			return
		}
		res, err = plans.FindPlansByType(r.Context(), plansType[0])
	} else {
		res, err = plans.FindPlans(r.Context())
	}
	if err != nil {
		logrus.WithError(err).Errorln("unable to get plans")
		response.ErrorInternal(w)
		return
	}
	if res == nil {
		res = []*plans.Plan{}
	}
	response.Resource(w, "plans", res)
}

func (s *ApiServer) GetPlan(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	plan, err := plans.FindPlanByID(r.Context(), id)
	if err != nil {
		logrus.WithError(err).Errorln("unable to get plans")
		response.ErrorInternal(w)
		return
	}
	if plan == nil {
		response.ErrorNotFound(w)
		return
	}
	response.Resource(w, "plan", plan)
}
