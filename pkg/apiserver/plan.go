package apiserver

import (
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/models/plans"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *ApiServer) GetPlans(w http.ResponseWriter, r *http.Request) {
	plansType := mux.Vars(r)["type"]
	if plansType != plans.TypeContainer && plansType != plans.TypeImage {
		response.ErrorBadRequest(w, "Invalid plan type.", nil)
		return
	}
	plans, err := containers.FindPlans(r.Context())
	if err != nil {
		logrus.WithError(err).Errorln("unable to get container plans")
		response.ErrorInternal(w)
		return
	}
	response.Resource(w, "plans", plans)
}
