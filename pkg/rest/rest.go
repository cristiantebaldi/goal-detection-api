package rest

import (
	"fmt"
	"net/http"

	"github.com/cristiantebaldi/goal-detection-api/internal/middleware"
	"github.com/cristiantebaldi/goal-detection-api/pkg/config"
	"github.com/cristiantebaldi/goal-detection-api/pkg/di"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

func Initialize() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("fatal error loading config: %w", err))
	}

	port := viper.GetString("server.http.port")
	if port == "" {
		port = cfg.Server.HTTP.Port
	}

	goalEventHub := di.NewGoalEventHub()
	healthController := di.NewHealthController()
	createSensorReadingController := di.NewCreateSensorReadingController(cfg, goalEventHub)
	listStreamGoalEventController := di.NewListStreamGoalEventController(goalEventHub)

	router := chi.NewRouter()
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	middleware.Cors(router)

	router.Group(func(publicRouter chi.Router) {
		publicRouter.Get("/health", healthController.Execute)
		publicRouter.Post("/sensor/position", createSensorReadingController.Execute)
		publicRouter.Get("/goal/event", listStreamGoalEventController.Execute)
	})

	fmt.Println("Server running on port " + port)
	http.ListenAndServe(":"+port, router)
}
