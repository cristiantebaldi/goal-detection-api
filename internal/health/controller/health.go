package controller

import (
	"encoding/json"
	"net/http"

	"github.com/cristiantebaldi/goal-detection-api/internal/health/contract"
	"github.com/cristiantebaldi/goal-detection-api/internal/health/model"
)

type healthController struct{}

func (controller *healthController) Execute(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(model.Health{
		Status:  "ok",
		Message: "service is running",
	})
}

func NewHealthController() contract.HealthController {
	return &healthController{}
}
