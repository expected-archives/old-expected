package apiserver

import (
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/models/plans"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *ApiServer) GetPlans(w http.ResponseWriter, r *http.Request) {
	plansType := mux.Vars(r)["type"]
	if plansType != string(plans.TypeContainer) && plansType != string(plans.TypeImage) {
		response.ErrorBadRequest(w, "Invalid plan type.", nil)
		return
	}
	p, err := plans.FindPlansByType(r.Context(), plansType)
	if err != nil {
		logrus.WithError(err).Errorln("unable to get plans")
		response.ErrorInternal(w)
		return
	}
	if p == nil {
		p = []*plans.Plan{}
	}
	response.Resource(w, "plans", p)
}
