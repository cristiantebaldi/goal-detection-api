package contract

import "net/http"

type HealthController interface {
	Execute(http.ResponseWriter, *http.Request)
}
