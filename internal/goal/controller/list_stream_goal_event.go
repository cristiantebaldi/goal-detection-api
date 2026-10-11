package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cristiantebaldi/goal-detection-api/internal/goal/contract"
	"github.com/cristiantebaldi/goal-detection-api/internal/goal/model"
	"github.com/cristiantebaldi/goal-detection-api/internal/util"
	"github.com/cristiantebaldi/goal-detection-api/pkg/formatter"
)

type listStreamGoalEventController struct {
	listStreamGoalEventUseCase contract.ListStreamGoalEventUseCase
}

func (controller *listStreamGoalEventController) Execute(response http.ResponseWriter, request *http.Request) {
	flusher, ok := response.(http.Flusher)
	if !ok {
		formatter.HttpErrorResponse(response, formatter.ErrInternalServer, "Streaming unsupported")
		return
	}

	events, err := controller.listStreamGoalEventUseCase.Execute(request.Context())
	if err != nil {
		util.HandleAPIError(response, err)
		return
	}

	response.Header().Set("Content-Type", "text/event-stream")
	response.Header().Set("Cache-Control", "no-cache")
	response.Header().Set("Connection", "keep-alive")
	response.Header().Set("X-Accel-Buffering", "no")
	response.WriteHeader(http.StatusOK)
	flusher.Flush()

	heartbeat := time.NewTicker(model.SSEHeartbeatInterval)
	defer heartbeat.Stop()

	for {
		select {
		case <-request.Context().Done():
			return
		case <-heartbeat.C:
			if _, err := fmt.Fprintf(response, ": ping\n\n"); err != nil {
				return
			}
			flusher.Flush()
		case event, ok := <-events:
			if !ok {
				return
			}

			payload, err := json.Marshal(event.ToStreamEvent())
			if err != nil {
				continue
			}

			if _, err := fmt.Fprintf(response, "data: %s\n\n", payload); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func NewListStreamGoalEventController(listStreamGoalEventUseCase contract.ListStreamGoalEventUseCase) contract.ListStreamGoalEventController {
	return &listStreamGoalEventController{
		listStreamGoalEventUseCase: listStreamGoalEventUseCase,
	}
}
