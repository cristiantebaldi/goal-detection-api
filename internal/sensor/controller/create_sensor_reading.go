package controller

import (
	"encoding/json"
	"net/http"

	"github.com/cristiantebaldi/goal-detection-api/internal/sensor/contract"
	"github.com/cristiantebaldi/goal-detection-api/internal/sensor/model"
	"github.com/cristiantebaldi/goal-detection-api/internal/util"
	"github.com/cristiantebaldi/goal-detection-api/pkg/validator"
)

type createSensorReadingController struct {
	createSensorReadingUseCase contract.CreateSensorReadingUseCase
}

func (controller *createSensorReadingController) Execute(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	var input model.CreateSensorReadingInput
	err := json.NewDecoder(request.Body).Decode(&input)
	if err != nil {
		util.HandleJSONError(response, err)
		return
	}

	if err := validator.ValidateStruct(&input); err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}

	output, err := controller.createSensorReadingUseCase.Execute(ctx, input)
	if err != nil {
		util.HandleAPIError(response, err)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(output)
}

func NewCreateSensorReadingController(createSensorReadingUseCase contract.CreateSensorReadingUseCase) contract.CreateSensorReadingController {
	return &createSensorReadingController{
		createSensorReadingUseCase: createSensorReadingUseCase,
	}
}
